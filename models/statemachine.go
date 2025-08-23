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

	// Structural integrity validation
	sm.validateStructuralIntegrity(context, errors)
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

	// UML constraint validations
	r.validateInitialStates(context, errors)
	r.validateVertexContainment(context, errors)
	r.validateTransitionScope(context, errors)

	// Structural integrity validation
	r.validateStructuralIntegrity(context, errors)
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

// validateInitialStates ensures at most one initial pseudostate per region
// UML Constraint: A Region can have at most one initial pseudostate
func (r *Region) validateInitialStates(context *ValidationContext, errors *ValidationErrors) {
	initialCount := 0
	var initialIndices []int

	// Check vertices for initial pseudostates using naming conventions
	for i, vertex := range r.Vertices {
		if vertex == nil {
			continue // This will be caught by collection validation
		}

		// Check if this vertex is an initial pseudostate
		if vertex.Type == "pseudostate" && r.isInitialPseudostate(vertex) {
			initialCount++
			initialIndices = append(initialIndices, i)
		}
	}

	// Also check states collection in case pseudostates are stored there
	for i, state := range r.States {
		if state == nil {
			continue
		}

		if state.Type == "pseudostate" && r.isInitialPseudostate(&state.Vertex) {
			initialCount++
			initialIndices = append(initialIndices, i)
		}
	}

	if initialCount > 1 {
		errors.AddError(
			ErrorTypeMultiplicity,
			"Region",
			"Vertices",
			fmt.Sprintf("Region can have at most one initial pseudostate, found %d at indices: %v (UML constraint)", initialCount, initialIndices),
			context.Path,
		)
	}
}

// validateVertexContainment verifies vertex collections are properly structured
// UML Constraint: States and Vertices collections should not overlap - states go in States, pseudostates/final states go in Vertices
func (r *Region) validateVertexContainment(context *ValidationContext, errors *ValidationErrors) {
	// Create maps for checking overlaps
	stateIDs := make(map[string]bool)
	vertexIDs := make(map[string]bool)

	// Collect state IDs
	for _, state := range r.States {
		if state != nil {
			stateIDs[state.ID] = true
		}
	}

	// Collect vertex IDs and check for overlaps with states
	for _, vertex := range r.Vertices {
		if vertex != nil {
			vertexIDs[vertex.ID] = true
			// Check if this vertex ID is also used by a state (which would be incorrect)
			if stateIDs[vertex.ID] {
				errors.AddError(
					ErrorTypeConstraint,
					"Region",
					"States",
					fmt.Sprintf("duplicate vertex ID '%s' found in vertices collection and states collection at index %d (structural integrity violation)", vertex.ID, 0),
					context.Path,
				)
			}
		}
	}

	// Validate that vertices have consistent containment
	// All vertices should logically belong to this region
	for i, vertex := range r.Vertices {
		if vertex == nil {
			continue
		}

		// Ensure vertex has proper identification
		if vertex.ID == "" {
			errors.AddError(
				ErrorTypeConstraint,
				"Region",
				"Vertices",
				fmt.Sprintf("vertex at index %d must have a valid ID for proper containment (UML constraint)", i),
				context.WithPathIndex("Vertices", i).Path,
			)
		}

		// Validate vertex type is appropriate for vertices collection
		// States should be in States collection, not Vertices collection
		if vertex.Type == "state" {
			errors.AddError(
				ErrorTypeConstraint,
				"Region",
				"Vertices",
				fmt.Sprintf("vertex at index %d has type 'state' but should be in States collection, not Vertices collection (UML constraint)", i),
				context.WithPathIndex("Vertices", i).Path,
			)
		} else {
			// For non-state vertices, validate they have appropriate types
			validTypes := []string{"pseudostate", "finalstate"}
			isValidType := false
			for _, validType := range validTypes {
				if vertex.Type == validType {
					isValidType = true
					break
				}
			}

			if !isValidType {
				errors.AddError(
					ErrorTypeConstraint,
					"Region",
					"Vertices",
					fmt.Sprintf("vertex at index %d has invalid type '%s' for vertices collection (UML constraint)", i, vertex.Type),
					context.WithPathIndex("Vertices", i).Path,
				)
			}
		}
	}
}

