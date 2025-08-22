package models

import (
	"fmt"
	"strings"
)

// ComplexPatternValidator provides validation for complex UML state machine patterns
type ComplexPatternValidator struct {
	context *ValidationContext
	errors  []ValidationError
}

// NewComplexPatternValidator creates a new complex pattern validator
func NewComplexPatternValidator(context *ValidationContext) *ComplexPatternValidator {
	return &ComplexPatternValidator{
		context: context,
		errors:  make([]ValidationError, 0),
	}
}

// ValidateOrthogonalRegions validates orthogonal regions in composite states
func (cpv *ComplexPatternValidator) ValidateOrthogonalRegions(state *State) error {
	if state == nil {
		return fmt.Errorf("state cannot be nil")
	}

	// Only composite states can have orthogonal regions
	if len(state.Regions) <= 1 {
		return nil // Not orthogonal, no validation needed
	}

	// Validate that all regions are properly orthogonal
	for i, region := range state.Regions {
		if region == nil {
			cpv.addError(ValidationError{
				Type:    ErrorTypeRequired,
				Object:  fmt.Sprintf("State[%s]", state.Name),
				Field:   fmt.Sprintf("Regions[%d]", i),
				Message: "orthogonal region cannot be nil",
				Path:    cpv.buildPath(state.Name, fmt.Sprintf("regions[%d]", i)),
			})
			continue
		}

		// Validate orthogonal regions don't share vertices
		if err := cpv.validateRegionSeparation(state.Regions, i); err != nil {
			cpv.addError(ValidationError{
				Type:    ErrorTypeConstraint,
				Object:  fmt.Sprintf("State[%s]", state.Name),
				Field:   "OrthogonalRegions",
				Message: err.Error(),
				Path:    cpv.buildPath(state.Name, "orthogonal_regions"),
			})
		}
	}

	return cpv.collectErrors()
}

// ValidateConnectionPointReferences validates connection point references and submachine interfaces
func (cpv *ComplexPatternValidator) ValidateConnectionPointReferences(stateMachine *StateMachine) error {
	if stateMachine == nil {
		return fmt.Errorf("state machine cannot be nil")
	}

	// Validate connection points are properly typed
	for i, cp := range stateMachine.ConnectionPoints {
		if cp == nil {
			cpv.addError(ValidationError{
				Type:    ErrorTypeRequired,
				Object:  fmt.Sprintf("StateMachine[%s]", stateMachine.Name),
				Field:   fmt.Sprintf("ConnectionPoints[%d]", i),
				Message: "connection point cannot be nil",
				Path:    cpv.buildPath(stateMachine.Name, fmt.Sprintf("connection_points[%d]", i)),
			})
			continue
		}

		// Connection points must be entry or exit pseudostates
		if cp.Kind != PseudostateKindEntryPoint && cp.Kind != PseudostateKindExitPoint {
			cpv.addError(ValidationError{
				Type:    ErrorTypeConstraint,
				Object:  fmt.Sprintf("Pseudostate[%s]", cp.Name),
				Field:   "Kind",
				Message: fmt.Sprintf("connection point must be entry or exit point, got %v", cp.Kind),
				Path:    cpv.buildPath(stateMachine.Name, cp.Name, "kind"),
			})
		}

		// Validate connection point references
		if err := cpv.validateConnectionPointReferences(cp, stateMachine); err != nil {
			cpv.addError(ValidationError{
				Type:    ErrorTypeReference,
				Object:  fmt.Sprintf("Pseudostate[%s]", cp.Name),
				Field:   "References",
				Message: err.Error(),
				Path:    cpv.buildPath(stateMachine.Name, cp.Name, "references"),
			})
		}
	}

	// Validate submachine states reference connection points correctly
	for _, region := range stateMachine.Regions {
		if err := cpv.validateSubmachineConnectionPoints(region); err != nil {
			cpv.addError(ValidationError{
				Type:    ErrorTypeConstraint,
				Object:  fmt.Sprintf("Region[%s]", region.Name),
				Field:   "SubmachineStates",
				Message: err.Error(),
				Path:    cpv.buildPath(stateMachine.Name, region.Name, "submachine_states"),
			})
		}
	}

	return cpv.collectErrors()
}

