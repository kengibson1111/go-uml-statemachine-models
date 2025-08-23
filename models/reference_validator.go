package models

import (
	"fmt"
	"reflect"
)

// ReferenceValidator provides structural integrity validation for UML state machine models
type ReferenceValidator struct {
	visited           map[string]bool
	errors            *ValidationErrors
	context           *ValidationContext
	referenceMap      map[string]interface{} // Maps object IDs to their instances
	bidirectionalRefs map[string][]string    // Maps object IDs to their bidirectional references
	containmentTree   map[string][]string    // Maps parent IDs to child IDs
	inheritanceTree   map[string]string      // Maps child IDs to parent IDs
}

// NewReferenceValidator creates a new reference validator
func NewReferenceValidator() *ReferenceValidator {
	return &ReferenceValidator{
		visited:           make(map[string]bool),
		errors:            &ValidationErrors{},
		referenceMap:      make(map[string]interface{}),
		bidirectionalRefs: make(map[string][]string),
		containmentTree:   make(map[string][]string),
		inheritanceTree:   make(map[string]string),
	}
}

// ValidateReferences performs comprehensive reference validation on a state machine
func (rv *ReferenceValidator) ValidateReferences(obj interface{}) error {
	rv.context = NewValidationContext()

	// First pass: build reference maps
	rv.buildReferenceMaps(obj, rv.context)

	// Second pass: validate references
	rv.validateObjectReferences(obj, rv.context)

	// Third pass: validate structural integrity
	rv.validateStructuralIntegrity()

	return rv.errors.ToError()
}

// ValidateReferencesInContext performs reference validation with provided context
func (rv *ReferenceValidator) ValidateReferencesInContext(obj interface{}, context *ValidationContext) error {
	rv.context = context
	if rv.context == nil {
		rv.context = NewValidationContext()
	}

	// First pass: build reference maps
	rv.buildReferenceMaps(obj, rv.context)

	// Second pass: validate references
	rv.validateObjectReferences(obj, rv.context)

	// Third pass: validate structural integrity
	rv.validateStructuralIntegrity()

	return rv.errors.ToError()
}

// buildReferenceMaps builds internal maps of object references for validation
func (rv *ReferenceValidator) buildReferenceMaps(obj interface{}, context *ValidationContext) {
	if obj == nil {
		return
	}

	objID := rv.getObjectID(obj)
	if objID == "" {
		return
	}

	// Check for duplicate IDs (different objects with same ID)
	if existingObj, exists := rv.referenceMap[objID]; exists {
		if existingObj != obj {
			// Different objects with same ID - this is a structural integrity issue
			rv.errors.AddError(
				ErrorTypeConstraint,
				rv.getObjectTypeName(obj),
				"DuplicateID",
				fmt.Sprintf("duplicate ID '%s' found in different objects", objID),
				context.Path,
			)
		}
		return // Don't process the same object twice
	}

	// Store object reference
	rv.referenceMap[objID] = obj

	// Avoid infinite recursion for containment
	if rv.visited[objID] {
		return
	}
	rv.visited[objID] = true

	// Build maps based on object type
	switch v := obj.(type) {
	case *StateMachine:
		rv.buildStateMachineReferences(v, context)
	case *Region:
		rv.buildRegionReferences(v, context)
	case *State:
		rv.buildStateReferences(v, context)
	case *Pseudostate:
		rv.buildPseudostateReferences(v, context)
	case *FinalState:
		rv.buildFinalStateReferences(v, context)
	case *Transition:
		rv.buildTransitionReferences(v, context)
	case *ConnectionPointReference:
		rv.buildConnectionPointReferences(v, context)
	}
}

