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

// AddErrorWithContext adds a validation error with additional context information
func (ve *ValidationErrors) AddErrorWithContext(errorType ValidationErrorType, object, field, message string, path []string, context map[string]interface{}) {
	ve.Add(&ValidationError{
		Type:    errorType,
		Object:  object,
		Field:   field,
		Message: message,
		Path:    path,
		Context: context,
	})
}

// GetErrorsByType returns all errors of a specific type
func (ve *ValidationErrors) GetErrorsByType(errorType ValidationErrorType) []*ValidationError {
	var result []*ValidationError
	for _, err := range ve.Errors {
		if err.Type == errorType {
			result = append(result, err)
		}
	}
	return result
}

// GetErrorsByObject returns all errors for a specific object
func (ve *ValidationErrors) GetErrorsByObject(objectName string) []*ValidationError {
	var result []*ValidationError
	for _, err := range ve.Errors {
		if err.Object == objectName {
			result = append(result, err)
		}
	}
	return result
}

// GetErrorsByPath returns all errors for a specific path prefix
func (ve *ValidationErrors) GetErrorsByPath(pathPrefix string) []*ValidationError {
	var result []*ValidationError
	for _, err := range ve.Errors {
		if len(err.Path) > 0 {
			fullPath := strings.Join(err.Path, ".")
			if strings.HasPrefix(fullPath, pathPrefix) {
				result = append(result, err)
			}
		}
	}
	return result
}

// GetSummary returns a summary of errors by type
func (ve *ValidationErrors) GetSummary() map[ValidationErrorType]int {
	summary := make(map[ValidationErrorType]int)
	for _, err := range ve.Errors {
		summary[err.Type]++
	}
	return summary
}

// Merge merges another ValidationErrors into this one
func (ve *ValidationErrors) Merge(other *ValidationErrors) {
	if other != nil {
		for _, err := range other.Errors {
			ve.Add(err)
		}
	}
}

// Clear removes all errors
func (ve *ValidationErrors) Clear() {
	ve.Errors = ve.Errors[:0]
}

// Count returns the number of errors
func (ve *ValidationErrors) Count() int {
	return len(ve.Errors)
}

// IsEmpty returns true if there are no errors
func (ve *ValidationErrors) IsEmpty() bool {
	return len(ve.Errors) == 0
}

// GetDetailedReport returns a detailed report of all errors
func (ve *ValidationErrors) GetDetailedReport() string {
	if len(ve.Errors) == 0 {
		return "No validation errors"
	}

	var report strings.Builder
	report.WriteString(fmt.Sprintf("Validation Report: %d error(s) found\n", len(ve.Errors)))
	report.WriteString(strings.Repeat("=", 50) + "\n")

	// Group errors by type
	errorsByType := make(map[ValidationErrorType][]*ValidationError)
	for _, err := range ve.Errors {
		errorsByType[err.Type] = append(errorsByType[err.Type], err)
	}

	// Report errors by type
	for errorType, errors := range errorsByType {
		report.WriteString(fmt.Sprintf("\n%s Errors (%d):\n", errorType.String(), len(errors)))
		report.WriteString(strings.Repeat("-", 30) + "\n")

		for i, err := range errors {
			report.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
			if len(err.Context) > 0 {
				report.WriteString("   Context: ")
				for k, v := range err.Context {
					report.WriteString(fmt.Sprintf("%s=%v ", k, v))
				}
				report.WriteString("\n")
			}
		}
	}

	return report.String()
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

// WithMetadata returns a new context with additional metadata
func (vc *ValidationContext) WithMetadata(key string, value interface{}) *ValidationContext {
	if vc == nil {
		vc = NewValidationContext()
	}
	newCtx := *vc
	if newCtx.Metadata == nil {
		newCtx.Metadata = make(map[string]interface{})
	} else {
		// Copy the metadata map to avoid modifying the original
		newCtx.Metadata = make(map[string]interface{})
		for k, v := range vc.Metadata {
			newCtx.Metadata[k] = v
		}
	}
	newCtx.Metadata[key] = value
	return &newCtx
}

// GetFullPath returns the full path including parent context information
func (vc *ValidationContext) GetFullPath() string {
	if vc == nil {
		return ""
	}

	var pathParts []string

	// Add state machine context if available
	if vc.StateMachine != nil && vc.StateMachine.ID != "" {
		pathParts = append(pathParts, fmt.Sprintf("StateMachine[%s]", vc.StateMachine.ID))
	}

	// Add region context if available
	if vc.Region != nil && vc.Region.ID != "" {
		pathParts = append(pathParts, fmt.Sprintf("Region[%s]", vc.Region.ID))
	}

	// Add the current path
	if len(vc.Path) > 0 {
		pathParts = append(pathParts, strings.Join(vc.Path, "."))
	}

	return strings.Join(pathParts, ".")
}

// GetContextInfo returns a map of context information for debugging
func (vc *ValidationContext) GetContextInfo() map[string]interface{} {
	info := make(map[string]interface{})

	if vc == nil {
		return info
	}

	if vc.StateMachine != nil {
		info["stateMachine"] = map[string]interface{}{
			"id":   vc.StateMachine.ID,
			"name": vc.StateMachine.Name,
		}
	}

	if vc.Region != nil {
		info["region"] = map[string]interface{}{
			"id":   vc.Region.ID,
			"name": vc.Region.Name,
		}
	}

	if len(vc.Path) > 0 {
		info["path"] = vc.Path
	}

	if len(vc.Metadata) > 0 {
		info["metadata"] = vc.Metadata
	}

	return info
}

// Clone creates a deep copy of the validation context
func (vc *ValidationContext) Clone() *ValidationContext {
	if vc == nil {
		return NewValidationContext()
	}

	newCtx := &ValidationContext{
		StateMachine: vc.StateMachine,
		Region:       vc.Region,
		Parent:       vc.Parent,
		Path:         make([]string, len(vc.Path)),
		Metadata:     make(map[string]interface{}),
	}

	// Copy path
	copy(newCtx.Path, vc.Path)

	// Copy metadata
	for k, v := range vc.Metadata {
		newCtx.Metadata[k] = v
	}

	return newCtx
}
