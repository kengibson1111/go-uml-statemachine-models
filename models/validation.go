package models

import "fmt"

// Validator interface defines the contract for objects that can be validated
type Validator interface {
	// Validate performs basic validation without context
	Validate() error

	// ValidateInContext performs validation with the provided context
	ValidateInContext(context *ValidationContext) error
}

// ContextualValidator interface for objects that require context for validation
type ContextualValidator interface {
	// ValidateInContext performs validation with the provided context
	ValidateInContext(context *ValidationContext) error
}

// ValidatorWithErrors interface for objects that can collect multiple validation errors
type ValidatorWithErrors interface {
	// ValidateWithErrors performs validation and collects all errors
	ValidateWithErrors(context *ValidationContext, errors *ValidationErrors)
}

// ValidationHelper provides utility functions for common validation patterns
type ValidationHelper struct{}

// NewValidationHelper creates a new validation helper
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{}
}

// ValidateRequired checks if a required field is present and non-empty
func (vh *ValidationHelper) ValidateRequired(value, fieldName, objectName string, context *ValidationContext, errors *ValidationErrors) {
	if value == "" {
		errors.AddError(
			ErrorTypeRequired,
			objectName,
			fieldName,
			"field is required and cannot be empty",
			context.Path,
		)
	}
}

// ValidateRequiredPointer checks if a required pointer field is not nil
func (vh *ValidationHelper) ValidateRequiredPointer(value interface{}, fieldName, objectName string, context *ValidationContext, errors *ValidationErrors) {
	if value == nil {
		errors.AddError(
			ErrorTypeRequired,
			objectName,
			fieldName,
			"field is required and cannot be nil",
			context.Path,
		)
	}
}

// ValidateEnum checks if a value is within a set of allowed values
func (vh *ValidationHelper) ValidateEnum(value, fieldName, objectName string, allowedValues []string, context *ValidationContext, errors *ValidationErrors) {
	for _, allowed := range allowedValues {
		if value == allowed {
			return
		}
	}

	errors.AddError(
		ErrorTypeInvalid,
		objectName,
		fieldName,
		"invalid value: must be one of "+formatStringSlice(allowedValues),
		context.Path,
	)
}

// ValidateCollection validates a collection of validators
func (vh *ValidationHelper) ValidateCollection(validators []Validator, collectionName, objectName string, context *ValidationContext, errors *ValidationErrors) {
	for i, validator := range validators {
		if validator == nil {
			errors.AddError(
				ErrorTypeReference,
				objectName,
				collectionName,
				"collection contains nil element",
				context.WithPathIndex(collectionName, i).Path,
			)
			continue
		}

		// Use ValidateWithErrors if available, otherwise fall back to ValidateInContext
		if validatorWithErrors, ok := validator.(ValidatorWithErrors); ok {
			validatorWithErrors.ValidateWithErrors(context.WithPathIndex(collectionName, i), errors)
		} else if contextualValidator, ok := validator.(ContextualValidator); ok {
			if err := contextualValidator.ValidateInContext(context.WithPathIndex(collectionName, i)); err != nil {
				// Convert single error to ValidationError, preserving error type if it's a ValidationErrors
				if validationErrors, ok := err.(*ValidationErrors); ok {
					// Add all errors from the nested validation
					for _, nestedError := range validationErrors.Errors {
						errors.Add(nestedError)
					}
				} else {
					// Single error case
					errors.AddError(
						ErrorTypeInvalid,
						objectName,
						collectionName,
						err.Error(),
						context.WithPathIndex(collectionName, i).Path,
					)
				}
			}
		} else {
			// Fall back to basic Validate method
			if err := validator.Validate(); err != nil {
				// Convert single error to ValidationError, preserving error type if it's a ValidationErrors
				if validationErrors, ok := err.(*ValidationErrors); ok {
					// Add all errors from the nested validation
					for _, nestedError := range validationErrors.Errors {
						errors.Add(nestedError)
					}
				} else {
					// Single error case
					errors.AddError(
						ErrorTypeInvalid,
						objectName,
						collectionName,
						err.Error(),
						context.WithPathIndex(collectionName, i).Path,
					)
				}
			}
		}
	}
}

