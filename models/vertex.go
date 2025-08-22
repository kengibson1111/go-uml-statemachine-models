package models

import (
	"fmt"
	"strings"
)

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

	// UML constraint validations
	s.validateCompositeConstraints(context, errors)
	s.validateSubmachineConstraints(context, errors)
	s.validateBehaviorConsistency(context, errors)
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

	// UML constraint validations
	ps.validateKindConstraints(context, errors)
	ps.validateMultiplicity(context, errors)
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

// validateKindConstraints validates kind-specific UML constraints for pseudostates
// UML Constraint: Each pseudostate kind has specific structural and behavioral constraints
func (ps *Pseudostate) validateKindConstraints(context *ValidationContext, errors *ValidationErrors) {
	switch ps.Kind {
	case PseudostateKindInitial:
		ps.validateInitialConstraints(context, errors)
	case PseudostateKindDeepHistory, PseudostateKindShallowHistory:
		ps.validateHistoryConstraints(context, errors)
	case PseudostateKindJoin:
		ps.validateJoinConstraints(context, errors)
	case PseudostateKindFork:
		ps.validateForkConstraints(context, errors)
	case PseudostateKindJunction:
		ps.validateJunctionConstraints(context, errors)
	case PseudostateKindChoice:
		ps.validateChoiceConstraints(context, errors)
	case PseudostateKindEntryPoint, PseudostateKindExitPoint:
		ps.validateConnectionPointConstraints(context, errors)
	case PseudostateKindTerminate:
		ps.validateTerminateConstraints(context, errors)
	}
}

// validateMultiplicity validates multiplicity constraints per pseudostate kind
// UML Constraint: Different pseudostate kinds have different multiplicity rules within regions
func (ps *Pseudostate) validateMultiplicity(context *ValidationContext, errors *ValidationErrors) {
	// Get the containing region from context
	region := context.Region
	if region == nil {
		// If no region context, we can't validate multiplicity constraints
		return
	}

	switch ps.Kind {
	case PseudostateKindInitial:
		ps.validateInitialMultiplicity(region, context, errors)
	case PseudostateKindDeepHistory, PseudostateKindShallowHistory:
		ps.validateHistoryMultiplicity(region, context, errors)
	case PseudostateKindTerminate:
		ps.validateTerminateMultiplicity(region, context, errors)
		// Other kinds don't have specific multiplicity constraints at region level
	}
}

// validateInitialConstraints validates constraints specific to initial pseudostates
func (ps *Pseudostate) validateInitialConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Initial pseudostates should not have incoming transitions (validated elsewhere)
	// Initial pseudostates must have exactly one outgoing transition (validated elsewhere)
	// Name should be appropriate for initial pseudostate
	if ps.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Name",
			"initial pseudostate should have a descriptive name (UML best practice)",
			context.Path,
		)
	}
}

// validateHistoryConstraints validates constraints for history pseudostates
func (ps *Pseudostate) validateHistoryConstraints(context *ValidationContext, errors *ValidationErrors) {
	// History pseudostates must be contained in composite states
	// This would require access to the containing state, which we validate through region context
	if context.Region == nil {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Kind",
			"history pseudostate must be contained within a region of a composite state (UML constraint)",
			context.Path,
		)
	}
}

// validateJoinConstraints validates constraints for join pseudostates
func (ps *Pseudostate) validateJoinConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Join pseudostates must have multiple incoming transitions and one outgoing transition
	// This is typically validated at the transition level, but we can add basic checks here
	if ps.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Name",
			"join pseudostate should have a descriptive name (UML best practice)",
			context.Path,
		)
	}
}

// validateForkConstraints validates constraints for fork pseudostates
func (ps *Pseudostate) validateForkConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Fork pseudostates must have one incoming transition and multiple outgoing transitions
	// This is typically validated at the transition level, but we can add basic checks here
	if ps.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Name",
			"fork pseudostate should have a descriptive name (UML best practice)",
			context.Path,
		)
	}
}

