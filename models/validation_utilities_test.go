package models

import (
	"strings"
	"testing"
	"time"
)

func TestStateMachineTraverser(t *testing.T) {
	t.Run("TraverseStateMachine basic functionality", func(t *testing.T) {
		// Create a simple state machine for testing
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "state1",
								Name: "State 1",
								Type: "state",
							},
						},
					},
					Vertices: []*Vertex{
						{
							ID:   "vertex1",
							Name: "Vertex 1",
							Type: "pseudostate",
						},
					},
				},
			},
		}

		traverser := NewStateMachineTraverser()
		visitedObjects := make([]string, 0)

		err := traverser.TraverseStateMachine(sm, func(obj interface{}, path []string, depth int) error {
			// Record the object type and path
			objType := ""
			switch obj.(type) {
			case *StateMachine:
				objType = "StateMachine"
			case *Region:
				objType = "Region"
			case *State:
				objType = "State"
			case *Vertex:
				objType = "Vertex"
			}
			visitedObjects = append(visitedObjects, objType)
			return nil
		})

		if err != nil {
			t.Fatalf("TraverseStateMachine failed: %v", err)
		}

		// Verify that all expected objects were visited
		expectedObjects := []string{"StateMachine", "Region", "State", "Vertex"}
		if len(visitedObjects) != len(expectedObjects) {
			t.Errorf("Expected %d objects, got %d", len(expectedObjects), len(visitedObjects))
		}

		// Check that StateMachine was visited first
		if len(visitedObjects) > 0 && visitedObjects[0] != "StateMachine" {
			t.Errorf("Expected StateMachine to be visited first, got %s", visitedObjects[0])
		}
	})

	t.Run("TraverseStateMachine with nil input", func(t *testing.T) {
		traverser := NewStateMachineTraverser()

		err := traverser.TraverseStateMachine(nil, func(obj interface{}, path []string, depth int) error {
			return nil
		})

		if err == nil {
			t.Error("Expected error for nil state machine")
		}
	})

	t.Run("TraverseStateMachine with cycle detection", func(t *testing.T) {
		// Create a state machine with potential cycles
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
				},
			},
		}

		// Add the same region twice to test cycle detection
		sm.Regions = append(sm.Regions, sm.Regions[0])

		traverser := NewStateMachineTraverser()
		visitCount := 0

		err := traverser.TraverseStateMachine(sm, func(obj interface{}, path []string, depth int) error {
			if _, ok := obj.(*Region); ok {
				visitCount++
			}
			return nil
		})

		if err != nil {
			t.Fatalf("TraverseStateMachine failed: %v", err)
		}

		// Should only visit the region once due to cycle detection
		if visitCount != 1 {
			t.Errorf("Expected region to be visited once, got %d visits", visitCount)
		}
	})

	t.Run("TraverseRegion functionality", func(t *testing.T) {
		region := &Region{
			ID:   "test-region",
			Name: "Test Region",
			States: []*State{
				{
					Vertex: Vertex{
						ID:   "state1",
						Name: "State 1",
						Type: "state",
					},
				},
			},
			Transitions: []*Transition{
				{
					ID: "transition1",
					Source: &Vertex{
						ID:   "source",
						Name: "Source",
						Type: "state",
					},
					Target: &Vertex{
						ID:   "target",
						Name: "Target",
						Type: "state",
					},
				},
			},
		}

		traverser := NewStateMachineTraverser()
		visitedTypes := make([]string, 0)

		err := traverser.TraverseRegion(region, func(obj interface{}, path []string, depth int) error {
			switch obj.(type) {
			case *Region:
				visitedTypes = append(visitedTypes, "Region")
			case *State:
				visitedTypes = append(visitedTypes, "State")
			case *Transition:
				visitedTypes = append(visitedTypes, "Transition")
			}
			return nil
		})

		if err != nil {
			t.Fatalf("TraverseRegion failed: %v", err)
		}

		// Should visit Region, State, and Transition
		expectedTypes := []string{"Region", "State", "Transition"}
		if len(visitedTypes) != len(expectedTypes) {
			t.Errorf("Expected %d types, got %d: %v", len(expectedTypes), len(visitedTypes), visitedTypes)
		}
	})
}

