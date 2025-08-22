package models

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
				// Convert single error to ValidationError
				errors.AddError(
					ErrorTypeInvalid,
					objectName,
					collectionName,
					err.Error(),
					context.WithPathIndex(collectionName, i).Path,
				)
			}
		} else {
			// Fall back to basic Validate method
			if err := validator.Validate(); err != nil {
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
			errors.AddError(
				ErrorTypeReference,
				objectName,
				fieldName,
				err.Error(),
				refContext.Path,
			)
		}
	} else {
		// Fall back to basic Validate method
		if err := validator.Validate(); err != nil {
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