// validateJunctionConstraints validates constraints for junction pseudostates
func (ps *Pseudostate) validateJunctionConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Junction pseudostates are static conditional branches
	// They must have at least one incoming and one outgoing transition
	if ps.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Name",
			"junction pseudostate should have a descriptive name (UML best practice)",
			context.Path,
		)
	}
}

// validateChoiceConstraints validates constraints for choice pseudostates
func (ps *Pseudostate) validateChoiceConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Choice pseudostates are dynamic conditional branches
	// They must have at least one incoming and one outgoing transition
	if ps.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Name",
			"choice pseudostate should have a descriptive name (UML best practice)",
			context.Path,
		)
	}
}

// validateConnectionPointConstraints validates constraints for entry/exit point pseudostates
func (ps *Pseudostate) validateConnectionPointConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Entry and exit points are used as connection points for submachine states
	// They should only appear in connection point collections of state machines
	// However, we need to be more lenient here as they can appear in various contexts
	// The main constraint is that they should be properly used as connection points

	// Only validate if we have enough context to determine proper usage
	if context.StateMachine == nil && context.Region == nil {
		// If we have no context at all, we can't validate properly
		// This might be a standalone validation, so we'll be lenient
		return
	}

	// If we're in a state machine context but not as a connection point, that might be an issue
	// But we need to check the path to see if we're actually in the connection points collection
	pathStr := strings.Join(context.Path, ".")
	if context.StateMachine != nil && !strings.Contains(pathStr, "ConnectionPoints") && !strings.Contains(pathStr, "Connections") {
		// This is a more nuanced check - only flag if we're clearly not in a connection point context
		// and we're directly in a region (which would be inappropriate for entry/exit points)
		if context.Region != nil && strings.Contains(pathStr, "Vertices") {
			errors.AddError(
				ErrorTypeConstraint,
				"Pseudostate",
				"Kind",
				fmt.Sprintf("%s pseudostate should be used as a connection point, not as a regular vertex in a region (UML constraint)", ps.Kind),
				context.Path,
			)
		}
	}
}

// validateTerminateConstraints validates constraints for terminate pseudostates
func (ps *Pseudostate) validateTerminateConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Terminate pseudostates should not have outgoing transitions
	// This is validated at the transition level, but we can add basic checks here
	if ps.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"Pseudostate",
			"Name",
			"terminate pseudostate should have a descriptive name (UML best practice)",
			context.Path,
		)
	}
}

// validateInitialMultiplicity validates that there is at most one initial pseudostate per region
func (ps *Pseudostate) validateInitialMultiplicity(region *Region, context *ValidationContext, errors *ValidationErrors) {
	if ps.Kind != PseudostateKindInitial {
		return
	}

	initialCount := 0
	var initialIndices []int

	// Count initial pseudostates in the region's vertices
	for i, vertex := range region.Vertices {
		if vertex == nil {
			continue
		}

		// Check if this is an initial pseudostate
		if vertex.Type == "pseudostate" {
			// We need to check if this vertex represents an initial pseudostate
			// Since we don't have direct access to the Pseudostate object from Vertex,
			// we use the same logic as in the region validation
			if ps.isInitialPseudostateVertex(vertex) {
				initialCount++
				initialIndices = append(initialIndices, i)
			}
		}
	}

	// If this pseudostate is initial and there are others, report error
	if initialCount > 1 {
		errors.AddError(
			ErrorTypeMultiplicity,
			"Pseudostate",
			"Kind",
			fmt.Sprintf("region can have at most one initial pseudostate, found %d initial pseudostates (UML constraint)", initialCount),
			context.Path,
		)
	}
}

