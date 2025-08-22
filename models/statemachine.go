package models

import (
	"time"
)

// StateMachine represents a UML state machine
type StateMachine struct {
	ID        string                 `json:"id" validate:"required"`
	Name      string                 `json:"name" validate:"required"`
	Version   string                 `json:"version" validate:"required"`
	Regions   []*Region              `json:"regions"`
	Entities  map[string]string      `json:"entities"` // entityID -> cache key mapping
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}

// Validate validates the StateMachine data integrity
func (sm *StateMachine) Validate() error {
	context := NewValidationContext().WithStateMachine(sm)
	errors := &ValidationErrors{}
	sm.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the StateMachine with the provided context
func (sm *StateMachine) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	sm.ValidateWithErrors(context.WithStateMachine(sm), errors)
	return errors.ToError()
}

// ValidateWithErrors validates the StateMachine and collects all errors
func (sm *StateMachine) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(sm.ID, "ID", "StateMachine", context, errors)
	helper.ValidateRequired(sm.Name, "Name", "StateMachine", context, errors)
	helper.ValidateRequired(sm.Version, "Version", "StateMachine", context, errors)

	// Validate regions collection
	regionValidators := make([]Validator, len(sm.Regions))
	for i, region := range sm.Regions {
		regionValidators[i] = region
	}
	helper.ValidateCollection(regionValidators, "Regions", "StateMachine", context, errors)
}

// Region represents a region within a state machine
type Region struct {
	ID          string        `json:"id" validate:"required"`
	Name        string        `json:"name" validate:"required"`
	States      []*State      `json:"states"`
	Transitions []*Transition `json:"transitions"`
	Vertices    []*Vertex     `json:"vertices"`
}

// Validate validates the Region data integrity
func (r *Region) Validate() error {
	context := NewValidationContext().WithRegion(r)
	errors := &ValidationErrors{}
	r.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Region with the provided context
func (r *Region) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	r.ValidateWithErrors(context.WithRegion(r), errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Region and collects all errors
func (r *Region) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(r.ID, "ID", "Region", context, errors)
	helper.ValidateRequired(r.Name, "Name", "Region", context, errors)

	// Validate states collection
	stateValidators := make([]Validator, len(r.States))
	for i, state := range r.States {
		stateValidators[i] = state
	}
	helper.ValidateCollection(stateValidators, "States", "Region", context, errors)

	// Validate transitions collection
	transitionValidators := make([]Validator, len(r.Transitions))
	for i, transition := range r.Transitions {
		transitionValidators[i] = transition
	}
	helper.ValidateCollection(transitionValidators, "Transitions", "Region", context, errors)

	// Validate vertices collection
	vertexValidators := make([]Validator, len(r.Vertices))
	for i, vertex := range r.Vertices {
		vertexValidators[i] = vertex
	}
	helper.ValidateCollection(vertexValidators, "Vertices", "Region", context, errors)
}