func TestValidationResultAggregator(t *testing.T) {
	t.Run("AddResult and GetResults", func(t *testing.T) {
		aggregator := NewValidationResultAggregator()

		// Create some validation errors
		errors1 := &ValidationErrors{}
		errors1.AddError(ErrorTypeRequired, "Object1", "Field1", "Error 1", []string{"path1"})

		errors2 := &ValidationErrors{}
		errors2.AddError(ErrorTypeInvalid, "Object2", "Field2", "Error 2", []string{"path2"})

		// Add results
		aggregator.AddResult("obj1", errors1)
		aggregator.AddResult("obj2", errors2)

		// Verify results
		results := aggregator.GetResults()
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		if !aggregator.HasErrors() {
			t.Error("HasErrors should return true")
		}

		totalErrors := aggregator.GetTotalErrorCount()
		if totalErrors != 2 {
			t.Errorf("Expected 2 total errors, got %d", totalErrors)
		}
	})

	t.Run("AddSingleError", func(t *testing.T) {
		aggregator := NewValidationResultAggregator()

		aggregator.AddSingleError("obj1", ErrorTypeRequired, "TestObject", "TestField", "Test message", []string{"test", "path"})

		results := aggregator.GetResults()
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		errors := results["obj1"]
		if len(errors.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errors.Errors))
		}

		err := errors.Errors[0]
		if err.Type != ErrorTypeRequired {
			t.Errorf("Expected ErrorTypeRequired, got %v", err.Type)
		}
	})

	t.Run("GetSummaryReport", func(t *testing.T) {
		aggregator := NewValidationResultAggregator()

		// Add some errors
		aggregator.AddSingleError("obj1", ErrorTypeRequired, "Object1", "Field1", "Error 1", []string{"path1"})
		aggregator.AddSingleError("obj2", ErrorTypeInvalid, "Object2", "Field2", "Error 2", []string{"path2"})

		report := aggregator.GetSummaryReport()

		// Verify report contains expected information
		if !strings.Contains(report, "Validation Summary") {
			t.Error("Report should contain 'Validation Summary'")
		}
		if !strings.Contains(report, "2 error(s)") {
			t.Error("Report should contain '2 error(s)'")
		}
		if !strings.Contains(report, "obj1") {
			t.Error("Report should contain 'obj1'")
		}
		if !strings.Contains(report, "obj2") {
			t.Error("Report should contain 'obj2'")
		}
	})

	t.Run("GetDetailedReport", func(t *testing.T) {
		aggregator := NewValidationResultAggregator()

		// Add errors with context
		errors := &ValidationErrors{}
		errors.AddErrorWithContext(ErrorTypeRequired, "TestObject", "TestField", "Test message", []string{"test", "path"}, map[string]interface{}{
			"context_key": "context_value",
		})
		aggregator.AddResult("obj1", errors)

		report := aggregator.GetDetailedReport()

		// Verify detailed report contains expected information
		if !strings.Contains(report, "Detailed Validation Report") {
			t.Error("Report should contain 'Detailed Validation Report'")
		}
		if !strings.Contains(report, "Generated:") {
			t.Error("Report should contain 'Generated:'")
		}
		if !strings.Contains(report, "Required Errors") {
			t.Error("Report should contain 'Required Errors'")
		}
		if !strings.Contains(report, "context_key=context_value") {
			t.Error("Report should contain context information")
		}
	})

	t.Run("Merge functionality", func(t *testing.T) {
		aggregator1 := NewValidationResultAggregator()
		aggregator2 := NewValidationResultAggregator()

		// Add errors to both aggregators
		aggregator1.AddSingleError("obj1", ErrorTypeRequired, "Object1", "Field1", "Error 1", []string{"path1"})
		aggregator2.AddSingleError("obj2", ErrorTypeInvalid, "Object2", "Field2", "Error 2", []string{"path2"})

		// Merge aggregator2 into aggregator1
		aggregator1.Merge(aggregator2)

		// Verify merged results
		results := aggregator1.GetResults()
		if len(results) != 2 {
			t.Errorf("Expected 2 results after merge, got %d", len(results))
		}

		totalErrors := aggregator1.GetTotalErrorCount()
		if totalErrors != 2 {
			t.Errorf("Expected 2 total errors after merge, got %d", totalErrors)
		}
	})

	t.Run("Clear functionality", func(t *testing.T) {
		aggregator := NewValidationResultAggregator()

		// Add some errors
		aggregator.AddSingleError("obj1", ErrorTypeRequired, "Object1", "Field1", "Error 1", []string{"path1"})

		// Verify errors exist
		if !aggregator.HasErrors() {
			t.Error("Should have errors before clear")
		}

		// Clear and verify
		aggregator.Clear()

		if aggregator.HasErrors() {
			t.Error("Should not have errors after clear")
		}

		if aggregator.GetTotalErrorCount() != 0 {
			t.Error("Should have 0 errors after clear")
		}
	})
}