// validateHistoryMultiplicity validates history pseudostate multiplicity constraints
func (ps *Pseudostate) validateHistoryMultiplicity(region *Region, context *ValidationContext, errors *ValidationErrors) {
	if ps.Kind != PseudostateKindDeepHistory && ps.Kind != PseudostateKindShallowHistory {
		return
	}

	// Count history pseudostates of the same kind in the region
	historyCount := 0
	for _, vertex := range region.Vertices {
		if vertex == nil {
			continue
		}

		if vertex.Type == "pseudostate" {
			// We would need access to the actual Pseudostate object to check the kind
			// For now, we use naming conventions as a heuristic
			if ps.isHistoryPseudostateVertex(vertex, ps.Kind) {
				historyCount++
			}
		}
	}

	// A region should typically have at most one history pseudostate of each kind
	if historyCount > 1 {
		errors.AddError(
			ErrorTypeMultiplicity,
			"Pseudostate",
			"Kind",
			fmt.Sprintf("region should have at most one %s pseudostate, found %d (UML best practice)", ps.Kind, historyCount),
			context.Path,
		)
	}
}

// validateTerminateMultiplicity validates terminate pseudostate multiplicity constraints
func (ps *Pseudostate) validateTerminateMultiplicity(region *Region, context *ValidationContext, errors *ValidationErrors) {
	if ps.Kind != PseudostateKindTerminate {
		return
	}

	// Count terminate pseudostates in the region
	terminateCount := 0
	for _, vertex := range region.Vertices {
		if vertex == nil {
			continue
		}

		if vertex.Type == "pseudostate" {
			if ps.isTerminatePseudostateVertex(vertex) {
				terminateCount++
			}
		}
	}

	// Multiple terminate pseudostates in a region might indicate design issues
	if terminateCount > 1 {
		errors.AddError(
			ErrorTypeMultiplicity,
			"Pseudostate",
			"Kind",
			fmt.Sprintf("region has %d terminate pseudostates, consider if this is intended (UML design consideration)", terminateCount),
			context.Path,
		)
	}
}

// Helper methods for identifying pseudostate types from vertex information

// isInitialPseudostateVertex checks if a vertex represents an initial pseudostate using naming conventions
func (ps *Pseudostate) isInitialPseudostateVertex(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	// Check common naming patterns for initial pseudostates
	name := vertex.Name
	id := vertex.ID

	initialPatterns := []string{
		"initial", "Initial", "INITIAL",
		"init", "Init", "INIT",
		"start", "Start", "START",
	}

	for _, pattern := range initialPatterns {
		if name == pattern || id == pattern {
			return true
		}
	}

	return false
}

// isHistoryPseudostateVertex checks if a vertex represents a history pseudostate of the specified kind
func (ps *Pseudostate) isHistoryPseudostateVertex(vertex *Vertex, kind PseudostateKind) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	name := vertex.Name
	id := vertex.ID

	if kind == PseudostateKindDeepHistory {
		deepHistoryPatterns := []string{
			"deepHistory", "DeepHistory", "DEEP_HISTORY",
			"deep_history", "deephistory", "H*",
		}
		for _, pattern := range deepHistoryPatterns {
			if name == pattern || id == pattern {
				return true
			}
		}
	} else if kind == PseudostateKindShallowHistory {
		shallowHistoryPatterns := []string{
			"shallowHistory", "ShallowHistory", "SHALLOW_HISTORY",
			"shallow_history", "shallowhistory", "H",
		}
		for _, pattern := range shallowHistoryPatterns {
			if name == pattern || id == pattern {
				return true
			}
		}
	}

	return false
}

// isTerminatePseudostateVertex checks if a vertex represents a terminate pseudostate
func (ps *Pseudostate) isTerminatePseudostateVertex(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	name := vertex.Name
	id := vertex.ID

	terminatePatterns := []string{
		"terminate", "Terminate", "TERMINATE",
		"term", "Term", "TERM",
		"end", "End", "END",
	}

	for _, pattern := range terminatePatterns {
		if name == pattern || id == pattern {
			return true
		}
	}

	return false
}