// validateTransitionScope ensures transitions connect appropriate vertices
// UML Constraint: Transitions must connect vertices that are appropriately scoped within the region
func (r *Region) validateTransitionScope(context *ValidationContext, errors *ValidationErrors) {
	// Create a map of vertex IDs for quick lookup
	vertexIDs := make(map[string]bool)
	for _, vertex := range r.Vertices {
		if vertex != nil {
			vertexIDs[vertex.ID] = true
		}
	}

	// Validate each transition
	for i, transition := range r.Transitions {
		if transition == nil {
			continue // This will be caught by collection validation
		}

		transitionContext := context.WithPathIndex("Transitions", i)

		// Validate source vertex is in this region
		if transition.Source != nil {
			if !vertexIDs[transition.Source.ID] {
				errors.AddError(
					ErrorTypeConstraint,
					"Region",
					"Transitions",
					fmt.Sprintf("transition at index %d has source vertex (ID: %s) that is not contained in this region (UML constraint)", i, transition.Source.ID),
					transitionContext.Path,
				)
			}
		}

		// Validate target vertex is in this region or is appropriately accessible
		if transition.Target != nil {
			if !vertexIDs[transition.Target.ID] {
				// For external transitions, the target might be in a different region
				// but for internal and local transitions, target must be in same region
				if transition.Kind == TransitionKindInternal || transition.Kind == TransitionKindLocal {
					errors.AddError(
						ErrorTypeConstraint,
						"Region",
						"Transitions",
						fmt.Sprintf("transition at index %d has target vertex (ID: %s) that is not contained in this region, but transition kind is %s (UML constraint)", i, transition.Target.ID, transition.Kind),
						transitionContext.Path,
					)
				}
				// For external transitions, we allow targets outside the region
				// but we should validate they exist somewhere in the state machine
			}
		}

		// Validate transition kind constraints
		if transition.Source != nil && transition.Target != nil {
			// Internal transitions must have the same source and target
			if transition.Kind == TransitionKindInternal {
				if transition.Source.ID != transition.Target.ID {
					errors.AddError(
						ErrorTypeConstraint,
						"Region",
						"Transitions",
						fmt.Sprintf("internal transition at index %d must have the same source and target vertex (UML constraint)", i),
						transitionContext.Path,
					)
				}
			}

			// Validate that source and target are compatible types
			r.validateTransitionVertexCompatibility(transition, i, transitionContext, errors)
		}
	}
}

// validateTransitionVertexCompatibility validates that source and target vertices are compatible
func (r *Region) validateTransitionVertexCompatibility(transition *Transition, index int, context *ValidationContext, errors *ValidationErrors) {
	if transition.Source == nil || transition.Target == nil {
		return // Already validated by required field validation
	}

	source := transition.Source
	target := transition.Target

	// Validate pseudostate transition rules
	if source.Type == "pseudostate" {
		// Initial pseudostates can only have outgoing transitions
		if source.Name == "Initial" || source.ID == "initial" {
			// This is handled by the pseudostate validation, but we can add region-specific checks
		}

		// Junction and choice pseudostates have specific rules
		// (We'd need access to PseudostateKind to implement these fully)
	}

	if target.Type == "pseudostate" {
		// Final states cannot have outgoing transitions (but can be targets)
		// Terminate pseudostates have specific rules
	}

	// Validate that final states don't have outgoing transitions
	if source.Type == "finalstate" {
		errors.AddError(
			ErrorTypeConstraint,
			"Region",
			"Transitions",
			fmt.Sprintf("transition at index %d has a final state as source, which is not allowed (UML constraint)", index),
			context.Path,
		)
	}

	// Additional compatibility checks can be added here based on UML rules
}