// buildStateMachineReferences builds reference maps for StateMachine
func (rv *ReferenceValidator) buildStateMachineReferences(sm *StateMachine, context *ValidationContext) {
	smContext := context.WithStateMachine(sm)

	// Build containment relationships for regions
	for i, region := range sm.Regions {
		if region != nil {
			rv.containmentTree[sm.ID] = append(rv.containmentTree[sm.ID], region.ID)
			rv.buildReferenceMaps(region, smContext.WithPathIndex("Regions", i))
		}
	}

	// Build references for connection points
	for i, cp := range sm.ConnectionPoints {
		if cp != nil {
			rv.containmentTree[sm.ID] = append(rv.containmentTree[sm.ID], cp.ID)
			rv.buildReferenceMaps(cp, smContext.WithPathIndex("ConnectionPoints", i))
		}
	}
}

// buildRegionReferences builds reference maps for Region
func (rv *ReferenceValidator) buildRegionReferences(region *Region, context *ValidationContext) {
	regionContext := context.WithRegion(region)

	// Build containment relationships for states
	for i, state := range region.States {
		if state != nil {
			rv.containmentTree[region.ID] = append(rv.containmentTree[region.ID], state.ID)
			rv.buildReferenceMaps(state, regionContext.WithPathIndex("States", i))
		}
	}

	// Build containment relationships for vertices
	for i, vertex := range region.Vertices {
		if vertex != nil {
			rv.containmentTree[region.ID] = append(rv.containmentTree[region.ID], vertex.ID)
			rv.buildReferenceMaps(vertex, regionContext.WithPathIndex("Vertices", i))
		}
	}

	// Build references for transitions
	for i, transition := range region.Transitions {
		if transition != nil {
			rv.containmentTree[region.ID] = append(rv.containmentTree[region.ID], transition.ID)
			rv.buildReferenceMaps(transition, regionContext.WithPathIndex("Transitions", i))
		}
	}
}

// buildStateReferences builds reference maps for State
func (rv *ReferenceValidator) buildStateReferences(state *State, context *ValidationContext) {
	stateContext := context.WithPath("State")

	// Build containment relationships for regions in composite states
	for i, region := range state.Regions {
		if region != nil {
			rv.containmentTree[state.ID] = append(rv.containmentTree[state.ID], region.ID)
			rv.buildReferenceMaps(region, stateContext.WithPathIndex("Regions", i))
		}
	}

	// Build references for submachine (inheritance-like relationship)
	if state.Submachine != nil {
		rv.inheritanceTree[state.ID] = state.Submachine.ID
		rv.buildReferenceMaps(state.Submachine, stateContext.WithPath("Submachine"))
	}

	// Build references for connection point references
	for i, conn := range state.Connections {
		if conn != nil {
			rv.containmentTree[state.ID] = append(rv.containmentTree[state.ID], conn.ID)
			rv.buildReferenceMaps(conn, stateContext.WithPathIndex("Connections", i))
		}
	}

	// Build references for behaviors
	if state.Entry != nil {
		rv.buildReferenceMaps(state.Entry, stateContext.WithPath("Entry"))
	}
	if state.Exit != nil {
		rv.buildReferenceMaps(state.Exit, stateContext.WithPath("Exit"))
	}
	if state.DoActivity != nil {
		rv.buildReferenceMaps(state.DoActivity, stateContext.WithPath("DoActivity"))
	}
}

// buildPseudostateReferences builds reference maps for Pseudostate
func (rv *ReferenceValidator) buildPseudostateReferences(ps *Pseudostate, context *ValidationContext) {
	// Pseudostates don't typically contain other objects, but we track them for reference validation
}

// buildFinalStateReferences builds reference maps for FinalState
func (rv *ReferenceValidator) buildFinalStateReferences(fs *FinalState, context *ValidationContext) {
	// Final states don't typically contain other objects, but we track them for reference validation
}

