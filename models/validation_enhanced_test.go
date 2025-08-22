package models

import (
	"strings"
	"testing"
)

func TestEnhancedValidationErrorHandling(t *testing.T) {
	t.Run("comprehensive error collection", func(t *testing.T) {
		// Create a state machine with multiple validation errors
		sm := &StateMachine{
			ID:      "", // Missing required field
			Name:    "", // Missing required field
			Version: "", // Missing required field
			Regions: []*Region{
				{
					ID:   "", // Missing required field
					Name: "", // Missing required field
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "", // Missing required field
								Name: "", // Missing required field
								Type: "", // Missing required field
							},
						},
					},
					Transitions: []*Transition{
						{
							ID:     "",  // Missing required field
							Source: nil, // Missing required reference
							Target: nil, // Missing required reference
							Kind:   "",  // Invalid kind
						},
					},
				},
			},
		}

		// Validate and collect all errors
		err := sm.Validate()
		if err == nil {
			t.Fatal("Expected validation errors, got nil")
		}

		validationErrors, ok := err.(*ValidationErrors)
		if !ok {
			t.Fatalf("Expected ValidationErrors, got %T", err)
		}

		// Verify multiple errors were collected
		if validationErrors.Count() < 5 {
			t.Errorf("Expected at least 5 validation errors, got %d", validationErrors.Count())
		}

		// Verify error types are properly categorized
		summary := validationErrors.GetSummary()
		if summary[ErrorTypeRequired] == 0 {
			t.Error("Expected required field errors")
		}
		if summary[ErrorTypeInvalid] == 0 {
			t.Error("Expected invalid value errors")
		}

		// Verify detailed report generation
		report := validationErrors.GetDetailedReport()
		if !strings.Contains(report, "Validation Report:") {
			t.Error("Expected detailed report header")
		}
		if !strings.Contains(report, "Required Errors") {
			t.Error("Expected required errors section in report")
		}
	})

	t.Run("path tracking through hierarchy", func(t *testing.T) {
		// Create a nested structure with errors at different levels
		sm := &StateMachine{
			ID:      "test-sm",
			Name:    "Test StateMachine",
			Version: "1.0",
			Regions: []*Region{
				{
					ID:   "test-region",
					Name: "Test Region",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "test-state",
								Name: "Test State",
								Type: "state",
							},
							Entry: &Behavior{
								ID:            "", // Missing required field - should have path
								Name:          "Entry Behavior",
								Specification: "", // Missing required field - should have path
							},
						},
					},
				},
			},
		}

		err := sm.Validate()
		if err == nil {
			t.Fatal("Expected validation errors, got nil")
		}

		validationErrors, ok := err.(*ValidationErrors)
		if !ok {
			t.Fatalf("Expected ValidationErrors, got %T", err)
		}

		// Find errors with specific paths
		foundEntryIDError := false
		foundEntrySpecError := false

		for _, validationError := range validationErrors.Errors {
			pathStr := strings.Join(validationError.Path, ".")
			if strings.Contains(pathStr, "Regions[0].States[0].Entry") {
				if validationError.Field == "ID" {
					foundEntryIDError = true
				}
				if validationError.Field == "Specification" {
					foundEntrySpecError = true
				}
			}
		}

		if !foundEntryIDError {
			t.Error("Expected to find entry behavior ID error with proper path")
		}
		if !foundEntrySpecError {
			t.Error("Expected to find entry behavior specification error with proper path")
		}
	})

	t.Run("context information in errors", func(t *testing.T) {
		// Create a validation context with metadata
		context := NewValidationContext()
		context.SetMetadata("validationPhase", "UMLConstraints")
		context.SetMetadata("validationLevel", "detailed")

		errors := &ValidationErrors{}
		helper := NewValidationHelper()

		// Add an error with context
		helper.ValidateRequired("", "TestField", "TestObject", context, errors)

		if errors.Count() != 1 {
			t.Fatalf("Expected 1 error, got %d", errors.Count())
		}

		// Verify context information is available through the context
		contextInfo := context.GetContextInfo()
		if contextInfo["metadata"] == nil {
			t.Error("Expected metadata in context info")
		}

		metadata := contextInfo["metadata"].(map[string]interface{})
		if metadata["validationPhase"] != "UMLConstraints" {
			t.Error("Expected validation phase metadata")
		}
	})

	t.Run("helper methods for common patterns", func(t *testing.T) {
		context := NewValidationContext()
		errors := &ValidationErrors{}
		helper := NewValidationHelper()

		// Test unique ID validation
		objects := []interface{}{
			&Region{ID: "region1", Name: "Region 1"},
			&Region{ID: "region1", Name: "Region 2"}, // Duplicate ID
			&Region{ID: "region2", Name: "Region 3"},
		}

		helper.ValidateUniqueIDs(objects, "Regions", "StateMachine", context, errors, func(obj interface{}) string {
			if region, ok := obj.(*Region); ok {
				return region.ID
			}
			return ""
		})

		if errors.Count() != 1 {
			t.Errorf("Expected 1 duplicate ID error, got %d", errors.Count())
		}

		duplicateError := errors.Errors[0]
		if duplicateError.Type != ErrorTypeConstraint {
			t.Errorf("Expected constraint error type, got %s", duplicateError.Type)
		}
		if !strings.Contains(duplicateError.Message, "duplicate ID") {
			t.Errorf("Expected duplicate ID message, got: %s", duplicateError.Message)
		}

		// Test conditional required validation
		errors.Clear()
		helper.ValidateConditionalRequired("", "ConnectionPoints", "StateMachine", true, "state machine is used as method", context, errors)

		if errors.Count() != 1 {
			t.Errorf("Expected 1 conditional required error, got %d", errors.Count())
		}

		conditionalError := errors.Errors[0]
		if !strings.Contains(conditionalError.Message, "when state machine is used as method") {
			t.Errorf("Expected conditional message, got: %s", conditionalError.Message)
		}

		// Test collection size validation
		errors.Clear()
		emptyCollection := []*Region{}
		helper.ValidateCollectionSize(emptyCollection, "Regions", "StateMachine", 1, 0, context, errors)

		if errors.Count() != 1 {
			t.Errorf("Expected 1 collection size error, got %d", errors.Count())
		}

		sizeError := errors.Errors[0]
		if sizeError.Type != ErrorTypeMultiplicity {
			t.Errorf("Expected multiplicity error type, got %s", sizeError.Type)
		}
	})

	t.Run("error filtering and querying", func(t *testing.T) {
		errors := &ValidationErrors{}

		// Add various types of errors
		errors.AddError(ErrorTypeRequired, "StateMachine", "ID", "ID is required", []string{"StateMachine"})
		errors.AddError(ErrorTypeRequired, "Region", "Name", "Name is required", []string{"Regions", "0"})
		errors.AddError(ErrorTypeConstraint, "StateMachine", "Regions", "Must have at least one region", []string{"StateMachine"})
		errors.AddError(ErrorTypeInvalid, "Vertex", "Type", "Invalid type", []string{"Regions", "0", "Vertices", "0"})

		// Test filtering by type
		requiredErrors := errors.GetErrorsByType(ErrorTypeRequired)
		if len(requiredErrors) != 2 {
			t.Errorf("Expected 2 required errors, got %d", len(requiredErrors))
		}

		constraintErrors := errors.GetErrorsByType(ErrorTypeConstraint)
		if len(constraintErrors) != 1 {
			t.Errorf("Expected 1 constraint error, got %d", len(constraintErrors))
		}

		// Test filtering by object
		stateMachineErrors := errors.GetErrorsByObject("StateMachine")
		if len(stateMachineErrors) != 2 {
			t.Errorf("Expected 2 StateMachine errors, got %d", len(stateMachineErrors))
		}

		// Test filtering by path
		regionErrors := errors.GetErrorsByPath("Regions")
		if len(regionErrors) != 2 {
			t.Errorf("Expected 2 errors under Regions path, got %d", len(regionErrors))
		}

		// Test summary
		summary := errors.GetSummary()
		if summary[ErrorTypeRequired] != 2 {
			t.Errorf("Expected 2 required errors in summary, got %d", summary[ErrorTypeRequired])
		}
		if summary[ErrorTypeConstraint] != 1 {
			t.Errorf("Expected 1 constraint error in summary, got %d", summary[ErrorTypeConstraint])
		}
		if summary[ErrorTypeInvalid] != 1 {
			t.Errorf("Expected 1 invalid error in summary, got %d", summary[ErrorTypeInvalid])
		}
	})

	t.Run("error merging and manipulation", func(t *testing.T) {
		errors1 := &ValidationErrors{}
		errors1.AddError(ErrorTypeRequired, "Object1", "Field1", "Error 1", []string{"path1"})

		errors2 := &ValidationErrors{}
		errors2.AddError(ErrorTypeInvalid, "Object2", "Field2", "Error 2", []string{"path2"})

		// Test merging
		errors1.Merge(errors2)
		if errors1.Count() != 2 {
			t.Errorf("Expected 2 errors after merge, got %d", errors1.Count())
		}

		// Test clearing
		errors1.Clear()
		if !errors1.IsEmpty() {
			t.Error("Expected errors to be empty after clear")
		}
	})
}