// isInitialPseudostate checks if a vertex represents an initial pseudostate
// This is a helper method that uses naming conventions to identify initial pseudostates
func (r *Region) isInitialPseudostate(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	// Check common naming patterns for initial pseudostates
	name := vertex.Name
	id := vertex.ID

	// Common patterns for initial pseudostates
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

// validateStructuralIntegrity performs structural integrity validation for StateMachine
func (sm *StateMachine) validateStructuralIntegrity(context *ValidationContext, errors *ValidationErrors) {
	// Create a reference validator for this state machine
	refValidator := NewReferenceValidator()

	// Validate references within this state machine
	if err := refValidator.ValidateReferencesInContext(sm, context); err != nil {
		// Extract errors from the reference validator and add them to our error collection
		if refErrors, ok := err.(*ValidationErrors); ok {
			for _, refError := range refErrors.Errors {
				errors.Add(refError)
			}
		} else {
			// Single error case
			errors.AddError(
				ErrorTypeReference,
				"StateMachine",
				"StructuralIntegrity",
				err.Error(),
				context.Path,
			)
		}
	}

	// Additional state machine specific structural validations
	sm.validateRegionConsistency(context, errors)
	sm.validateConnectionPointConsistency(context, errors)
}

// validateRegionConsistency validates consistency between regions
func (sm *StateMachine) validateRegionConsistency(context *ValidationContext, errors *ValidationErrors) {
	// Check for duplicate region IDs
	regionIDs := make(map[string]int)
	for i, region := range sm.Regions {
		if region == nil {
			continue
		}

		if prevIndex, exists := regionIDs[region.ID]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				"StateMachine",
				"Regions",
				fmt.Sprintf("duplicate region ID '%s' found at indices %d and %d (structural integrity violation)", region.ID, prevIndex, i),
				context.WithPathIndex("Regions", i).Path,
			)
		} else {
			regionIDs[region.ID] = i
		}
	}

	// Validate region names are unique (best practice)
	regionNames := make(map[string]int)
	for i, region := range sm.Regions {
		if region == nil || region.Name == "" {
			continue
		}

		if prevIndex, exists := regionNames[region.Name]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				"StateMachine",
				"Regions",
				fmt.Sprintf("duplicate region name '%s' found at indices %d and %d (may cause confusion)", region.Name, prevIndex, i),
				context.WithPathIndex("Regions", i).Path,
			)
		} else {
			regionNames[region.Name] = i
		}
	}
}

// validateConnectionPointConsistency validates consistency of connection points
func (sm *StateMachine) validateConnectionPointConsistency(context *ValidationContext, errors *ValidationErrors) {
	// Check for duplicate connection point IDs
	cpIDs := make(map[string]int)
	for i, cp := range sm.ConnectionPoints {
		if cp == nil {
			continue
		}

		if prevIndex, exists := cpIDs[cp.ID]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				"StateMachine",
				"ConnectionPoints",
				fmt.Sprintf("duplicate connection point ID '%s' found at indices %d and %d (structural integrity violation)", cp.ID, prevIndex, i),
				context.WithPathIndex("ConnectionPoints", i).Path,
			)
		} else {
			cpIDs[cp.ID] = i
		}
	}

	// Validate connection point names are unique within their kind
	entryNames := make(map[string]int)
	exitNames := make(map[string]int)

	for i, cp := range sm.ConnectionPoints {
		if cp == nil || cp.Name == "" {
			continue
		}

		switch cp.Kind {
		case PseudostateKindEntryPoint:
			if prevIndex, exists := entryNames[cp.Name]; exists {
				errors.AddError(
					ErrorTypeConstraint,
					"StateMachine",
					"ConnectionPoints",
					fmt.Sprintf("duplicate entry point name '%s' found at indices %d and %d (may cause confusion)", cp.Name, prevIndex, i),
					context.WithPathIndex("ConnectionPoints", i).Path,
				)
			} else {
				entryNames[cp.Name] = i
			}
		case PseudostateKindExitPoint:
			if prevIndex, exists := exitNames[cp.Name]; exists {
				errors.AddError(
					ErrorTypeConstraint,
					"StateMachine",
					"ConnectionPoints",
					fmt.Sprintf("duplicate exit point name '%s' found at indices %d and %d (may cause confusion)", cp.Name, prevIndex, i),
					context.WithPathIndex("ConnectionPoints", i).Path,
				)
			} else {
				exitNames[cp.Name] = i
			}
		}
	}
}

// validateStructuralIntegrity performs structural integrity validation for Region
func (r *Region) validateStructuralIntegrity(context *ValidationContext, errors *ValidationErrors) {
	// Validate vertex ID consistency between states and vertices collections
	r.validateVertexIDConsistency(context, errors)

	// Validate transition reference consistency
	r.validateTransitionReferenceConsistency(context, errors)

	// Validate containment relationships
	r.validateContainmentRelationships(context, errors)
}