// ValidateReference validates a single reference to another validator
func (vh *ValidationHelper) ValidateReference(validator Validator, fieldName, objectName string, context *ValidationContext, errors *ValidationErrors, required bool) {
	// Check for nil interface or nil underlying value
	if validator == nil || isNilInterface(validator) {
		if required {
			errors.AddError(
				ErrorTypeRequired,
				objectName,
				fieldName,
				"required reference cannot be nil",
				context.Path,
			)
		}
		return
	}

	// Create a new context for the referenced object
	refContext := context.WithPath(fieldName)

	// Use ValidateWithErrors if available, otherwise fall back to ValidateInContext
	if validatorWithErrors, ok := validator.(ValidatorWithErrors); ok {
		validatorWithErrors.ValidateWithErrors(refContext, errors)
	} else if contextualValidator, ok := validator.(ContextualValidator); ok {
		if err := contextualValidator.ValidateInContext(refContext); err != nil {
			// Convert error to ValidationError, preserving error type if it's a ValidationErrors
			if validationErrors, ok := err.(*ValidationErrors); ok {
				// Add all errors from the nested validation
				for _, nestedError := range validationErrors.Errors {
					errors.Add(nestedError)
				}
			} else {
				// Single error case
				errors.AddError(
					ErrorTypeReference,
					objectName,
					fieldName,
					err.Error(),
					refContext.Path,
				)
			}
		}
	} else {
		// Fall back to basic Validate method
		if err := validator.Validate(); err != nil {
			// Convert error to ValidationError, preserving error type if it's a ValidationErrors
			if validationErrors, ok := err.(*ValidationErrors); ok {
				// Add all errors from the nested validation
				for _, nestedError := range validationErrors.Errors {
					errors.Add(nestedError)
				}
			} else {
				// Single error case
				errors.AddError(
					ErrorTypeReference,
					objectName,
					fieldName,
					err.Error(),
					refContext.Path,
				)
			}
		}
	}
}

// isNilInterface checks if an interface contains a nil value
func isNilInterface(i interface{}) bool {
	if i == nil {
		return true
	}
	// Use type assertions to check if the underlying value is nil
	// This handles cases like (*SomeType)(nil) which are not == nil
	switch v := i.(type) {
	case *Behavior:
		return v == nil
	case *Constraint:
		return v == nil
	case *StateMachine:
		return v == nil
	case *Region:
		return v == nil
	case *State:
		return v == nil
	case *Pseudostate:
		return v == nil
	case *FinalState:
		return v == nil
	case *Vertex:
		return v == nil
	case *Transition:
		return v == nil
	case *Trigger:
		return v == nil
	case *Event:
		return v == nil
	case *ConnectionPointReference:
		return v == nil
	default:
		return false
	}
}

// ValidateUniqueIDs validates that all objects in a collection have unique IDs
func (vh *ValidationHelper) ValidateUniqueIDs(objects []interface{}, collectionName, objectName string, context *ValidationContext, errors *ValidationErrors, getID func(interface{}) string) {
	idMap := make(map[string]int)

	for i, obj := range objects {
		if obj == nil {
			continue
		}

		id := getID(obj)
		if id == "" {
			continue // Empty IDs will be caught by required field validation
		}

		if prevIndex, exists := idMap[id]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				objectName,
				collectionName,
				fmt.Sprintf("duplicate ID '%s' found at indices %d and %d", id, prevIndex, i),
				context.WithPathIndex(collectionName, i).Path,
			)
		} else {
			idMap[id] = i
		}
	}
}

// ValidateUniqueNames validates that all objects in a collection have unique names
func (vh *ValidationHelper) ValidateUniqueNames(objects []interface{}, collectionName, objectName string, context *ValidationContext, errors *ValidationErrors, getName func(interface{}) string) {
	nameMap := make(map[string]int)

	for i, obj := range objects {
		if obj == nil {
			continue
		}

		name := getName(obj)
		if name == "" {
			continue // Empty names are allowed in some cases
		}

		if prevIndex, exists := nameMap[name]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				objectName,
				collectionName,
				fmt.Sprintf("duplicate name '%s' found at indices %d and %d (may cause confusion)", name, prevIndex, i),
				context.WithPathIndex(collectionName, i).Path,
			)
		} else {
			nameMap[name] = i
		}
	}
}

// ValidateConditionalRequired validates that a field is required under certain conditions
func (vh *ValidationHelper) ValidateConditionalRequired(value, fieldName, objectName string, condition bool, conditionDescription string, context *ValidationContext, errors *ValidationErrors) {
	if condition && value == "" {
		errors.AddError(
			ErrorTypeRequired,
			objectName,
			fieldName,
			fmt.Sprintf("field is required when %s", conditionDescription),
			context.Path,
		)
	}
}