// buildTransitionReferences builds reference maps for Transition
func (rv *ReferenceValidator) buildTransitionReferences(transition *Transition, context *ValidationContext) {
	transitionContext := context.WithPath("Transition")

	// Build bidirectional references for source and target
	if transition.Source != nil {
		rv.bidirectionalRefs[transition.ID] = append(rv.bidirectionalRefs[transition.ID], transition.Source.ID)
		rv.bidirectionalRefs[transition.Source.ID] = append(rv.bidirectionalRefs[transition.Source.ID], transition.ID)
	}

	if transition.Target != nil {
		rv.bidirectionalRefs[transition.ID] = append(rv.bidirectionalRefs[transition.ID], transition.Target.ID)
		rv.bidirectionalRefs[transition.Target.ID] = append(rv.bidirectionalRefs[transition.Target.ID], transition.ID)
	}

	// Build references for triggers
	for i, trigger := range transition.Triggers {
		if trigger != nil {
			rv.buildReferenceMaps(trigger, transitionContext.WithPathIndex("Triggers", i))
		}
	}

	// Build references for guard and effect
	if transition.Guard != nil {
		rv.buildReferenceMaps(transition.Guard, transitionContext.WithPath("Guard"))
	}
	if transition.Effect != nil {
		rv.buildReferenceMaps(transition.Effect, transitionContext.WithPath("Effect"))
	}
}

// buildConnectionPointReferences builds reference maps for ConnectionPointReference
func (rv *ReferenceValidator) buildConnectionPointReferences(cpr *ConnectionPointReference, context *ValidationContext) {
	cprContext := context.WithPath("ConnectionPointReference")

	// Build references for entry and exit pseudostates
	for i, entry := range cpr.Entry {
		if entry != nil {
			rv.bidirectionalRefs[cpr.ID] = append(rv.bidirectionalRefs[cpr.ID], entry.ID)
			rv.bidirectionalRefs[entry.ID] = append(rv.bidirectionalRefs[entry.ID], cpr.ID)
			rv.buildReferenceMaps(entry, cprContext.WithPathIndex("Entry", i))
		}
	}

	for i, exit := range cpr.Exit {
		if exit != nil {
			rv.bidirectionalRefs[cpr.ID] = append(rv.bidirectionalRefs[cpr.ID], exit.ID)
			rv.bidirectionalRefs[exit.ID] = append(rv.bidirectionalRefs[exit.ID], cpr.ID)
			rv.buildReferenceMaps(exit, cprContext.WithPathIndex("Exit", i))
		}
	}
}

// validateObjectReferences validates individual object references
func (rv *ReferenceValidator) validateObjectReferences(obj interface{}, context *ValidationContext) {
	if obj == nil {
		return
	}

	objID := rv.getObjectID(obj)
	if objID == "" {
		return
	}

	// Validate based on object type
	switch v := obj.(type) {
	case *StateMachine:
		rv.validateStateMachineReferences(v, context)
	case *Region:
		rv.validateRegionReferences(v, context)
	case *State:
		rv.validateStateReferences(v, context)
	case *Transition:
		rv.validateTransitionReferences(v, context)
	case *ConnectionPointReference:
		rv.validateConnectionPointReferenceReferences(v, context)
	}
}

// validateStateMachineReferences validates StateMachine references
func (rv *ReferenceValidator) validateStateMachineReferences(sm *StateMachine, context *ValidationContext) {
	smContext := context.WithStateMachine(sm)

	// Validate region references
	for i, region := range sm.Regions {
		if region == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"StateMachine",
				"Regions",
				fmt.Sprintf("region at index %d is nil", i),
				smContext.WithPathIndex("Regions", i).Path,
			)
			continue
		}

		// Validate region exists in reference map
		if _, exists := rv.referenceMap[region.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"StateMachine",
				"Regions",
				fmt.Sprintf("region at index %d (ID: %s) not found in reference map", i, region.ID),
				smContext.WithPathIndex("Regions", i).Path,
			)
		}

		rv.validateObjectReferences(region, smContext.WithPathIndex("Regions", i))
	}

	// Validate connection point references
	for i, cp := range sm.ConnectionPoints {
		if cp == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"StateMachine",
				"ConnectionPoints",
				fmt.Sprintf("connection point at index %d is nil", i),
				smContext.WithPathIndex("ConnectionPoints", i).Path,
			)
			continue
		}

		// Validate connection point exists in reference map
		if _, exists := rv.referenceMap[cp.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"StateMachine",
				"ConnectionPoints",
				fmt.Sprintf("connection point at index %d (ID: %s) not found in reference map", i, cp.ID),
				smContext.WithPathIndex("ConnectionPoints", i).Path,
			)
		}

		rv.validateObjectReferences(cp, smContext.WithPathIndex("ConnectionPoints", i))
	}
}

