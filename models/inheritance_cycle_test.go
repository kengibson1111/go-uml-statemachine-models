package models

import (
	"strings"
	"testing"
)

func TestReferenceValidator_InheritanceCycleDetection_Manual(t *testing.T) {
	// Test inheritance cycle detection by manually setting up the inheritance tree
	rv := NewReferenceValidator()

	// Create objects with different IDs
	sm1 := &StateMachine{ID: "sm1", Name: "StateMachine1", Version: "1.0"}
	sm2 := &StateMachine{ID: "sm2", Name: "StateMachine2", Version: "1.0"}
	sm3 := &StateMachine{ID: "sm3", Name: "StateMachine3", Version: "1.0"}

	// Set up reference map
	rv.referenceMap = map[string]interface{}{
		"sm1": sm1,
		"sm2": sm2,
		"sm3": sm3,
	}

	// Set up inheritance tree with a cycle: sm1 -> sm2 -> sm3 -> sm1
	rv.inheritanceTree = map[string]string{
		"sm1": "sm2",
		"sm2": "sm3",
		"sm3": "sm1", // Creates cycle
	}

	rv.errors = &ValidationErrors{}
	rv.context = NewValidationContext()

	// Test inheritance cycle detection
	rv.validateInheritanceRelationships()

	if !rv.errors.HasErrors() {
		t.Errorf("Expected inheritance cycle error but got none")
		return
	}

	errStr := rv.errors.Error()
	if !strings.Contains(errStr, "inheritance cycle detected") {
		t.Errorf("Expected inheritance cycle error, got: %v", errStr)
	}
}

func TestReferenceValidator_DirectInheritanceCycle(t *testing.T) {
	// Test direct self-reference in inheritance
	rv := NewReferenceValidator()

	// Create object
	sm1 := &StateMachine{ID: "sm1", Name: "StateMachine1", Version: "1.0"}

	// Set up reference map
	rv.referenceMap = map[string]interface{}{
		"sm1": sm1,
	}

	// Set up inheritance tree with direct self-reference
	rv.inheritanceTree = map[string]string{
		"sm1": "sm1", // Direct self-reference
	}

	rv.errors = &ValidationErrors{}
	rv.context = NewValidationContext()

	// Test inheritance cycle detection
	rv.validateInheritanceRelationships()

	if !rv.errors.HasErrors() {
		t.Errorf("Expected inheritance cycle error but got none")
		return
	}

	errStr := rv.errors.Error()
	if !strings.Contains(errStr, "inheritance cycle detected") && !strings.Contains(errStr, "direct self-reference") {
		t.Errorf("Expected inheritance cycle error, got: %v", errStr)
	}
}