// validateCompositeConstraints ensures composite states have regions
// UML Constraint: A composite state must have at least one region
func (s *State) validateCompositeConstraints(context *ValidationContext, errors *ValidationErrors) {
	if s.IsComposite {
		// Composite states must have at least one region
		if len(s.Regions) == 0 {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Regions",
				"composite state must have at least one region (UML constraint)",
				context.Path,
			)
		}

		// Composite states cannot be simple states
		if s.IsSimple {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"IsSimple",
				"state cannot be both composite and simple (UML constraint)",
				context.Path,
			)
		}

		// Validate each region in the composite state
		for i, region := range s.Regions {
			if region == nil {
				continue // This will be caught by collection validation
			}

			regionContext := context.WithPathIndex("Regions", i)

			// Ensure region has proper identification
			if region.ID == "" {
				errors.AddError(
					ErrorTypeConstraint,
					"State",
					"Regions",
					fmt.Sprintf("region at index %d must have a valid ID (UML constraint)", i),
					regionContext.Path,
				)
			}

			// Ensure region name is meaningful
			if region.Name == "" {
				errors.AddError(
					ErrorTypeConstraint,
					"State",
					"Regions",
					fmt.Sprintf("region at index %d should have a descriptive name (UML best practice)", i),
					regionContext.Path,
				)
			}
		}

		// If orthogonal, must have multiple regions
		if s.IsOrthogonal && len(s.Regions) < 2 {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Regions",
				"orthogonal composite state must have at least two regions (UML constraint)",
				context.Path,
			)
		}
	} else {
		// Non-composite states should not have regions
		if len(s.Regions) > 0 {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Regions",
				"non-composite state cannot have regions (UML constraint)",
				context.Path,
			)
		}

		// Non-composite states cannot be orthogonal
		if s.IsOrthogonal {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"IsOrthogonal",
				"non-composite state cannot be orthogonal (UML constraint)",
				context.Path,
			)
		}
	}
}

// validateSubmachineConstraints validates submachine state constraints
// UML Constraint: A submachine state must reference a valid state machine and have proper connection points
func (s *State) validateSubmachineConstraints(context *ValidationContext, errors *ValidationErrors) {
	if s.IsSubmachineState {
		// Submachine states must reference a state machine
		if s.Submachine == nil {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Submachine",
				"submachine state must reference a valid state machine (UML constraint)",
				context.Path,
			)
		} else {
			// Validate the referenced submachine
			submachineContext := context.WithPath("Submachine")

			// Ensure submachine has proper identification
			if s.Submachine.ID == "" {
				errors.AddError(
					ErrorTypeConstraint,
					"State",
					"Submachine",
					"referenced submachine must have a valid ID (UML constraint)",
					submachineContext.Path,
				)
			}

			// Ensure submachine name is meaningful
			if s.Submachine.Name == "" {
				errors.AddError(
					ErrorTypeConstraint,
					"State",
					"Submachine",
					"referenced submachine should have a descriptive name (UML best practice)",
					submachineContext.Path,
				)
			}

			// Validate connection point compatibility
			s.validateConnectionPointCompatibility(context, errors)
		}

		// Submachine states should not be composite in the traditional sense
		if s.IsComposite {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"IsComposite",
				"submachine state should not be marked as composite (use submachine reference instead) (UML constraint)",
				context.Path,
			)
		}

		// Submachine states should not have their own regions
		if len(s.Regions) > 0 {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Regions",
				"submachine state should not have its own regions (use submachine reference instead) (UML constraint)",
				context.Path,
			)
		}
	} else {
		// Non-submachine states should not reference a submachine
		if s.Submachine != nil {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Submachine",
				"non-submachine state should not reference a submachine (UML constraint)",
				context.Path,
			)
		}

		// Non-submachine states should not have connection point references
		if len(s.Connections) > 0 {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Connections",
				"non-submachine state should not have connection point references (UML constraint)",
				context.Path,
			)
		}
	}
}