func TestValidationContextEnhancements(t *testing.T) {
	t.Run("context cloning and metadata", func(t *testing.T) {
		// Create original context
		original := NewValidationContext()
		original.SetMetadata("key1", "value1")
		original.SetMetadata("key2", 42)

		// Clone the context
		cloned := original.Clone()

		// Modify the clone
		cloned.SetMetadata("key3", "value3")

		// Verify original is unchanged
		if _, exists := original.GetMetadata("key3"); exists {
			t.Error("Original context should not have key3")
		}

		// Verify clone has all metadata
		if value, exists := cloned.GetMetadata("key1"); !exists || value != "value1" {
			t.Error("Cloned context should have key1")
		}
		if value, exists := cloned.GetMetadata("key3"); !exists || value != "value3" {
			t.Error("Cloned context should have key3")
		}
	})

	t.Run("full path generation", func(t *testing.T) {
		sm := &StateMachine{ID: "sm1", Name: "StateMachine1"}
		region := &Region{ID: "r1", Name: "Region1"}

		context := NewValidationContext().
			WithStateMachine(sm).
			WithRegion(region).
			WithPath("States").
			WithPathIndex("Transitions", 0)

		fullPath := context.GetFullPath()
		expected := "StateMachine[sm1].Region[r1].States.Transitions[0]"
		if fullPath != expected {
			t.Errorf("Expected full path '%s', got '%s'", expected, fullPath)
		}
	})

	t.Run("context info generation", func(t *testing.T) {
		sm := &StateMachine{ID: "sm1", Name: "StateMachine1"}
		region := &Region{ID: "r1", Name: "Region1"}

		context := NewValidationContext().
			WithStateMachine(sm).
			WithRegion(region).
			WithPath("States").
			WithMetadata("phase", "validation")

		info := context.GetContextInfo()

		// Verify state machine info
		if smInfo, ok := info["stateMachine"].(map[string]interface{}); ok {
			if smInfo["id"] != "sm1" || smInfo["name"] != "StateMachine1" {
				t.Error("State machine info not correct")
			}
		} else {
			t.Error("Expected state machine info")
		}

		// Verify region info
		if regionInfo, ok := info["region"].(map[string]interface{}); ok {
			if regionInfo["id"] != "r1" || regionInfo["name"] != "Region1" {
				t.Error("Region info not correct")
			}
		} else {
			t.Error("Expected region info")
		}

		// Verify path info
		if pathInfo, ok := info["path"].([]string); ok {
			if len(pathInfo) != 1 || pathInfo[0] != "States" {
				t.Error("Path info not correct")
			}
		} else {
			t.Error("Expected path info")
		}

		// Verify metadata info
		if metadataInfo, ok := info["metadata"].(map[string]interface{}); ok {
			if metadataInfo["phase"] != "validation" {
				t.Error("Metadata info not correct")
			}
		} else {
			t.Error("Expected metadata info")
		}
	})
}