// validateRegionReferences validates Region references
func (rv *ReferenceValidator) validateRegionReferences(region *Region, context *ValidationContext) {
	regionContext := context.WithRegion(region)

	// Validate state references
	for i, state := range region.States {
		if state == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"Region",
				"States",
				fmt.Sprintf("state at index %d is nil", i),
				regionContext.WithPathIndex("States", i).Path,
			)
			continue
		}

		rv.validateObjectReferences(state, regionContext.WithPathIndex("States", i))
	}

	// Validate vertex references
	for i, vertex := range region.Vertices {
		if vertex == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"Region",
				"Vertices",
				fmt.Sprintf("vertex at index %d is nil", i),
				regionContext.WithPathIndex("Vertices", i).Path,
			)
			continue
		}

		rv.validateObjectReferences(vertex, regionContext.WithPathIndex("Vertices", i))
	}

	// Validate transition references
	for i, transition := range region.Transitions {
		if transition == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"Region",
				"Transitions",
				fmt.Sprintf("transition at index %d is nil", i),
				regionContext.WithPathIndex("Transitions", i).Path,
			)
			continue
		}

		rv.validateObjectReferences(transition, regionContext.WithPathIndex("Transitions", i))
	}
}

// validateStateReferences validates State references
func (rv *ReferenceValidator) validateStateReferences(state *State, context *ValidationContext) {
	stateContext := context.WithPath("State")

	// Validate submachine reference
	if state.Submachine != nil {
		if _, exists := rv.referenceMap[state.Submachine.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"State",
				"Submachine",
				fmt.Sprintf("submachine (ID: %s) not found in reference map", state.Submachine.ID),
				stateContext.WithPath("Submachine").Path,
			)
		}

		rv.validateObjectReferences(state.Submachine, stateContext.WithPath("Submachine"))
	}

	// Validate region references in composite states
	for i, region := range state.Regions {
		if region == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"State",
				"Regions",
				fmt.Sprintf("region at index %d is nil", i),
				stateContext.WithPathIndex("Regions", i).Path,
			)
			continue
		}

		rv.validateObjectReferences(region, stateContext.WithPathIndex("Regions", i))
	}

	// Validate connection point references
	for i, conn := range state.Connections {
		if conn == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"State",
				"Connections",
				fmt.Sprintf("connection point reference at index %d is nil", i),
				stateContext.WithPathIndex("Connections", i).Path,
			)
			continue
		}

		rv.validateObjectReferences(conn, stateContext.WithPathIndex("Connections", i))
	}
}

// validateTransitionReferences validates Transition references
func (rv *ReferenceValidator) validateTransitionReferences(transition *Transition, context *ValidationContext) {
	transitionContext := context.WithPath("Transition")

	// Validate source reference
	if transition.Source == nil {
		rv.errors.AddError(
			ErrorTypeReference,
			"Transition",
			"Source",
			"source vertex is required and cannot be nil",
			transitionContext.Path,
		)
	} else {
		if _, exists := rv.referenceMap[transition.Source.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"Transition",
				"Source",
				fmt.Sprintf("source vertex (ID: %s) not found in reference map", transition.Source.ID),
				transitionContext.WithPath("Source").Path,
			)
		}
	}

	// Validate target reference
	if transition.Target == nil {
		rv.errors.AddError(
			ErrorTypeReference,
			"Transition",
			"Target",
			"target vertex is required and cannot be nil",
			transitionContext.Path,
		)
	} else {
		if _, exists := rv.referenceMap[transition.Target.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"Transition",
				"Target",
				fmt.Sprintf("target vertex (ID: %s) not found in reference map", transition.Target.ID),
				transitionContext.WithPath("Target").Path,
			)
		}
	}
}