// validateVertexIDConsistency validates that vertex IDs are consistent across collections
func (r *Region) validateVertexIDConsistency(context *ValidationContext, errors *ValidationErrors) {
	// Check for duplicate vertex IDs within the region
	vertexIDs := make(map[string]string) // ID -> collection type

	// Check vertices collection
	for i, vertex := range r.Vertices {
		if vertex == nil {
			continue
		}

		if collectionType, exists := vertexIDs[vertex.ID]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				"Region",
				"Vertices",
				fmt.Sprintf("duplicate vertex ID '%s' found in %s collection and vertices collection at index %d (structural integrity violation)", vertex.ID, collectionType, i),
				context.WithPathIndex("Vertices", i).Path,
			)
		} else {
			vertexIDs[vertex.ID] = "vertices"
		}
	}

	// Check states collection
	for i, state := range r.States {
		if state == nil {
			continue
		}

		if collectionType, exists := vertexIDs[state.ID]; exists {
			errors.AddError(
				ErrorTypeConstraint,
				"Region",
				"States",
				fmt.Sprintf("duplicate vertex ID '%s' found in %s collection and states collection at index %d (structural integrity violation)", state.ID, collectionType, i),
				context.WithPathIndex("States", i).Path,
			)
		} else {
			vertexIDs[state.ID] = "states"
		}
	}
}

// validateTransitionReferenceConsistency validates that transitions reference valid vertices
func (r *Region) validateTransitionReferenceConsistency(context *ValidationContext, errors *ValidationErrors) {
	// Build a map of all available vertices in this region
	availableVertices := make(map[string]bool)

	for _, vertex := range r.Vertices {
		if vertex != nil {
			availableVertices[vertex.ID] = true
		}
	}

	for _, state := range r.States {
		if state != nil {
			availableVertices[state.ID] = true
		}
	}

	// Validate each transition's source and target references
	for i, transition := range r.Transitions {
		if transition == nil {
			continue
		}

		transitionContext := context.WithPathIndex("Transitions", i)

		// Check source reference
		if transition.Source != nil {
			if !availableVertices[transition.Source.ID] {
				// For internal and local transitions, source must be in this region
				if transition.Kind == TransitionKindInternal || transition.Kind == TransitionKindLocal {
					errors.AddError(
						ErrorTypeReference,
						"Region",
						"Transitions",
						fmt.Sprintf("transition at index %d references source vertex '%s' that is not available in this region (structural integrity violation)", i, transition.Source.ID),
						transitionContext.Path,
					)
				}
			}
		}

		// Check target reference
		if transition.Target != nil {
			if !availableVertices[transition.Target.ID] {
				// For internal and local transitions, target must be in this region
				if transition.Kind == TransitionKindInternal || transition.Kind == TransitionKindLocal {
					errors.AddError(
						ErrorTypeReference,
						"Region",
						"Transitions",
						fmt.Sprintf("transition at index %d references target vertex '%s' that is not available in this region (structural integrity violation)", i, transition.Target.ID),
						transitionContext.Path,
					)
				}
			}
		}
	}
}

// validateContainmentRelationships validates containment relationships within the region
func (r *Region) validateContainmentRelationships(context *ValidationContext, errors *ValidationErrors) {
	// Validate that composite states properly contain their regions
	for i, state := range r.States {
		if state == nil || !state.IsComposite {
			continue
		}

		stateContext := context.WithPathIndex("States", i)

		// Check that composite state regions don't conflict with this region's vertices
		for j, subRegion := range state.Regions {
			if subRegion == nil {
				continue
			}

			subRegionContext := stateContext.WithPathIndex("Regions", j)

			// Validate that sub-region vertices don't have ID conflicts with parent region
			for k, subVertex := range subRegion.Vertices {
				if subVertex == nil {
					continue
				}

				// Check if this vertex ID exists in the parent region
				for _, parentVertex := range r.Vertices {
					if parentVertex != nil && parentVertex.ID == subVertex.ID {
						errors.AddError(
							ErrorTypeConstraint,
							"Region",
							"ContainmentRelationships",
							fmt.Sprintf("composite state region contains vertex with ID '%s' that conflicts with parent region vertex (structural integrity violation)", subVertex.ID),
							subRegionContext.WithPathIndex("Vertices", k).Path,
						)
					}
				}
			}
		}
	}
}
