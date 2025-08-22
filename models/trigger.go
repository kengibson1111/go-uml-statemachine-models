package models

import "fmt"

// EventType represents the type of event
type EventType string

const (
	EventTypeCall       EventType = "call"
	EventTypeSignal     EventType = "signal"
	EventTypeChange     EventType = "change"
	EventTypeTime       EventType = "time"
	EventTypeAnyReceive EventType = "anyReceive"
)

// IsValid checks if the EventType is valid
func (et EventType) IsValid() bool {
	validTypes := map[EventType]bool{
		EventTypeCall:       true,
		EventTypeSignal:     true,
		EventTypeChange:     true,
		EventTypeTime:       true,
		EventTypeAnyReceive: true,
	}
	return validTypes[et]
}

// Event represents an event that can trigger a transition
type Event struct {
	ID   string    `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required"`
	Type EventType `json:"type" validate:"required"`
}

// Validate validates the Event data integrity
func (e *Event) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	e.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Event with the provided context
func (e *Event) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	e.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Event and collects all errors
func (e *Event) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(e.ID, "ID", "Event", context, errors)
	helper.ValidateRequired(e.Name, "Name", "Event", context, errors)

	// Validate type
	if !e.Type.IsValid() {
		errors.AddError(
			ErrorTypeInvalid,
			"Event",
			"Type",
			fmt.Sprintf("invalid EventType: %s", e.Type),
			context.Path,
		)
	}
}

// Trigger represents a trigger for a transition
type Trigger struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Event *Event `json:"event" validate:"required"`
}

// Validate validates the Trigger data integrity
func (tr *Trigger) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	tr.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Trigger with the provided context
func (tr *Trigger) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	tr.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Trigger and collects all errors
func (tr *Trigger) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(tr.ID, "ID", "Trigger", context, errors)
	helper.ValidateRequired(tr.Name, "Name", "Trigger", context, errors)

	// Validate required reference
	helper.ValidateReference(tr.Event, "Event", "Trigger", context, errors, true)
}
