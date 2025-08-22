package models

import (
	"strings"
	"testing"
)

func TestValidationInfrastructure(t *testing.T) {
	t.Run("ValidationContext creation and path tracking", func(t *testing.T) {
		ctx := NewValidationContext()
		if ctx == nil {
			t.Fatal("NewValidationContext() returned nil")
		}

		// Test path building
		ctx2 := ctx.WithPath("test")
		if len(ctx2.Path) != 1 || ctx2.Path[0] != "test" {
			t.Errorf("WithPath() failed, got path: %v", ctx2.Path)
		}

		ctx3 := ctx2.WithPathIndex("items", 5)
		if len(ctx3.Path) != 2 || ctx3.Path[1] != "items[5]" {
			t.Errorf("WithPathIndex() failed, got path: %v", ctx3.Path)
		}

		pathStr := ctx3.GetPath()
		expected := "test.items[5]"
		if pathStr != expected {
			t.Errorf("GetPath() = %s, want %s", pathStr, expected)
		}
	})

	t.Run("ValidationErrors collection", func(t *testing.T) {
		errors := &ValidationErrors{}

		// Test adding errors
		errors.AddError(ErrorTypeRequired, "TestObject", "TestField", "test message", []string{"path", "to", "field"})

		if !errors.HasErrors() {
			t.Error("HasErrors() should return true after adding error")
		}

		if len(errors.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errors.Errors))
		}

		err := errors.Errors[0]
		if err.Type != ErrorTypeRequired {
			t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeRequired)
		}
		if err.Object != "TestObject" {
			t.Errorf("Error object = %s, want TestObject", err.Object)
		}
		if err.Field != "TestField" {
			t.Errorf("Error field = %s, want TestField", err.Field)
		}

		// Test error message formatting
		errMsg := err.Error()
		if !strings.Contains(errMsg, "[Required]") {
			t.Errorf("Error message should contain [Required], got: %s", errMsg)
		}
		if !strings.Contains(errMsg, "TestObject.TestField") {
			t.Errorf("Error message should contain object.field, got: %s", errMsg)
		}
		if !strings.Contains(errMsg, "path.to.field") {
			t.Errorf("Error message should contain path, got: %s", errMsg)
		}
	})

	t.Run("ValidationHelper required field validation", func(t *testing.T) {
		helper := NewValidationHelper()
		ctx := NewValidationContext()
		errors := &ValidationErrors{}

		// Test empty string validation
		helper.ValidateRequired("", "TestField", "TestObject", ctx, errors)

		if !errors.HasErrors() {
			t.Error("Should have validation error for empty string")
		}

		if len(errors.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errors.Errors))
		}

		err := errors.Errors[0]
		if err.Type != ErrorTypeRequired {
			t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeRequired)
		}
	})

	t.Run("ValidationHelper enum validation", func(t *testing.T) {
		helper := NewValidationHelper()
		ctx := NewValidationContext()
		errors := &ValidationErrors{}

		// Test invalid enum value
		allowedValues := []string{"value1", "value2", "value3"}
		helper.ValidateEnum("invalid", "TestField", "TestObject", allowedValues, ctx, errors)

		if !errors.HasErrors() {
			t.Error("Should have validation error for invalid enum value")
		}

		err := errors.Errors[0]
		if err.Type != ErrorTypeInvalid {
			t.Errorf("Error type = %v, want %v", err.Type, ErrorTypeInvalid)
		}

		// Test valid enum value
		errors2 := &ValidationErrors{}
		helper.ValidateEnum("value2", "TestField", "TestObject", allowedValues, ctx, errors2)

		if errors2.HasErrors() {
			t.Error("Should not have validation error for valid enum value")
		}
	})

	t.Run("Enhanced validation with multiple errors", func(t *testing.T) {
		// Create an invalid StateMachine to test multiple error collection
		sm := &StateMachine{
			// Missing required fields ID, Name, Version
			Regions: []*Region{
				{
					// Missing required fields ID, Name
				},
			},
		}

		err := sm.Validate()
		if err == nil {
			t.Fatal("Expected validation errors for invalid StateMachine")
		}

		errMsg := err.Error()

		// Should contain multiple errors
		if !strings.Contains(errMsg, "multiple validation errors") {
			t.Errorf("Should contain 'multiple validation errors', got: %s", errMsg)
		}

		// Should contain path information
		if !strings.Contains(errMsg, "Regions[0]") {
			t.Errorf("Should contain path information 'Regions[0]', got: %s", errMsg)
		}

		// Should contain error types
		if !strings.Contains(errMsg, "[Required]") {
			t.Errorf("Should contain error type '[Required]', got: %s", errMsg)
		}
	})
}

func TestValidationContextNilHandling(t *testing.T) {
	t.Run("nil context handling", func(t *testing.T) {
		var ctx *ValidationContext

		// Test that nil context methods don't panic
		ctx2 := ctx.WithPath("test")
		if ctx2 == nil {
			t.Error("WithPath should handle nil context")
		}

		ctx3 := ctx.WithStateMachine(&StateMachine{})
		if ctx3 == nil {
			t.Error("WithStateMachine should handle nil context")
		}

		path := ctx.GetPath()
		if path != "" {
			t.Errorf("GetPath on nil context should return empty string, got: %s", path)
		}
	})
}
