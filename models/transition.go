package models

import "fmt"

// TransitionKind represents the kind of transition
type TransitionKind string

const (
	TransitionKindInternal TransitionKind = "internal"
	TransitionKindLocal    TransitionKind = "local"
	TransitionKindExternal TransitionKind = "external"
)

// IsValid checks if the TransitionKind is valid
func (tk TransitionKind) IsValid() bool {
	validKinds := map[TransitionKind]bool{
		TransitionKindInternal: true,
		TransitionKindLocal:    true,
		TransitionKindExternal: true,
	}
	return validKinds[tk]
}

// Transition represents a transition between vertices in a state machine
type Transition struct {
	ID       string         `json:"id" validate:"required"`
	Name     string         `json:"name,omitempty"`
	Source   *Vertex        `json:"source" validate:"required"`
	Target   *Vertex        `json:"target" validate:"required"`
	Kind     TransitionKind `json:"kind" validate:"required"`
	Triggers []*Trigger     `json:"triggers,omitempty"`
	Guard    *Constraint    `json:"guard,omitempty"`
	Effect   *Behavior      `json:"effect,omitempty"`
	// Container *Region       `json:"-"` // Parent region (not serialized)
}

// Validate validates the Transition data integrity
func (t *Transition) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	t.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Transition with the provided context
func (t *Transition) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	t.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Transition and collects all errors
func (t *Transition) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(t.ID, "ID", "Transition", context, errors)

	// Validate required references
	helper.ValidateReference(t.Source, "Source", "Transition", context, errors, true)
	helper.ValidateReference(t.Target, "Target", "Transition", context, errors, true)

	// Validate kind
	if !t.Kind.IsValid() {
		errors.AddError(
			ErrorTypeInvalid,
			"Transition",
			"Kind",
			fmt.Sprintf("invalid TransitionKind: %s", t.Kind),
			context.Path,
		)
	}

	// Validate triggers collection
	triggerValidators := make([]Validator, len(t.Triggers))
	for i, trigger := range t.Triggers {
		triggerValidators[i] = trigger
	}
	helper.ValidateCollection(triggerValidators, "Triggers", "Transition", context, errors)

	// Validate optional references
	helper.ValidateReference(t.Guard, "Guard", "Transition", context, errors, false)
	helper.ValidateReference(t.Effect, "Effect", "Transition", context, errors, false)
}