// validateConnectionPointCompatibility validates connection points between submachine state and referenced submachine
func (s *State) validateConnectionPointCompatibility(context *ValidationContext, errors *ValidationErrors) {
	if s.Submachine == nil {
		return
	}

	// Create maps of available connection points in the submachine
	submachineEntryPoints := make(map[string]*Pseudostate)
	submachineExitPoints := make(map[string]*Pseudostate)

	for _, cp := range s.Submachine.ConnectionPoints {
		if cp == nil {
			continue
		}

		switch cp.Kind {
		case PseudostateKindEntryPoint:
			submachineEntryPoints[cp.ID] = cp
		case PseudostateKindExitPoint:
			submachineExitPoints[cp.ID] = cp
		}
	}

	// Validate connection point references
	for i, conn := range s.Connections {
		if conn == nil {
			continue
		}

		connContext := context.WithPathIndex("Connections", i)

		// Validate entry point references
		for j, entry := range conn.Entry {
			if entry == nil {
				continue
			}

			if _, exists := submachineEntryPoints[entry.ID]; !exists {
				errors.AddError(
					ErrorTypeConstraint,
					"State",
					"Connections",
					fmt.Sprintf("connection point reference at index %d references entry point '%s' that does not exist in submachine (UML constraint)", i, entry.ID),
					connContext.WithPathIndex("Entry", j).Path,
				)
			}
		}

		// Validate exit point references
		for j, exit := range conn.Exit {
			if exit == nil {
				continue
			}

			if _, exists := submachineExitPoints[exit.ID]; !exists {
				errors.AddError(
					ErrorTypeConstraint,
					"State",
					"Connections",
					fmt.Sprintf("connection point reference at index %d references exit point '%s' that does not exist in submachine (UML constraint)", i, exit.ID),
					connContext.WithPathIndex("Exit", j).Path,
				)
			}
		}
	}
}

// validateBehaviorConsistency validates state behavior consistency (entry/exit/do activities)
// UML Constraint: State behaviors must be consistent and properly defined
func (s *State) validateBehaviorConsistency(context *ValidationContext, errors *ValidationErrors) {
	// Validate entry behavior
	if s.Entry != nil {
		entryContext := context.WithPath("Entry")
		s.validateBehavior(s.Entry, "entry", entryContext, errors)
	}

	// Validate exit behavior
	if s.Exit != nil {
		exitContext := context.WithPath("Exit")
		s.validateBehavior(s.Exit, "exit", exitContext, errors)
	}

	// Validate do activity behavior
	if s.DoActivity != nil {
		doContext := context.WithPath("DoActivity")
		s.validateBehavior(s.DoActivity, "do activity", doContext, errors)
	}

	// Validate behavior consistency rules
	s.validateBehaviorInteractions(context, errors)
}

// validateBehavior validates a single behavior object
func (s *State) validateBehavior(behavior *Behavior, behaviorType string, context *ValidationContext, errors *ValidationErrors) {
	if behavior == nil {
		return
	}

	// Validate behavior has proper identification
	if behavior.ID == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"State",
			"Behavior",
			fmt.Sprintf("%s behavior must have a valid ID (UML constraint)", behaviorType),
			context.Path,
		)
	}

	// Validate behavior name is meaningful (optional but recommended)
	if behavior.Name == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"State",
			"Behavior",
			fmt.Sprintf("%s behavior should have a descriptive name (UML best practice)", behaviorType),
			context.Path,
		)
	}

	// Validate behavior specification exists
	if behavior.Specification == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"State",
			"Behavior",
			fmt.Sprintf("%s behavior must have a valid specification (UML constraint)", behaviorType),
			context.Path,
		)
	}

	// Validate behavior language consistency
	if behavior.Language != "" && behavior.Specification == "" {
		errors.AddError(
			ErrorTypeConstraint,
			"State",
			"Behavior",
			fmt.Sprintf("%s behavior specifies language '%s' but has no specification content (UML constraint)", behaviorType, behavior.Language),
			context.Path,
		)
	}
}

