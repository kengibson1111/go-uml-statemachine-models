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

	// UML constraint validations
	t.validateSourceTarget(context, errors)
	t.validateKindConstraints(context, errors)
	t.validateContainment(context, errors)
}

// validateSourceTarget ensures source/target compatibility
// UML Constraint: Transition source and target must be compatible according to UML rules
func (t *Transition) validateSourceTarget(context *ValidationContext, errors *ValidationErrors) {
	if t.Source == nil || t.Target == nil {
		// Already validated by required field validation
		return
	}

	source := t.Source
	target := t.Target

	// Validate source vertex constraints
	t.validateSourceConstraints(source, context, errors)

	// Validate target vertex constraints
	t.validateTargetConstraints(target, context, errors)

	// Validate source-target compatibility
	t.validateSourceTargetCompatibility(source, target, context, errors)
}

// validateKindConstraints validates internal/local/external transition rules
// UML Constraint: Different transition kinds have specific structural requirements
func (t *Transition) validateKindConstraints(context *ValidationContext, errors *ValidationErrors) {
	if t.Source == nil || t.Target == nil {
		// Already validated by required field validation
		return
	}

	switch t.Kind {
	case TransitionKindInternal:
		t.validateInternalTransitionConstraints(context, errors)
	case TransitionKindLocal:
		t.validateLocalTransitionConstraints(context, errors)
	case TransitionKindExternal:
		t.validateExternalTransitionConstraints(context, errors)
	default:
		// Invalid kind already caught by basic validation
	}
}

// validateContainment validates transition containment within appropriate regions
// UML Constraint: Transitions must be properly contained within regions that contain their source/target vertices
func (t *Transition) validateContainment(context *ValidationContext, errors *ValidationErrors) {
	if t.Source == nil || t.Target == nil {
		// Already validated by required field validation
		return
	}

	// Get the containing region from context
	region := context.Region
	if region == nil {
		// If no region context, we can't validate containment constraints
		// This is acceptable for standalone transition validation
		return
	}

	// Validate source vertex containment
	t.validateVertexContainment(t.Source, "Source", region, context, errors)

	// Validate target vertex containment based on transition kind
	t.validateTargetContainment(t.Target, region, context, errors)
}

// validateSourceConstraints validates constraints specific to source vertices
func (t *Transition) validateSourceConstraints(source *Vertex, context *ValidationContext, errors *ValidationErrors) {
	// Final states cannot have outgoing transitions
	if source.Type == "finalstate" {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Source",
			"final state cannot be the source of a transition (UML constraint)",
			context.Path,
		)
	}

	// Validate pseudostate-specific source constraints
	if source.Type == "pseudostate" {
		t.validatePseudostateSourceConstraints(source, context, errors)
	}
}

// validateTargetConstraints validates constraints specific to target vertices
func (t *Transition) validateTargetConstraints(target *Vertex, context *ValidationContext, errors *ValidationErrors) {
	// Validate pseudostate-specific target constraints
	if target.Type == "pseudostate" {
		t.validatePseudostateTargetConstraints(target, context, errors)
	}

	// Initial pseudostates cannot be targets of transitions (except from outside the region)
	if target.Type == "pseudostate" && t.isInitialPseudostate(target) {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Target",
			"initial pseudostate cannot be the target of a transition within the same region (UML constraint)",
			context.Path,
		)
	}
}

// validateSourceTargetCompatibility validates compatibility between source and target
func (t *Transition) validateSourceTargetCompatibility(source, target *Vertex, context *ValidationContext, errors *ValidationErrors) {
	// Validate reflexive transitions (self-transitions)
	if source.ID == target.ID {
		// Self-transitions should typically be internal or local
		if t.Kind == TransitionKindExternal {
			// This is allowed but might indicate a design issue
			// We'll issue a warning-level constraint error
			errors.AddError(
				ErrorTypeConstraint,
				"Transition",
				"Kind",
				"self-transition with external kind may cause exit/entry actions to be executed (UML design consideration)",
				context.Path,
			)
		}
	}

	// Validate cross-region transitions
	t.validateCrossRegionTransition(source, target, context, errors)
}

// validateInternalTransitionConstraints validates constraints for internal transitions
func (t *Transition) validateInternalTransitionConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Internal transitions must have the same source and target
	if t.Source.ID != t.Target.ID {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Kind",
			"internal transition must have the same source and target vertex (UML constraint)",
			context.Path,
		)
	}

	// Internal transitions should not cause state exit/entry
	// This is more of a semantic constraint that affects behavior
	if t.Source.Type != "state" {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Source",
			"internal transition source should be a state (UML constraint)",
			context.Path,
		)
	}
}