// validateConnectionPointReferenceReferences validates ConnectionPointReference references
func (rv *ReferenceValidator) validateConnectionPointReferenceReferences(cpr *ConnectionPointReference, context *ValidationContext) {
	cprContext := context.WithPath("ConnectionPointReference")

	// Validate entry pseudostate references
	for i, entry := range cpr.Entry {
		if entry == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"ConnectionPointReference",
				"Entry",
				fmt.Sprintf("entry pseudostate at index %d is nil", i),
				cprContext.WithPathIndex("Entry", i).Path,
			)
			continue
		}

		if _, exists := rv.referenceMap[entry.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"ConnectionPointReference",
				"Entry",
				fmt.Sprintf("entry pseudostate at index %d (ID: %s) not found in reference map", i, entry.ID),
				cprContext.WithPathIndex("Entry", i).Path,
			)
		}
	}

	// Validate exit pseudostate references
	for i, exit := range cpr.Exit {
		if exit == nil {
			rv.errors.AddError(
				ErrorTypeReference,
				"ConnectionPointReference",
				"Exit",
				fmt.Sprintf("exit pseudostate at index %d is nil", i),
				cprContext.WithPathIndex("Exit", i).Path,
			)
			continue
		}

		if _, exists := rv.referenceMap[exit.ID]; !exists {
			rv.errors.AddError(
				ErrorTypeReference,
				"ConnectionPointReference",
				"Exit",
				fmt.Sprintf("exit pseudostate at index %d (ID: %s) not found in reference map", i, exit.ID),
				cprContext.WithPathIndex("Exit", i).Path,
			)
		}
	}
}

// validateStructuralIntegrity performs high-level structural integrity validation
func (rv *ReferenceValidator) validateStructuralIntegrity() {
	// Validate bidirectional relationship consistency
	rv.validateBidirectionalConsistency()

	// Validate containment hierarchy
	rv.validateContainmentHierarchy()

	// Validate inheritance relationships and detect cycles
	rv.validateInheritanceRelationships()
}

// validateBidirectionalConsistency validates that bidirectional relationships are consistent
func (rv *ReferenceValidator) validateBidirectionalConsistency() {
	for objID, refs := range rv.bidirectionalRefs {
		obj, exists := rv.referenceMap[objID]
		if !exists {
			continue
		}

		objContext := NewValidationContext().WithPath(fmt.Sprintf("Object[%s]", objID))

		// Validate that each referenced object also references back
		for _, refID := range refs {
			refObj, refExists := rv.referenceMap[refID]
			if !refExists {
				rv.errors.AddError(
					ErrorTypeReference,
					rv.getObjectTypeName(obj),
					"BidirectionalReference",
					fmt.Sprintf("object (ID: %s) references non-existent object (ID: %s)", objID, refID),
					objContext.Path,
				)
				continue
			}

			// Check if the referenced object references back
			refRefs, hasRefs := rv.bidirectionalRefs[refID]
			if !hasRefs {
				rv.errors.AddError(
					ErrorTypeConstraint,
					rv.getObjectTypeName(obj),
					"BidirectionalReference",
					fmt.Sprintf("object (ID: %s) references object (ID: %s) but the reference is not bidirectional", objID, refID),
					objContext.Path,
				)
				continue
			}

			// Check if objID is in the referenced object's references
			found := false
			for _, backRef := range refRefs {
				if backRef == objID {
					found = true
					break
				}
			}

			if !found {
				rv.errors.AddError(
					ErrorTypeConstraint,
					rv.getObjectTypeName(obj),
					"BidirectionalReference",
					fmt.Sprintf("object (ID: %s) references object (ID: %s) but the reference is not properly bidirectional", objID, refID),
					objContext.Path,
				)
			}

			// Validate reference type compatibility
			rv.validateReferenceTypeCompatibility(obj, refObj, objContext)
		}
	}
}

