package models

import "fmt"

// Vertex represents a vertex in a state machine (base type for states and pseudostates)
type Vertex struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"` // "state", "pseudostate", "finalstate"
	// Container *Region `json:"-"` // Parent region (not serialized)
}

// Validate validates the Vertex data integrity
func (v *Vertex) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	v.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Vertex with the provided context
func (v *Vertex) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	v.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Vertex and collects all errors
func (v *Vertex) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate required fields
	helper.ValidateRequired(v.ID, "ID", "Vertex", context, errors)
	helper.ValidateRequired(v.Name, "Name", "Vertex", context, errors)
	helper.ValidateRequired(v.Type, "Type", "Vertex", context, errors)

	// Validate type is one of the allowed values
	validTypes := []string{"state", "pseudostate", "finalstate"}
	helper.ValidateEnum(v.Type, "Type", "Vertex", validTypes, context, errors)
}

// State represents a state in a state machine
type State struct {
	Vertex                                        // Embedded vertex
	IsComposite       bool                        `json:"is_composite"`
	IsOrthogonal      bool                        `json:"is_orthogonal"`
	IsSimple          bool                        `json:"is_simple"`
	IsSubmachineState bool                        `json:"is_submachine_state"`
	Regions           []*Region                   `json:"regions,omitempty"`
	Entry             *Behavior                   `json:"entry,omitempty"`
	Exit              *Behavior                   `json:"exit,omitempty"`
	DoActivity        *Behavior                   `json:"do_activity,omitempty"`
	Submachine        *StateMachine               `json:"submachine,omitempty"`
	Connections       []*ConnectionPointReference `json:"connections,omitempty"`
}

// Validate validates the State data integrity
func (s *State) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	s.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the State with the provided context
func (s *State) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	s.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the State and collects all errors
func (s *State) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate embedded vertex
	s.Vertex.ValidateWithErrors(context.WithPath("Vertex"), errors)

	// Validate that type is "state"
	if s.Type != "state" {
		errors.AddError(
			ErrorTypeConstraint,
			"State",
			"Type",
			fmt.Sprintf("State must have type 'state', got: %s", s.Type),
			context.Path,
		)
	}

	// Validate regions if composite
	if s.IsComposite {
		regionValidators := make([]Validator, len(s.Regions))
		for i, region := range s.Regions {
			regionValidators[i] = region
		}
		helper.ValidateCollection(regionValidators, "Regions", "State", context, errors)
	}

	// Validate behaviors
	helper.ValidateReference(s.Entry, "Entry", "State", context, errors, false)
	helper.ValidateReference(s.Exit, "Exit", "State", context, errors, false)
	helper.ValidateReference(s.DoActivity, "DoActivity", "State", context, errors, false)

	// Validate submachine if present
	helper.ValidateReference(s.Submachine, "Submachine", "State", context, errors, false)

	// Validate connections
	connectionValidators := make([]Validator, len(s.Connections))
	for i, conn := range s.Connections {
		connectionValidators[i] = conn
	}
	helper.ValidateCollection(connectionValidators, "Connections", "State", context, errors)
}

// PseudostateKind represents the kind of pseudostate
type PseudostateKind string

const (
	PseudostateKindInitial        PseudostateKind = "initial"
	PseudostateKindDeepHistory    PseudostateKind = "deepHistory"
	PseudostateKindShallowHistory PseudostateKind = "shallowHistory"
	PseudostateKindJoin           PseudostateKind = "join"
	PseudostateKindFork           PseudostateKind = "fork"
	PseudostateKindJunction       PseudostateKind = "junction"
	PseudostateKindChoice         PseudostateKind = "choice"
	PseudostateKindEntryPoint     PseudostateKind = "entryPoint"
	PseudostateKindExitPoint      PseudostateKind = "exitPoint"
	PseudostateKindTerminate      PseudostateKind = "terminate"
)