func TestValidationDebugger(t *testing.T) {
	t.Run("DebugStateMachine basic functionality", func(t *testing.T) {
		// Create a test state machine
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "state1",
								Name: "State 1",
								Type: "state",
							},
						},
					},
				},
			},
		}

		debugger := NewValidationDebugger()
		report, err := debugger.DebugStateMachine(sm)

		if err != nil {
			t.Fatalf("DebugStateMachine failed: %v", err)
		}

		if report == nil {
			t.Fatal("Expected debug report, got nil")
		}

		// Verify report contents
		if report.StateMachineID != "test-sm" {
			t.Errorf("Expected StateMachineID 'test-sm', got '%s'", report.StateMachineID)
		}

		if report.TotalObjects == 0 {
			t.Error("Expected TotalObjects > 0")
		}

		if len(report.Objects) == 0 {
			t.Error("Expected Objects map to be populated")
		}

		// Verify timestamp is recent
		if time.Since(report.Timestamp) > time.Minute {
			t.Error("Report timestamp should be recent")
		}
	})

	t.Run("DebugStateMachine with nil input", func(t *testing.T) {
		debugger := NewValidationDebugger()
		report, err := debugger.DebugStateMachine(nil)

		if err == nil {
			t.Error("Expected error for nil state machine")
		}

		if report != nil {
			t.Error("Expected nil report for nil state machine")
		}
	})

	t.Run("DebugStateMachine with validation errors", func(t *testing.T) {
		// Create an invalid state machine (missing required fields)
		sm := &StateMachine{
			// Missing ID, Name, Version
			Regions: []*Region{
				{
					// Missing ID, Name
				},
			},
		}

		debugger := NewValidationDebugger()
		report, err := debugger.DebugStateMachine(sm)

		if err != nil {
			t.Fatalf("DebugStateMachine failed: %v", err)
		}

		// Should have validation errors
		if report.TotalErrors == 0 {
			t.Error("Expected validation errors for invalid state machine")
		}

		if len(report.ValidationResults) == 0 {
			t.Error("Expected ValidationResults to be populated")
		}
	})

	t.Run("ValidationDebugReport GetSummary", func(t *testing.T) {
		report := &ValidationDebugReport{
			StateMachineID: "test-sm",
			Timestamp:      time.Now(),
			TotalObjects:   5,
			TotalErrors:    2,
			Objects: map[string]*ObjectDebugInfo{
				"obj1": {
					ID:   "obj1",
					Type: "StateMachine",
					Path: "StateMachine",
				},
				"obj2": {
					ID:   "obj2",
					Type: "Region",
					Path: "StateMachine.Regions[0]",
				},
			},
			ValidationResults: map[string]*ValidationErrors{
				"test-sm": {
					Errors: []*ValidationError{
						{
							Type:    ErrorTypeRequired,
							Object:  "StateMachine",
							Field:   "ID",
							Message: "field is required",
						},
					},
				},
			},
		}

		summary := report.GetSummary()

		// Verify summary contains expected information
		if !strings.Contains(summary, "test-sm") {
			t.Error("Summary should contain StateMachine ID")
		}
		if !strings.Contains(summary, "Total Objects: 5") {
			t.Error("Summary should contain total objects count")
		}
		if !strings.Contains(summary, "Total Errors: 2") {
			t.Error("Summary should contain total errors count")
		}
		if !strings.Contains(summary, "Object Distribution") {
			t.Error("Summary should contain object distribution")
		}
		if !strings.Contains(summary, "StateMachine: 1") {
			t.Error("Summary should contain StateMachine count")
		}
		if !strings.Contains(summary, "Region: 1") {
			t.Error("Summary should contain Region count")
		}
	})
}