// ValidateStateMachineInheritance validates state machine inheritance and redefinition
func (cpv *ComplexPatternValidator) ValidateStateMachineInheritance(stateMachine *StateMachine) error {
	if stateMachine == nil {
		return fmt.Errorf("state machine cannot be nil")
	}

	// Check for inheritance cycles
	if err := cpv.detectInheritanceCycles(stateMachine); err != nil {
		cpv.addError(ValidationError{
			Type:    ErrorTypeConstraint,
			Object:  fmt.Sprintf("StateMachine[%s]", stateMachine.Name),
			Field:   "Inheritance",
			Message: err.Error(),
			Path:    cpv.buildPath(stateMachine.Name, "inheritance"),
		})
	}

	// Validate redefinition constraints
	if err := cpv.validateRedefinitionConstraints(stateMachine); err != nil {
		cpv.addError(ValidationError{
			Type:    ErrorTypeConstraint,
			Object:  fmt.Sprintf("StateMachine[%s]", stateMachine.Name),
			Field:   "Redefinition",
			Message: err.Error(),
			Path:    cpv.buildPath(stateMachine.Name, "redefinition"),
		})
	}

	// Validate extended state machine compatibility
	if err := cpv.validateExtendedStateMachineCompatibility(stateMachine); err != nil {
		cpv.addError(ValidationError{
			Type:    ErrorTypeConstraint,
			Object:  fmt.Sprintf("StateMachine[%s]", stateMachine.Name),
			Field:   "Extension",
			Message: err.Error(),
			Path:    cpv.buildPath(stateMachine.Name, "extension"),
		})
	}

	return cpv.collectErrors()
}

// Helper methods

func (cpv *ComplexPatternValidator) validateRegionSeparation(regions []*Region, currentIndex int) error {
	currentRegion := regions[currentIndex]
	if currentRegion == nil {
		return fmt.Errorf("region at index %d is nil", currentIndex)
	}

	// Check that vertices in orthogonal regions don't overlap
	currentVertices := make(map[string]*Vertex)
	for _, vertex := range currentRegion.Vertices {
		if vertex != nil {
			currentVertices[vertex.Name] = vertex
		}
	}

	// Compare with other regions
	for i, otherRegion := range regions {
		if i == currentIndex || otherRegion == nil {
			continue
		}

		for _, vertex := range otherRegion.Vertices {
			if vertex != nil {
				if _, exists := currentVertices[vertex.Name]; exists {
					return fmt.Errorf("vertex '%s' appears in multiple orthogonal regions", vertex.Name)
				}
			}
		}
	}

	return nil
}

func (cpv *ComplexPatternValidator) validateConnectionPointReferences(cp *Pseudostate, stateMachine *StateMachine) error {
	// Validate that connection points are properly referenced by transitions
	referencedByTransitions := false

	// Check all regions for transitions that reference this connection point
	for _, region := range stateMachine.Regions {
		if region == nil {
			continue
		}

		for _, transition := range region.Transitions {
			if transition == nil {
				continue
			}

			// Check if transition references this connection point by comparing IDs
			sourceMatches := transition.Source != nil && transition.Source.ID == cp.ID
			targetMatches := transition.Target != nil && transition.Target.ID == cp.ID

			if sourceMatches || targetMatches {
				referencedByTransitions = true

				// Validate transition kind for connection points
				if cp.Kind == PseudostateKindEntryPoint && !targetMatches {
					return fmt.Errorf("entry point '%s' can only be transition target", cp.Name)
				}
				if cp.Kind == PseudostateKindExitPoint && !sourceMatches {
					return fmt.Errorf("exit point '%s' can only be transition source", cp.Name)
				}
			}
		}
	}

	// Connection points should be referenced by at least one transition
	if !referencedByTransitions {
		return fmt.Errorf("connection point '%s' is not referenced by any transitions", cp.Name)
	}

	return nil
}

func (cpv *ComplexPatternValidator) validateSubmachineConnectionPoints(region *Region) error {
	if region == nil {
		return nil
	}

	// Check states in the region for submachine states
	for _, state := range region.States {
		if state != nil && state.Submachine != nil {
			// Validate that submachine connection points are properly mapped
			if err := cpv.validateSubmachineMapping(state); err != nil {
				return fmt.Errorf("submachine state '%s': %v", state.Name, err)
			}
		}
	}

	return nil
}