// validateBehaviorInteractions validates interactions between different state behaviors
func (s *State) validateBehaviorInteractions(context *ValidationContext, errors *ValidationErrors) {
	// Entry and exit behaviors should not conflict
	if s.Entry != nil && s.Exit != nil {
		// Check for potential naming conflicts
		if s.Entry.Name == s.Exit.Name && s.Entry.Name != "" {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				"entry and exit behaviors should have distinct names to avoid confusion (UML best practice)",
				context.Path,
			)
		}

		// Check for language consistency
		if s.Entry.Language != "" && s.Exit.Language != "" && s.Entry.Language != s.Exit.Language {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				fmt.Sprintf("entry behavior uses language '%s' while exit behavior uses '%s', consider consistency (UML best practice)", s.Entry.Language, s.Exit.Language),
				context.Path,
			)
		}

		// Check for specification conflicts (same ID but different specifications)
		if s.Entry.ID == s.Exit.ID && s.Entry.Specification != s.Exit.Specification {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				"entry and exit behaviors have the same ID but different specifications, which may cause confusion (UML best practice)",
				context.Path,
			)
		}
	}

	// Do activity should be compatible with entry/exit behaviors
	if s.DoActivity != nil {
		if s.Entry != nil && s.Entry.Language != "" && s.DoActivity.Language != "" && s.Entry.Language != s.DoActivity.Language {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				fmt.Sprintf("entry behavior uses language '%s' while do activity uses '%s', consider consistency (UML best practice)", s.Entry.Language, s.DoActivity.Language),
				context.Path,
			)
		}

		if s.Exit != nil && s.Exit.Language != "" && s.DoActivity.Language != "" && s.Exit.Language != s.DoActivity.Language {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				fmt.Sprintf("exit behavior uses language '%s' while do activity uses '%s', consider consistency (UML best practice)", s.Exit.Language, s.DoActivity.Language),
				context.Path,
			)
		}

		// Check for ID conflicts between do activity and entry/exit behaviors
		if s.Entry != nil && s.DoActivity.ID == s.Entry.ID {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				"do activity and entry behavior have the same ID, which may cause confusion (UML best practice)",
				context.Path,
			)
		}

		if s.Exit != nil && s.DoActivity.ID == s.Exit.ID {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Behaviors",
				"do activity and exit behavior have the same ID, which may cause confusion (UML best practice)",
				context.Path,
			)
		}
	}

	// Validate behavior semantic consistency
	s.validateBehaviorSemantics(context, errors)
}

// validateBehaviorSemantics validates semantic consistency of state behaviors
func (s *State) validateBehaviorSemantics(context *ValidationContext, errors *ValidationErrors) {
	// Entry behavior should set up state conditions
	if s.Entry != nil && s.Entry.Specification != "" {
		// Check if entry behavior specification suggests it's doing cleanup (which should be in exit)
		entrySpec := strings.ToLower(s.Entry.Specification)
		if strings.Contains(entrySpec, "cleanup") || strings.Contains(entrySpec, "destroy") || strings.Contains(entrySpec, "finalize") {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Entry",
				"entry behavior specification suggests cleanup operations, which should typically be in exit behavior (UML semantics)",
				context.Path,
			)
		}
	}

	// Exit behavior should clean up state conditions
	if s.Exit != nil && s.Exit.Specification != "" {
		// Check if exit behavior specification suggests it's doing initialization (which should be in entry)
		exitSpec := strings.ToLower(s.Exit.Specification)
		if strings.Contains(exitSpec, "initialize") || strings.Contains(exitSpec, "setup") || strings.Contains(exitSpec, "create") {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"Exit",
				"exit behavior specification suggests initialization operations, which should typically be in entry behavior (UML semantics)",
				context.Path,
			)
		}
	}

	// Do activity should represent ongoing behavior
	if s.DoActivity != nil && s.DoActivity.Specification != "" {
		doSpec := strings.ToLower(s.DoActivity.Specification)
		// Check if do activity suggests one-time operations (which should be in entry/exit)
		if strings.Contains(doSpec, "initialize") || strings.Contains(doSpec, "setup") || strings.Contains(doSpec, "cleanup") || strings.Contains(doSpec, "destroy") {
			errors.AddError(
				ErrorTypeConstraint,
				"State",
				"DoActivity",
				"do activity specification suggests one-time operations, which should typically be in entry or exit behaviors (UML semantics)",
				context.Path,
			)
		}
	}
}