func TestCommonValidationPatterns(t *testing.T) {
	t.Run("ValidateStateMachineStructure", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test with valid state machine
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
				},
			},
		}

		patterns.ValidateStateMachineStructure(sm, context, errors)

		if errors.HasErrors() {
			t.Errorf("Valid state machine should not have errors: %v", errors.Error())
		}

		// Test with invalid state machine (missing required fields)
		errors.Clear()
		invalidSM := &StateMachine{
			// Missing ID, Name, Version
			Regions: []*Region{}, // Empty regions (violates UML constraint)
		}

		patterns.ValidateStateMachineStructure(invalidSM, context, errors)

		if !errors.HasErrors() {
			t.Error("Invalid state machine should have errors")
		}

		// Should have errors for missing ID, Name, Version, and empty regions
		if len(errors.Errors) < 4 {
			t.Errorf("Expected at least 4 errors, got %d", len(errors.Errors))
		}
	})

	t.Run("ValidateStateMachineStructure with connection points", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test with invalid connection points
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
				},
			},
			ConnectionPoints: []*Pseudostate{
				{
					Vertex: Vertex{
						ID:   "cp1",
						Name: "Connection Point 1",
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial, // Invalid for connection point
				},
			},
		}

		patterns.ValidateStateMachineStructure(sm, context, errors)

		if !errors.HasErrors() {
			t.Error("Should have error for invalid connection point kind")
		}

		// Check for specific error about connection point kind
		found := false
		for _, err := range errors.Errors {
			if strings.Contains(err.Message, "connection point must be entry or exit point") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should have specific error about connection point kind")
		}
	})

	t.Run("ValidateStateMachineStructure method constraints", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test method state machine with connection points (invalid)
		sm := &StateMachine{
			ID:       "test-sm",
			Name:     "Test StateMachine",
			Version:  "1.0",
			IsMethod: true,
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
				},
			},
			ConnectionPoints: []*Pseudostate{
				{
					Vertex: Vertex{
						ID:   "cp1",
						Name: "Connection Point 1",
						Type: "pseudostate",
					},
					Kind: PseudostateKindEntryPoint,
				},
			},
		}

		patterns.ValidateStateMachineStructure(sm, context, errors)

		if !errors.HasErrors() {
			t.Error("Method state machine with connection points should have error")
		}

		// Check for specific error about method constraints
		found := false
		for _, err := range errors.Errors {
			if strings.Contains(err.Message, "state machine used as method cannot have connection points") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should have specific error about method constraints")
		}
	})

	t.Run("ValidateRegionStructure", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test with valid region
		region := &Region{
			ID:   "test-region",
			Name: "Test Region",
			Vertices: []*Vertex{
				{
					ID:   "initial",
					Name: "Initial",
					Type: "pseudostate",
				},
			},
		}

		patterns.ValidateRegionStructure(region, context, errors)

		if errors.HasErrors() {
			t.Errorf("Valid region should not have errors: %v", errors.Error())
		}

		// Test with invalid region (missing required fields)
		errors.Clear()
		invalidRegion := &Region{
			// Missing ID, Name
		}

		patterns.ValidateRegionStructure(invalidRegion, context, errors)

		if !errors.HasErrors() {
			t.Error("Invalid region should have errors")
		}
	})

	t.Run("ValidateRegionStructure multiple initial states", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test with multiple initial pseudostates (invalid)
		region := &Region{
			ID:   "test-region",
			Name: "Test Region",
			Vertices: []*Vertex{
				{
					ID:   "initial1",
					Name: "Initial1",
					Type: "pseudostate",
				},
				{
					ID:   "initial2",
					Name: "Initial2",
					Type: "pseudostate",
				},
			},
		}

		patterns.ValidateRegionStructure(region, context, errors)

		if !errors.HasErrors() {
			t.Error("Region with multiple initial states should have error")
		}

		// Check for specific error about multiple initial states
		found := false
		for _, err := range errors.Errors {
			if strings.Contains(err.Message, "region can have at most one initial pseudostate") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should have specific error about multiple initial states")
		}
	})

	t.Run("ValidateTransitionStructure", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test with valid transition
		transition := &Transition{
			ID: "test-transition",
			Source: &Vertex{
				ID:   "source",
				Name: "Source",
				Type: "state",
			},
			Target: &Vertex{
				ID:   "target",
				Name: "Target",
				Type: "state",
			},
			Kind: TransitionKindExternal,
		}

		patterns.ValidateTransitionStructure(transition, context, errors)

		if errors.HasErrors() {
			t.Errorf("Valid transition should not have errors: %v", errors.Error())
		}

		// Test with invalid internal transition (different source and target)
		errors.Clear()
		invalidTransition := &Transition{
			ID: "test-transition",
			Source: &Vertex{
				ID:   "source",
				Name: "Source",
				Type: "state",
			},
			Target: &Vertex{
				ID:   "target", // Different from source
				Name: "Target",
				Type: "state",
			},
			Kind: TransitionKindInternal, // Internal transitions must have same source and target
		}

		patterns.ValidateTransitionStructure(invalidTransition, context, errors)

		if !errors.HasErrors() {
			t.Error("Invalid internal transition should have error")
		}

		// Check for specific error about internal transition constraints
		found := false
		for _, err := range errors.Errors {
			if strings.Contains(err.Message, "internal transition must have same source and target") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Should have specific error about internal transition constraints")
		}
	})

	t.Run("ValidateObjectHierarchy", func(t *testing.T) {
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		// Test with valid state machine
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
				},
			},
		}

		patterns.ValidateObjectHierarchy(sm, context, errors)

		if errors.HasErrors() {
			t.Errorf("Valid hierarchy should not have errors: %v", errors.Error())
		}

		// Test with object missing ID
		errors.Clear()
		invalidSM := &StateMachine{
			// Missing ID
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "region1",
					Name: "Region 1",
				},
			},
		}

		patterns.ValidateObjectHierarchy(invalidSM, context, errors)

		if !errors.HasErrors() {
			t.Error("Hierarchy with missing ID should have error")
		}
	})
}