// ValidateConditionalRequiredPointer validates that a pointer field is required under certain conditions
func (vh *ValidationHelper) ValidateConditionalRequiredPointer(value interface{}, fieldName, objectName string, condition bool, conditionDescription string, context *ValidationContext, errors *ValidationErrors) {
	if condition && value == nil {
		errors.AddError(
			ErrorTypeRequired,
			objectName,
			fieldName,
			fmt.Sprintf("field is required when %s", conditionDescription),
			context.Path,
		)
	}
}

// ValidateMutuallyExclusive validates that only one of the specified fields is set
func (vh *ValidationHelper) ValidateMutuallyExclusive(values map[string]interface{}, objectName string, context *ValidationContext, errors *ValidationErrors) {
	setFields := make([]string, 0)

	for fieldName, value := range values {
		if value != nil && !isEmptyValue(value) {
			setFields = append(setFields, fieldName)
		}
	}

	if len(setFields) > 1 {
		errors.AddError(
			ErrorTypeConstraint,
			objectName,
			"MutuallyExclusive",
			fmt.Sprintf("fields %v are mutually exclusive, but multiple are set", setFields),
			context.Path,
		)
	}
}

// ValidateAtLeastOne validates that at least one of the specified fields is set
func (vh *ValidationHelper) ValidateAtLeastOne(values map[string]interface{}, objectName string, context *ValidationContext, errors *ValidationErrors) {
	hasValue := false
	fieldNames := make([]string, 0, len(values))

	for fieldName, value := range values {
		fieldNames = append(fieldNames, fieldName)
		if value != nil && !isEmptyValue(value) {
			hasValue = true
			break
		}
	}

	if !hasValue {
		errors.AddError(
			ErrorTypeConstraint,
			objectName,
			"AtLeastOne",
			fmt.Sprintf("at least one of the following fields must be set: %v", fieldNames),
			context.Path,
		)
	}
}

// ValidateStringLength validates string length constraints
func (vh *ValidationHelper) ValidateStringLength(value, fieldName, objectName string, minLength, maxLength int, context *ValidationContext, errors *ValidationErrors) {
	length := len(value)

	if minLength > 0 && length < minLength {
		errors.AddError(
			ErrorTypeConstraint,
			objectName,
			fieldName,
			fmt.Sprintf("field must be at least %d characters long, got %d", minLength, length),
			context.Path,
		)
	}

	if maxLength > 0 && length > maxLength {
		errors.AddError(
			ErrorTypeConstraint,
			objectName,
			fieldName,
			fmt.Sprintf("field must be at most %d characters long, got %d", maxLength, length),
			context.Path,
		)
	}
}

// ValidateCollectionSize validates collection size constraints
func (vh *ValidationHelper) ValidateCollectionSize(collection interface{}, collectionName, objectName string, minSize, maxSize int, context *ValidationContext, errors *ValidationErrors) {
	var size int

	// Use reflection to get the size of different collection types
	switch v := collection.(type) {
	case []interface{}:
		size = len(v)
	case []string:
		size = len(v)
	case []Validator:
		size = len(v)
	case []*Region:
		size = len(v)
	case []*State:
		size = len(v)
	case []*Transition:
		size = len(v)
	case []*Vertex:
		size = len(v)
	case []*Trigger:
		size = len(v)
	case []*Pseudostate:
		size = len(v)
	default:
		// For unknown types, skip validation
		return
	}

	if minSize > 0 && size < minSize {
		errors.AddError(
			ErrorTypeMultiplicity,
			objectName,
			collectionName,
			fmt.Sprintf("collection must have at least %d elements, got %d", minSize, size),
			context.Path,
		)
	}

	if maxSize > 0 && size > maxSize {
		errors.AddError(
			ErrorTypeMultiplicity,
			objectName,
			collectionName,
			fmt.Sprintf("collection must have at most %d elements, got %d", maxSize, size),
			context.Path,
		)
	}
}

// isEmptyValue checks if a value is considered empty
func isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	default:
		return false
	}
}

// formatStringSlice formats a slice of strings for display in error messages
func formatStringSlice(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	if len(values) == 1 {
		return "[" + values[0] + "]"
	}

	result := "["
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		result += value
	}
	result += "]"
	return result
}