func (cpv *ComplexPatternValidator) validateSubmachineMapping(state *State) error {
	if state.Submachine == nil {
		return nil
	}

	// Validate that all connection points in submachine are accessible
	for _, cp := range state.Submachine.ConnectionPoints {
		if cp == nil {
			continue
		}

		// Check if connection point is properly mapped to state's interface
		if !cpv.isConnectionPointMapped(cp, state) {
			return fmt.Errorf("connection point '%s' in submachine is not properly mapped", cp.Name)
		}
	}

	return nil
}

func (cpv *ComplexPatternValidator) isConnectionPointMapped(cp *Pseudostate, state *State) bool {
	// In a real implementation, this would check the state's connection point mappings
	// For now, we'll assume proper mapping if the connection point exists
	return cp != nil && state != nil
}

func (cpv *ComplexPatternValidator) detectInheritanceCycles(stateMachine *StateMachine) error {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	return cpv.detectCyclesRecursive(stateMachine, visited, recursionStack)
}

func (cpv *ComplexPatternValidator) detectCyclesRecursive(sm *StateMachine, visited, recursionStack map[string]bool) error {
	if sm == nil {
		return nil
	}

	smKey := sm.Name
	if recursionStack[smKey] {
		return fmt.Errorf("inheritance cycle detected involving state machine '%s'", sm.Name)
	}

	if visited[smKey] {
		return nil
	}

	visited[smKey] = true
	recursionStack[smKey] = true

	// Check extended state machines (inheritance relationships)
	// In a real implementation, this would traverse the inheritance hierarchy
	// For now, we'll check if there are any obvious cycles in the structure

	recursionStack[smKey] = false
	return nil
}

func (cpv *ComplexPatternValidator) validateRedefinitionConstraints(stateMachine *StateMachine) error {
	// Validate that redefined elements maintain compatibility
	// This is a complex validation that would check:
	// 1. Redefined states maintain their interface
	// 2. Redefined transitions maintain their semantics
	// 3. Added elements don't conflict with inherited ones

	// For now, we'll implement basic checks
	if stateMachine == nil {
		return fmt.Errorf("cannot validate redefinition on nil state machine")
	}

	// Check for name conflicts in redefinition
	elementNames := make(map[string]bool)

	for _, region := range stateMachine.Regions {
		if region == nil {
			continue
		}

		if elementNames[region.Name] {
			return fmt.Errorf("duplicate region name '%s' in redefinition", region.Name)
		}
		elementNames[region.Name] = true

		// Check vertices for conflicts
		for _, vertex := range region.Vertices {
			if vertex == nil {
				continue
			}

			vertexName := vertex.Name
			if elementNames[vertexName] {
				return fmt.Errorf("duplicate vertex name '%s' in redefinition", vertexName)
			}
			elementNames[vertexName] = true
		}
	}

	return nil
}

func (cpv *ComplexPatternValidator) validateExtendedStateMachineCompatibility(stateMachine *StateMachine) error {
	// Validate that extended state machines are compatible
	// This would check interface compatibility, behavioral compatibility, etc.

	if stateMachine == nil {
		return fmt.Errorf("cannot validate extension on nil state machine")
	}

	// Basic compatibility checks
	if len(stateMachine.Regions) == 0 {
		return fmt.Errorf("extended state machine must have at least one region")
	}

	// Validate that connection points are compatible with extension
	for _, cp := range stateMachine.ConnectionPoints {
		if cp == nil {
			continue
		}

		// Extended state machines with connection points must maintain interface compatibility
		if cp.Kind != PseudostateKindEntryPoint && cp.Kind != PseudostateKindExitPoint {
			return fmt.Errorf("extended state machine connection point '%s' must be entry or exit point", cp.Name)
		}
	}

	return nil
}

func (cpv *ComplexPatternValidator) addError(err ValidationError) {
	cpv.errors = append(cpv.errors, err)
}

func (cpv *ComplexPatternValidator) collectErrors() error {
	if len(cpv.errors) == 0 {
		return nil
	}

	var messages []string
	for _, err := range cpv.errors {
		messages = append(messages, fmt.Sprintf("%s.%s: %s", err.Object, err.Field, err.Message))
	}

	return fmt.Errorf("complex pattern validation failed:\n%s", strings.Join(messages, "\n"))
}

func (cpv *ComplexPatternValidator) buildPath(elements ...string) []string {
	path := make([]string, 0, len(elements))
	for _, element := range elements {
		if element != "" {
			path = append(path, element)
		}
	}
	return path
}
