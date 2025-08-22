package models

// Constraint represents a constraint (guard condition)
type Constraint struct {
	ID            string `json:"id" validate:"required"`
	Name          string `json:"name,omitempty"`
	Specification string `json:"specification" validate:"required"`
	Language      string `json:"language,omitempty"`
}

// Validate validates the Constraint data integrity
func (c *Constraint) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	c.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Constraint with the provided context
func (c *Constraint) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	c.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Constraint and collects all errors
func (c *Constraint) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(c.ID, "ID", "Constraint", context, errors)
	helper.ValidateRequired(c.Specification, "Specification", "Constraint", context, errors)
}

// Behavior represents a behavior (action/activity)
type Behavior struct {
	ID            string `json:"id" validate:"required"`
	Name          string `json:"name,omitempty"`
	Specification string `json:"specification" validate:"required"`
	Language      string `json:"language,omitempty"`
}

// Validate validates the Behavior data integrity
func (b *Behavior) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	b.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Behavior with the provided context
func (b *Behavior) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	b.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Behavior and collects all errors
func (b *Behavior) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(b.ID, "ID", "Behavior", context, errors)
	helper.ValidateRequired(b.Specification, "Specification", "Behavior", context, errors)
}

// Effect is an alias for Behavior to maintain semantic clarity
type Effect = Behavior