// validateLocalTransitionConstraints validates constraints for local transitions
func (t *Transition) validateLocalTransitionConstraints(context *ValidationContext, errors *ValidationErrors) {
	// Local transitions are within the same composite state
	// The source and target must be in the same region or in nested regions of the same composite state

	// For now, we validate that both source and target are proper vertices
	if t.Source.Type == "pseudostate" && t.isConnectionPoint(t.Source) {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Source",
			"local transition should not originate from connection points (UML constraint)",
			context.Path,
		)
	}

	if t.Target.Type == "pseudostate" && t.isConnectionPoint(t.Target) {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Target",
			"local transition should not target connection points (UML constraint)",
			context.Path,
		)
	}
}

// validateExternalTransitionConstraints validates constraints for external transitions
func (t *Transition) validateExternalTransitionConstraints(context *ValidationContext, errors *ValidationErrors) {
	// External transitions can cross region boundaries
	// They cause exit from source state and entry to target state

	// Validate that external self-transitions may cause exit/entry actions
	if t.Source.ID == t.Target.ID && t.Source.Type == "state" {
		// This is allowed but might indicate a design issue
		// We'll issue a warning-level constraint error
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Kind",
			"self-transition with external kind may cause exit/entry actions to be executed (UML design consideration)",
			context.Path,
		)
	}

	// External transitions can use connection points
	// No additional constraints for external transitions in basic UML
}

// validatePseudostateSourceConstraints validates source constraints for pseudostates
func (t *Transition) validatePseudostateSourceConstraints(source *Vertex, context *ValidationContext, errors *ValidationErrors) {
	// Different pseudostate kinds have different outgoing transition constraints
	// We use naming conventions to identify pseudostate kinds

	if t.isTerminatePseudostate(source) {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			"Source",
			"terminate pseudostate cannot have outgoing transitions (UML constraint)",
			context.Path,
		)
	}

	// Junction and choice pseudostates must have at least one outgoing transition
	// This is validated at the region level, but we can check for basic consistency
	if t.isJunctionOrChoice(source) {
		// These should have guards or else conditions
		// This is a more complex validation that would require analyzing all transitions
	}
}

// validatePseudostateTargetConstraints validates target constraints for pseudostates
func (t *Transition) validatePseudostateTargetConstraints(target *Vertex, context *ValidationContext, errors *ValidationErrors) {
	// Initial pseudostates should not have incoming transitions from within the same region
	if t.isInitialPseudostate(target) {
		// This is already handled in validateTargetConstraints
	}

	// History pseudostates have specific incoming transition rules
	if t.isHistoryPseudostate(target) {
		// History pseudostates should have at most one incoming transition
		// This would require analyzing all transitions in the region
	}
}

// validateVertexContainment validates that a vertex is contained in the specified region
func (t *Transition) validateVertexContainment(vertex *Vertex, vertexRole string, region *Region, context *ValidationContext, errors *ValidationErrors) {
	// Check if vertex is in the region's vertices collection
	found := false
	for _, regionVertex := range region.Vertices {
		if regionVertex != nil && regionVertex.ID == vertex.ID {
			found = true
			break
		}
	}

	if !found {
		// Also check in the states collection
		for _, state := range region.States {
			if state != nil && state.ID == vertex.ID {
				found = true
				break
			}
		}
	}

	if !found {
		errors.AddError(
			ErrorTypeConstraint,
			"Transition",
			vertexRole,
			fmt.Sprintf("%s vertex (ID: %s) is not contained in the transition's region (UML constraint)", vertexRole, vertex.ID),
			context.Path,
		)
	}
}

// validateTargetContainment validates target containment based on transition kind
func (t *Transition) validateTargetContainment(target *Vertex, region *Region, context *ValidationContext, errors *ValidationErrors) {
	switch t.Kind {
	case TransitionKindInternal:
		// Internal transitions must have target in same region
		t.validateVertexContainment(target, "Target", region, context, errors)
	case TransitionKindLocal:
		// Local transitions must have target in same composite state (same or nested region)
		t.validateVertexContainment(target, "Target", region, context, errors)
	case TransitionKindExternal:
		// External transitions can have targets outside the region
		// We still validate if the target is in this region, but don't require it
		// The target should exist somewhere in the state machine
		t.validateExternalTargetAccessibility(target, region, context, errors)
	}
}

