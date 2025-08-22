package models

import (
	"fmt"
	"time"
)

// StateMachine represents a UML state machine
type StateMachine struct {
	ID               string                 `json:"id" validate:"required"`
	Name             string                 `json:"name" validate:"required"`
	Version          string                 `json:"version" validate:"required"`
	Regions          []*Region              `json:"regions"`
	ConnectionPoints []*Pseudostate         `json:"connection_points,omitempty"` // UML connection points (entry/exit pseudostates)
	IsMethod         bool                   `json:"is_method"`                   // True if this state machine is used as a method
	Entities         map[string]string      `json:"entities"`                    // entityID -> cache key mapping
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"created_at"`
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

	// Validate connection points collection
	connectionPointValidators := make([]Validator, len(sm.ConnectionPoints))
	for i, cp := range sm.ConnectionPoints {
		connectionPointValidators[i] = cp
	}
	helper.ValidateCollection(connectionPointValidators, "ConnectionPoints", "StateMachine", context, errors)

	// UML constraint validations
	sm.validateConnectionPoints(context, errors)
	sm.validateRegionMultiplicity(context, errors)
	sm.validateMethodConstraints(context, errors)
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

// validateConnectionPoints ensures connection points are entry/exit pseudostates
// UML Constraint: StateMachine connection points must be entry point or exit point pseudostates
func (sm *StateMachine) validateConnectionPoints(context *ValidationContext, errors *ValidationErrors) {
	for i, cp := range sm.ConnectionPoints {
		if cp == nil {
			continue // This will be caught by the collection validation
		}

		// Connection points must be entry point or exit point pseudostates
		if cp.Kind != PseudostateKindEntryPoint && cp.Kind != PseudostateKindExitPoint {
			errors.AddError(
				ErrorTypeConstraint,
				"StateMachine",
				"ConnectionPoints",
				fmt.Sprintf("connection point at index %d must be an entry point or exit point pseudostate, got: %s", i, cp.Kind),
				context.WithPathIndex("ConnectionPoints", i).Path,
			)
		}

		// Verify the pseudostate type is correct
		if cp.Type != "pseudostate" {
			errors.AddError(
				ErrorTypeConstraint,
				"StateMachine",
				"ConnectionPoints",
				fmt.Sprintf("connection point at index %d must have type 'pseudostate', got: %s", i, cp.Type),
				context.WithPathIndex("ConnectionPoints", i).Path,
			)
		}
	}
}

// validateRegionMultiplicity ensures at least one region exists
// UML Constraint: A StateMachine must have at least one region
func (sm *StateMachine) validateRegionMultiplicity(context *ValidationContext, errors *ValidationErrors) {
	if len(sm.Regions) == 0 {
		errors.AddError(
			ErrorTypeMultiplicity,
			"StateMachine",
			"Regions",
			"StateMachine must have at least one region (UML constraint)",
			context.Path,
		)
	}
}

// validateMethodConstraints checks method-specific constraints
// UML Constraint: If a StateMachine is used as a method, it cannot have connection points
func (sm *StateMachine) validateMethodConstraints(context *ValidationContext, errors *ValidationErrors) {
	if sm.IsMethod && len(sm.ConnectionPoints) > 0 {
		errors.AddError(
			ErrorTypeConstraint,
			"StateMachine",
			"ConnectionPoints",
			"StateMachine used as method cannot have connection points (UML constraint)",
			context.Path,
		)
	}
}
