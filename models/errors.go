package models

import (
	"fmt"
	"strings"
)

// ValidationErrorType represents the type of validation error
type ValidationErrorType int

const (
	ErrorTypeRequired ValidationErrorType = iota
	ErrorTypeInvalid
	ErrorTypeConstraint
	ErrorTypeReference
	ErrorTypeMultiplicity
)

// String returns the string representation of ValidationErrorType
func (vet ValidationErrorType) String() string {
	switch vet {
	case ErrorTypeRequired:
		return "Required"
	case ErrorTypeInvalid:
		return "Invalid"
	case ErrorTypeConstraint:
		return "Constraint"
	case ErrorTypeReference:
		return "Reference"
	case ErrorTypeMultiplicity:
		return "Multiplicity"
	default:
		return "Unknown"
	}
}

// ValidationError represents a validation error with enhanced context
type ValidationError struct {
	Type    ValidationErrorType    `json:"type"`
	Object  string                 `json:"object"`
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Path    []string               `json:"path"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// Error implements the error interface
func (ve *ValidationError) Error() string {
	pathStr := ""
	if len(ve.Path) > 0 {
		pathStr = fmt.Sprintf(" at %s", strings.Join(ve.Path, "."))
	}
	return fmt.Sprintf("[%s] %s.%s: %s%s", ve.Type.String(), ve.Object, ve.Field, ve.Message, pathStr)
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors struct {
	Errors []*ValidationError `json:"errors"`
}

// Error implements the error interface for ValidationErrors
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "no validation errors"
	}
	if len(ve.Errors) == 1 {
		return ve.Errors[0].Error()
	}

	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("multiple validation errors:\n  - %s", strings.Join(messages, "\n  - "))
}

// Add adds a validation error to the collection
func (ve *ValidationErrors) Add(err *ValidationError) {
	ve.Errors = append(ve.Errors, err)
}

// AddError adds a simple error as a validation error
func (ve *ValidationErrors) AddError(errorType ValidationErrorType, object, field, message string, path []string) {
	ve.Add(&ValidationError{
		Type:    errorType,
		Object:  object,
		Field:   field,
		Message: message,
		Path:    path,
	})
}

// HasErrors returns true if there are any validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// ToError returns the ValidationErrors as an error if there are any errors, nil otherwise
func (ve *ValidationErrors) ToError() error {
	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidationContext provides context for validation operations
type ValidationContext struct {
	StateMachine *StateMachine          `json:"state_machine,omitempty"`
	Region       *Region                `json:"region,omitempty"`
	Parent       interface{}            `json:"-"` // Parent object (not serialized due to potential cycles)
	Path         []string               `json:"path"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// NewValidationContext creates a new validation context
func NewValidationContext() *ValidationContext {
	return &ValidationContext{
		Path:     make([]string, 0),
		Metadata: make(map[string]interface{}),
	}
}

// WithStateMachine returns a new context with the specified state machine
func (vc *ValidationContext) WithStateMachine(sm *StateMachine) *ValidationContext {
	if vc == nil {
		vc = NewValidationContext()
	}
	newCtx := *vc
	newCtx.StateMachine = sm
	return &newCtx
}

// WithRegion returns a new context with the specified region
func (vc *ValidationContext) WithRegion(region *Region) *ValidationContext {
	if vc == nil {
		vc = NewValidationContext()
	}
	newCtx := *vc
	newCtx.Region = region
	return &newCtx
}

// WithParent returns a new context with the specified parent
func (vc *ValidationContext) WithParent(parent interface{}) *ValidationContext {
	if vc == nil {
		vc = NewValidationContext()
	}
	newCtx := *vc
	newCtx.Parent = parent
	return &newCtx
}

// WithPath returns a new context with the specified path element added
func (vc *ValidationContext) WithPath(pathElement string) *ValidationContext {
	if vc == nil {
		vc = NewValidationContext()
	}
	newCtx := *vc
	newCtx.Path = make([]string, len(vc.Path)+1)
	copy(newCtx.Path, vc.Path)
	newCtx.Path[len(vc.Path)] = pathElement
	return &newCtx
}

// WithPathIndex returns a new context with an indexed path element added
func (vc *ValidationContext) WithPathIndex(pathElement string, index int) *ValidationContext {
	return vc.WithPath(fmt.Sprintf("%s[%d]", pathElement, index))
}

// GetPath returns the current validation path as a string
func (vc *ValidationContext) GetPath() string {
	if vc == nil || vc.Path == nil {
		return ""
	}
	return strings.Join(vc.Path, ".")
}

// SetMetadata sets a metadata value in the context
func (vc *ValidationContext) SetMetadata(key string, value interface{}) {
	if vc.Metadata == nil {
		vc.Metadata = make(map[string]interface{})
	}
	vc.Metadata[key] = value
}

// GetMetadata gets a metadata value from the context
func (vc *ValidationContext) GetMetadata(key string) (interface{}, bool) {
	if vc.Metadata == nil {
		return nil, false
	}
	value, exists := vc.Metadata[key]
	return value, exists
}
