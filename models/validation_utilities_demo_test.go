package models

import (
	"testing"
)

func TestValidationUtilitiesDemo(t *testing.T) {
	t.Run("ValidationUtilitiesDemo runs without panic", func(t *testing.T) {
		// This test ensures the demo runs without panicking
		// We can't easily test the output, but we can ensure it doesn't crash
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ValidationUtilitiesDemo panicked: %v", r)
			}
		}()

		ValidationUtilitiesDemo()
	})

	t.Run("createDemoStateMachine creates valid state machine", func(t *testing.T) {
		sm := createDemoStateMachine()

		if sm == nil {
			t.Fatal("createDemoStateMachine returned nil")
		}

		if sm.ID == "" {
			t.Error("Demo state machine should have ID")
		}

		if sm.Name == "" {
			t.Error("Demo state machine should have Name")
		}

		if len(sm.Regions) == 0 {
			t.Error("Demo state machine should have regions")
		}

		// Validate the demo state machine
		err := sm.Validate()
		if err != nil {
			t.Errorf("Demo state machine should be valid: %v", err)
		}
	})

	t.Run("demonstrateTraversal works with valid state machine", func(t *testing.T) {
		sm := createDemoStateMachine()

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("demonstrateTraversal panicked: %v", r)
			}
		}()

		demonstrateTraversal(sm)
	})

	t.Run("demonstrateValidationAggregation works with valid state machine", func(t *testing.T) {
		sm := createDemoStateMachine()

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("demonstrateValidationAggregation panicked: %v", r)
			}
		}()

		demonstrateValidationAggregation(sm)
	})

	t.Run("demonstrateValidationDebugging works with valid state machine", func(t *testing.T) {
		sm := createDemoStateMachine()

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("demonstrateValidationDebugging panicked: %v", r)
			}
		}()

		demonstrateValidationDebugging(sm)
	})

	t.Run("demonstrateCommonPatterns works with valid state machine", func(t *testing.T) {
		sm := createDemoStateMachine()

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("demonstrateCommonPatterns panicked: %v", r)
			}
		}()

		demonstrateCommonPatterns(sm)
	})

	t.Run("demonstrateErrorScenarios works with invalid state machine", func(t *testing.T) {
		// This should not panic even with invalid data
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("demonstrateErrorScenarios panicked: %v", r)
			}
		}()

		demonstrateErrorScenarios()
	})

	t.Run("helper functions work correctly", func(t *testing.T) {
		sm := createDemoStateMachine()

		// Test getObjectTypeName
		typeName := getObjectTypeName(sm)
		if typeName != "StateMachine" {
			t.Errorf("Expected 'StateMachine', got '%s'", typeName)
		}

		// Test getObjectID
		id := getObjectID(sm)
		if id != sm.ID {
			t.Errorf("Expected '%s', got '%s'", sm.ID, id)
		}

		// Test with nil
		nilTypeName := getObjectTypeName(nil)
		if nilTypeName != "Unknown" {
			t.Errorf("Expected 'Unknown' for nil, got '%s'", nilTypeName)
		}

		nilID := getObjectID(nil)
		if nilID != "" {
			t.Errorf("Expected empty string for nil ID, got '%s'", nilID)
		}
	})
}

// TestValidationUtilitiesDemoIntegration tests the demo with a more comprehensive approach
func TestValidationUtilitiesDemoIntegration(t *testing.T) {
	t.Run("Full demo integration test", func(t *testing.T) {
		// Create demo state machine
		sm := createDemoStateMachine()

		// Test all components work together
		traverser := NewStateMachineTraverser()
		aggregator := NewValidationResultAggregator()
		debugger := NewValidationDebugger()
		patterns := NewCommonValidationPatterns()

		// 1. Traverse and count objects
		objectCount := 0
		err := traverser.TraverseStateMachine(sm, func(obj interface{}, path []string, depth int) error {
			objectCount++
			return nil
		})

		if err != nil {
			t.Fatalf("Traversal failed: %v", err)
		}

		if objectCount == 0 {
			t.Error("Should have found objects during traversal")
		}

		// 2. Validate and aggregate results
		validationErr := sm.Validate()
		if validationErr != nil {
			if validationErrors, ok := validationErr.(*ValidationErrors); ok {
				aggregator.AddResult(sm.ID, validationErrors)
			}
		}

		// 3. Generate debug report
		report, err := debugger.DebugStateMachine(sm)
		if err != nil {
			t.Fatalf("Debug failed: %v", err)
		}

		if report.TotalObjects == 0 {
			t.Error("Debug report should have objects")
		}

		// 4. Apply common patterns
		context := NewValidationContext()
		errors := &ValidationErrors{}
		patterns.ValidateStateMachineStructure(sm, context, errors)

		// The demo state machine should be valid
		if errors.HasErrors() {
			t.Errorf("Demo state machine should pass common patterns validation: %v", errors.Error())
		}

		// 5. Verify consistency between different validation approaches
		// The traversal object count should be reasonable compared to debug report
		if objectCount < report.TotalObjects {
			t.Errorf("Traversal found fewer objects (%d) than debug report (%d)", objectCount, report.TotalObjects)
		}
	})
}