// validateContainmentHierarchy validates the containment hierarchy for consistency
func (rv *ReferenceValidator) validateContainmentHierarchy() {
	for parentID, childIDs := range rv.containmentTree {
		parentObj, exists := rv.referenceMap[parentID]
		if !exists {
			continue
		}

		parentContext := NewValidationContext().WithPath(fmt.Sprintf("Parent[%s]", parentID))

		// Validate each child exists and is properly contained
		for _, childID := range childIDs {
			childObj, childExists := rv.referenceMap[childID]
			if !childExists {
				rv.errors.AddError(
					ErrorTypeReference,
					rv.getObjectTypeName(parentObj),
					"ContainmentHierarchy",
					fmt.Sprintf("parent (ID: %s) contains non-existent child (ID: %s)", parentID, childID),
					parentContext.Path,
				)
				continue
			}

			// Validate containment type compatibility
			rv.validateContainmentTypeCompatibility(parentObj, childObj, parentContext)

			// Check for containment cycles
			rv.validateContainmentCycle(parentID, childID, parentContext)
		}
	}
}

// validateInheritanceRelationships validates inheritance relationships and detects cycles
func (rv *ReferenceValidator) validateInheritanceRelationships() {
	for childID, parentID := range rv.inheritanceTree {
		childObj, childExists := rv.referenceMap[childID]
		if !childExists {
			continue
		}

		parentObj, parentExists := rv.referenceMap[parentID]
		if !parentExists {
			childContext := NewValidationContext().WithPath(fmt.Sprintf("Child[%s]", childID))
			rv.errors.AddError(
				ErrorTypeReference,
				rv.getObjectTypeName(childObj),
				"InheritanceRelationship",
				fmt.Sprintf("child (ID: %s) inherits from non-existent parent (ID: %s)", childID, parentID),
				childContext.Path,
			)
			continue
		}

		// Validate inheritance type compatibility
		childContext := NewValidationContext().WithPath(fmt.Sprintf("Child[%s]", childID))
		rv.validateInheritanceTypeCompatibility(childObj, parentObj, childContext)

		// Check for inheritance cycles
		rv.validateInheritanceCycle(childID, parentID, childContext)
	}
}

// validateReferenceTypeCompatibility validates that referenced objects are of compatible types
func (rv *ReferenceValidator) validateReferenceTypeCompatibility(obj, refObj interface{}, context *ValidationContext) {
	objType := rv.getObjectTypeName(obj)
	refType := rv.getObjectTypeName(refObj)

	// Define incompatible reference types - be more permissive for bidirectional references
	// Only validate strict incompatibilities, not all possible references
	incompatibleRefs := map[string][]string{
		"Vertex":      {"StateMachine", "Region"}, // Vertices shouldn't directly reference containers
		"Pseudostate": {"StateMachine", "Region"}, // Pseudostates shouldn't directly reference containers
		"FinalState":  {"StateMachine", "Region"}, // Final states shouldn't directly reference containers
	}

	if disallowedTypes, exists := incompatibleRefs[objType]; exists {
		for _, disallowedType := range disallowedTypes {
			if refType == disallowedType {
				rv.errors.AddError(
					ErrorTypeConstraint,
					objType,
					"ReferenceTypeCompatibility",
					fmt.Sprintf("%s should not directly reference %s (type incompatibility)", objType, refType),
					context.Path,
				)
				break
			}
		}
	}
}