func TestValidationUtilities(t *testing.T) {
	t.Run("NewValidationUtilities", func(t *testing.T) {
		utils := NewValidationUtilities()
		if utils == nil {
			t.Fatal("NewValidationUtilities returned nil")
		}
		if utils.helper == nil {
			t.Error("ValidationUtilities should have helper")
		}
	})

	t.Run("Integration test with real state machine", func(t *testing.T) {
		// Create a comprehensive test state machine
		sm := &StateMachine{
			ID:      "integration-test-sm",
			Name:    "Integration Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "main-region",
					Name: "Main Region",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "state1",
								Name: "State 1",
								Type: "state",
							},
							IsComposite: true,
							Regions: []*Region{
								{
									ID:   "nested-region",
									Name: "Nested Region",
									Vertices: []*Vertex{
										{
											ID:   "nested-initial",
											Name: "Initial",
											Type: "pseudostate",
										},
									},
								},
							},
						},
					},
					Vertices: []*Vertex{
						{
							ID:   "initial",
							Name: "Initial",
							Type: "pseudostate",
						},
					},
					Transitions: []*Transition{
						{
							ID:   "transition1",
							Name: "Transition 1",
							Source: &Vertex{
								ID:   "initial",
								Name: "Initial",
								Type: "pseudostate",
							},
							Target: &Vertex{
								ID:   "state1",
								Name: "State 1",
								Type: "state",
							},
							Kind: TransitionKindExternal,
						},
					},
				},
			},
		}

		// Test traversal
		traverser := NewStateMachineTraverser()
		objectCount := 0
		err := traverser.TraverseStateMachine(sm, func(obj interface{}, path []string, depth int) error {
			objectCount++
			return nil
		})

		if err != nil {
			t.Fatalf("Traversal failed: %v", err)
		}

		if objectCount == 0 {
			t.Error("Should have traversed some objects")
		}

		// Test validation aggregation
		aggregator := NewValidationResultAggregator()
		validationErr := sm.Validate()
		if validationErr != nil {
			if validationErrors, ok := validationErr.(*ValidationErrors); ok {
				aggregator.AddResult(sm.ID, validationErrors)
			}
		}

		// Test debugging
		debugger := NewValidationDebugger()
		report, err := debugger.DebugStateMachine(sm)
		if err != nil {
			t.Fatalf("Debugging failed: %v", err)
		}

		if report.TotalObjects == 0 {
			t.Error("Debug report should have objects")
		}

		// Test common patterns
		patterns := NewCommonValidationPatterns()
		context := NewValidationContext()
		errors := &ValidationErrors{}

		patterns.ValidateStateMachineStructure(sm, context, errors)
		patterns.ValidateObjectHierarchy(sm, context, errors)

		// The state machine should be valid
		if errors.HasErrors() {
			t.Errorf("Integration test state machine should be valid: %v", errors.Error())
		}
	})
}