// IsValid checks if the PseudostateKind is valid
func (pk PseudostateKind) IsValid() bool {
	validKinds := map[PseudostateKind]bool{
		PseudostateKindInitial:        true,
		PseudostateKindDeepHistory:    true,
		PseudostateKindShallowHistory: true,
		PseudostateKindJoin:           true,
		PseudostateKindFork:           true,
		PseudostateKindJunction:       true,
		PseudostateKindChoice:         true,
		PseudostateKindEntryPoint:     true,
		PseudostateKindExitPoint:      true,
		PseudostateKindTerminate:      true,
	}
	return validKinds[pk]
}

// Pseudostate represents a pseudostate in a state machine
type Pseudostate struct {
	Vertex                 // Embedded vertex
	Kind   PseudostateKind `json:"kind" validate:"required"`
}

// Validate validates the Pseudostate data integrity
func (ps *Pseudostate) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	ps.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the Pseudostate with the provided context
func (ps *Pseudostate) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	ps.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the Pseudostate and collects all errors
func (ps *Pseudostate) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	// Validate embedded vertex
	ps.Vertex.ValidateWithErrors(context.WithPath("Vertex"), errors)

	// Validate that type is "pseudostate"
	if ps.Type != "pseudostate" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Type",
			fmt.Sprintf("Pseudostate must have type 'pseudostate', got: %s", ps.Type),
			context.Path,
		)
	}

	// Validate kind
	if !ps.Kind.IsValid() {
		errors.AddError(
			ErrorTypeInvalid,
			"Pseudostate",
			"Kind",
			fmt.Sprintf("invalid PseudostateKind: %s", ps.Kind),
			context.Path,
		)
	}
}

// FinalState represents a final state in a state machine
type FinalState struct {
	Vertex // Embedded vertex
}

// Validate validates the FinalState data integrity
func (fs *FinalState) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	fs.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the FinalState with the provided context
func (fs *FinalState) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	fs.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the FinalState and collects all errors
func (fs *FinalState) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	// Validate embedded vertex
	fs.Vertex.ValidateWithErrors(context.WithPath("Vertex"), errors)

	// Validate that type is "finalstate"
	if fs.Type != "finalstate" {
		errors.AddError(
			ErrorTypeConstraint,
			"FinalState",
			"Type",
			fmt.Sprintf("FinalState must have type 'finalstate', got: %s", fs.Type),
			context.Path,
		)
	}
}

// ConnectionPointReference represents a connection point reference
type ConnectionPointReference struct {
	Vertex                // Embedded vertex
	Entry  []*Pseudostate `json:"entry,omitempty"`
	Exit   []*Pseudostate `json:"exit,omitempty"`
	// State  *State         `json:"-"` // Parent state (not serialized)
}

// Validate validates the ConnectionPointReference data integrity
func (cpr *ConnectionPointReference) Validate() error {
	context := NewValidationContext()
	errors := &ValidationErrors{}
	cpr.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateInContext validates the ConnectionPointReference with the provided context
func (cpr *ConnectionPointReference) ValidateInContext(context *ValidationContext) error {
	errors := &ValidationErrors{}
	cpr.ValidateWithErrors(context, errors)
	return errors.ToError()
}

// ValidateWithErrors validates the ConnectionPointReference and collects all errors
func (cpr *ConnectionPointReference) ValidateWithErrors(context *ValidationContext, errors *ValidationErrors) {
	if context == nil {
		context = NewValidationContext()
	}
	if errors == nil {
		return
	}

	helper := NewValidationHelper()

	// Validate embedded vertex
	cpr.Vertex.ValidateWithErrors(context.WithPath("Vertex"), errors)

	// Validate entry pseudostates
	entryValidators := make([]Validator, len(cpr.Entry))
	for i, entry := range cpr.Entry {
		entryValidators[i] = entry
	}
	helper.ValidateCollection(entryValidators, "Entry", "ConnectionPointReference", context, errors)

	// Validate exit pseudostates
	exitValidators := make([]Validator, len(cpr.Exit))
	for i, exit := range cpr.Exit {
		exitValidators[i] = exit
	}
	helper.ValidateCollection(exitValidators, "Exit", "ConnectionPointReference", context, errors)
}