// validateContainmentTypeCompatibility validates that containment relationships are type-compatible
func (rv *ReferenceValidator) validateContainmentTypeCompatibility(parent, child interface{}, context *ValidationContext) {
	parentType := rv.getObjectTypeName(parent)
	childType := rv.getObjectTypeName(child)

	// Define valid containment relationships
	validContainments := map[string][]string{
		"StateMachine":             {"Region", "Pseudostate"},
		"Region":                   {"State", "Vertex", "Pseudostate", "FinalState", "Transition"},
		"State":                    {"Region", "ConnectionPointReference"},
		"ConnectionPointReference": {"Pseudostate"},
	}

	if allowedChildren, exists := validContainments[parentType]; exists {
		valid := false
		for _, allowedChild := range allowedChildren {
			if childType == allowedChild {
				valid = true
				break
			}
		}

		if !valid {
			rv.errors.AddError(
				ErrorTypeConstraint,
				parentType,
				"ContainmentTypeCompatibility",
				fmt.Sprintf("%s cannot contain %s (containment rule violation)", parentType, childType),
				context.Path,
			)
		}
	}
}

// validateInheritanceTypeCompatibility validates that inheritance relationships are type-compatible
func (rv *ReferenceValidator) validateInheritanceTypeCompatibility(child, parent interface{}, context *ValidationContext) {
	childType := rv.getObjectTypeName(child)
	parentType := rv.getObjectTypeName(parent)

	// Define valid inheritance relationships (for UML state machines, mainly submachine relationships)
	validInheritances := map[string][]string{
		"State": {"StateMachine"}, // State can reference a submachine
	}

	if allowedParents, exists := validInheritances[childType]; exists {
		valid := false
		for _, allowedParent := range allowedParents {
			if parentType == allowedParent {
				valid = true
				break
			}
		}

		if !valid {
			rv.errors.AddError(
				ErrorTypeConstraint,
				childType,
				"InheritanceTypeCompatibility",
				fmt.Sprintf("%s cannot inherit from %s (inheritance rule violation)", childType, parentType),
				context.Path,
			)
		}
	}
}

// validateContainmentCycle checks for cycles in the containment hierarchy
func (rv *ReferenceValidator) validateContainmentCycle(parentID, childID string, context *ValidationContext) {
	visited := make(map[string]bool)
	path := []string{parentID}

	if rv.hasContainmentCycleWithDepthLimit(childID, parentID, visited, path, 0, 100) {
		rv.errors.AddError(
			ErrorTypeConstraint,
			"ContainmentHierarchy",
			"CycleDetection",
			fmt.Sprintf("containment cycle detected: %v", path),
			context.Path,
		)
	}
}

// validateInheritanceCycle checks for cycles in inheritance relationships
func (rv *ReferenceValidator) validateInheritanceCycle(childID, parentID string, context *ValidationContext) {
	// Direct self-reference is an immediate cycle
	if childID == parentID {
		rv.errors.AddError(
			ErrorTypeConstraint,
			"InheritanceHierarchy",
			"CycleDetection",
			fmt.Sprintf("inheritance cycle detected: direct self-reference %s -> %s", childID, parentID),
			context.Path,
		)
		return
	}

	visited := make(map[string]bool)
	path := []string{childID}

	if rv.hasInheritanceCycleWithDepthLimit(parentID, childID, visited, path, 0, 100) {
		rv.errors.AddError(
			ErrorTypeConstraint,
			"InheritanceHierarchy",
			"CycleDetection",
			fmt.Sprintf("inheritance cycle detected: %v", path),
			context.Path,
		)
	}
}

// hasContainmentCycle recursively checks for containment cycles
func (rv *ReferenceValidator) hasContainmentCycle(currentID, targetID string, visited map[string]bool, path []string) bool {
	return rv.hasContainmentCycleWithDepthLimit(currentID, targetID, visited, path, 0, 100)
}