// validateExternalTargetAccessibility validates that external transition targets are accessible
func (t *Transition) validateExternalTargetAccessibility(target *Vertex, region *Region, context *ValidationContext, errors *ValidationErrors) {
	// For external transitions, the target might be in a different region
	// We should validate that the target exists somewhere in the state machine

	// Check if target is in current region
	found := false
	for _, regionVertex := range region.Vertices {
		if regionVertex != nil && regionVertex.ID == target.ID {
			found = true
			break
		}
	}

	if !found {
		// Check in states collection
		for _, state := range region.States {
			if state != nil && state.ID == target.ID {
				found = true
				break
			}
		}
	}

	// If not found in current region and we have state machine context, check other regions
	if !found && context.StateMachine != nil {
		// Check other regions in the state machine
		for _, otherRegion := range context.StateMachine.Regions {
			if otherRegion == nil || otherRegion.ID == region.ID {
				continue // Skip current region or nil regions
			}

			// Check vertices in other regions
			for _, vertex := range otherRegion.Vertices {
				if vertex != nil && vertex.ID == target.ID {
					found = true
					break
				}
			}

			if found {
				break
			}

			// Check states in other regions
			for _, state := range otherRegion.States {
				if state != nil && state.ID == target.ID {
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		// If still not found, it's an error
		if !found {
			errors.AddError(
				ErrorTypeConstraint,
				"Transition",
				"Target",
				fmt.Sprintf("external transition target (ID: %s) not found in any region of the state machine (UML constraint)", target.ID),
				context.Path,
			)
		}
	}
	// If no state machine context, we can't validate cross-region accessibility
}

// validateCrossRegionTransition validates transitions that cross region boundaries
func (t *Transition) validateCrossRegionTransition(source, target *Vertex, context *ValidationContext, errors *ValidationErrors) {
	// This is a complex validation that would require full state machine context
	// For now, we validate basic constraints

	if t.Kind == TransitionKindInternal && source.ID != target.ID {
		// Already handled in validateInternalTransitionConstraints
		return
	}

	// Cross-region transitions should use appropriate connection points
	if t.Kind == TransitionKindExternal {
		// External transitions crossing composite state boundaries should use connection points
		// This requires more context about the state hierarchy
	}
}

// Helper methods for identifying pseudostate types

// isInitialPseudostate checks if a vertex is an initial pseudostate
func (t *Transition) isInitialPseudostate(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	// Use naming conventions to identify initial pseudostates
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

// isTerminatePseudostate checks if a vertex is a terminate pseudostate
func (t *Transition) isTerminatePseudostate(vertex *Vertex) bool {
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

// isHistoryPseudostate checks if a vertex is a history pseudostate
func (t *Transition) isHistoryPseudostate(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	name := vertex.Name
	id := vertex.ID

	historyPatterns := []string{
		"history", "History", "HISTORY",
		"deepHistory", "DeepHistory", "DEEP_HISTORY",
		"shallowHistory", "ShallowHistory", "SHALLOW_HISTORY",
		"H", "H*",
	}

	for _, pattern := range historyPatterns {
		if name == pattern || id == pattern {
			return true
		}
	}

	return false
}

// isJunctionOrChoice checks if a vertex is a junction or choice pseudostate
func (t *Transition) isJunctionOrChoice(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	name := vertex.Name
	id := vertex.ID

	junctionChoicePatterns := []string{
		"junction", "Junction", "JUNCTION",
		"choice", "Choice", "CHOICE",
		"decision", "Decision", "DECISION",
	}

	for _, pattern := range junctionChoicePatterns {
		if name == pattern || id == pattern {
			return true
		}
	}

	return false
}

// isConnectionPoint checks if a vertex is a connection point (entry/exit point)
func (t *Transition) isConnectionPoint(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	name := vertex.Name
	id := vertex.ID

	connectionPointPatterns := []string{
		"entryPoint", "EntryPoint", "ENTRY_POINT",
		"exitPoint", "ExitPoint", "EXIT_POINT",
		"entry", "Entry", "ENTRY",
		"exit", "Exit", "EXIT",
	}

	for _, pattern := range connectionPointPatterns {
		if name == pattern || id == pattern {
			return true
		}
	}

	return false
}