// hasContainmentCycleWithDepthLimit recursively checks for containment cycles with depth limit
func (rv *ReferenceValidator) hasContainmentCycleWithDepthLimit(currentID, targetID string, visited map[string]bool, path []string, depth, maxDepth int) bool {
	if depth > maxDepth {
		// Prevent infinite recursion by limiting depth
		return false
	}

	if currentID == targetID {
		return true
	}

	if visited[currentID] {
		return false
	}

	visited[currentID] = true
	// Create a new path slice to avoid modifying the original
	newPath := make([]string, len(path), len(path)+1)
	copy(newPath, path)
	newPath = append(newPath, currentID)

	if children, exists := rv.containmentTree[currentID]; exists {
		for _, childID := range children {
			if rv.hasContainmentCycleWithDepthLimit(childID, targetID, visited, newPath, depth+1, maxDepth) {
				return true
			}
		}
	}

	// Clean up visited state for this branch
	visited[currentID] = false
	return false
}

// hasInheritanceCycle recursively checks for inheritance cycles
func (rv *ReferenceValidator) hasInheritanceCycle(currentID, targetID string, visited map[string]bool, path []string) bool {
	return rv.hasInheritanceCycleWithDepthLimit(currentID, targetID, visited, path, 0, 100)
}

// hasInheritanceCycleWithDepthLimit recursively checks for inheritance cycles with depth limit
func (rv *ReferenceValidator) hasInheritanceCycleWithDepthLimit(currentID, targetID string, visited map[string]bool, path []string, depth, maxDepth int) bool {
	if depth > maxDepth {
		// Prevent infinite recursion by limiting depth
		return false
	}

	if currentID == targetID {
		return true
	}

	if visited[currentID] {
		return false
	}

	visited[currentID] = true
	// Create a new path slice to avoid modifying the original
	newPath := make([]string, len(path), len(path)+1)
	copy(newPath, path)
	newPath = append(newPath, currentID)

	if parentID, exists := rv.inheritanceTree[currentID]; exists {
		if rv.hasInheritanceCycleWithDepthLimit(parentID, targetID, visited, newPath, depth+1, maxDepth) {
			return true
		}
	}

	// Clean up visited state for this branch
	visited[currentID] = false
	return false
}

// getObjectID extracts the ID from an object
func (rv *ReferenceValidator) getObjectID(obj interface{}) string {
	if obj == nil {
		return ""
	}

	// Use reflection to get the ID field
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	idField := v.FieldByName("ID")
	if !idField.IsValid() || idField.Kind() != reflect.String {
		return ""
	}

	return idField.String()
}

// getObjectTypeName returns the type name of an object
func (rv *ReferenceValidator) getObjectTypeName(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

// ValidateContainment validates containment relationships for a specific object
func (rv *ReferenceValidator) ValidateContainment(parent, child interface{}) error {
	if parent == nil || child == nil {
		return fmt.Errorf("parent and child cannot be nil")
	}

	parentID := rv.getObjectID(parent)
	childID := rv.getObjectID(child)

	if parentID == "" || childID == "" {
		return fmt.Errorf("parent and child must have valid IDs")
	}

	// Reset validator state for this specific validation
	rv.errors = &ValidationErrors{}
	rv.context = NewValidationContext()

	// Validate containment type compatibility
	rv.validateContainmentTypeCompatibility(parent, child, rv.context)

	// Check for containment cycles if we have containment tree data
	if len(rv.containmentTree) > 0 {
		rv.validateContainmentCycle(parentID, childID, rv.context)
	}

	return rv.errors.ToError()
}

// DebugInfo returns debug information about the validator state
func (rv *ReferenceValidator) DebugInfo() map[string]interface{} {
	return map[string]interface{}{
		"referenceMap":      rv.referenceMap,
		"inheritanceTree":   rv.inheritanceTree,
		"containmentTree":   rv.containmentTree,
		"bidirectionalRefs": rv.bidirectionalRefs,
		"visited":           rv.visited,
	}
}
