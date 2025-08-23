package models

import (
	"fmt"
	"log"
)

// ValidationUtilitiesDemo demonstrates the usage of validation utilities and helpers
func ValidationUtilitiesDemo() {
	fmt.Println("=== UML State Machine Validation Utilities Demo ===")

	// Create a sample state machine for demonstration
	sm := createDemoStateMachine()

	// 1. Demonstrate State Machine Traversal
	fmt.Println("1. State Machine Traversal Demo")
	fmt.Println("--------------------------------")
	demonstrateTraversal(sm)

	// 2. Demonstrate Validation Result Aggregation
	fmt.Println("\n2. Validation Result Aggregation Demo")
	fmt.Println("-------------------------------------")
	demonstrateValidationAggregation(sm)

	// 3. Demonstrate Validation Debugging
	fmt.Println("\n3. Validation Debugging Demo")
	fmt.Println("-----------------------------")
	demonstrateValidationDebugging(sm)

	// 4. Demonstrate Common Validation Patterns
	fmt.Println("\n4. Common Validation Patterns Demo")
	fmt.Println("-----------------------------------")
	demonstrateCommonPatterns(sm)

	// 5. Demonstrate Error Scenarios
	fmt.Println("\n5. Error Scenarios Demo")
	fmt.Println("-----------------------")
	demonstrateErrorScenarios()

	fmt.Println("\n=== Demo Complete ===")
}

// createDemoStateMachine creates a comprehensive state machine for demonstration
func createDemoStateMachine() *StateMachine {
	// Create vertices first so we can reference them in transitions
	initialVertex := &Vertex{
		ID:   "initial-pseudostate",
		Name: "Initial",
		Type: "pseudostate",
	}
	finalVertex := &Vertex{
		ID:   "final-state",
		Name: "Final",
		Type: "finalstate",
	}

	idleState := &State{
		Vertex: Vertex{
			ID:   "idle-state",
			Name: "Idle State",
			Type: "state",
		},
		IsSimple: true,
	}

	activeState := &State{
		Vertex: Vertex{
			ID:   "active-state",
			Name: "Active State",
			Type: "state",
		},
		IsSimple: true,
	}

	return &StateMachine{
		ID:      "demo-statemachine",
		Name:    "Demo State Machine",
		Version: "1.0.0",
		Regions: []*Region{
			{
				ID:       "main-region",
				Name:     "Main Region",
				States:   []*State{idleState, activeState},
				Vertices: []*Vertex{initialVertex, finalVertex},
				Transitions: []*Transition{
					{
						ID:     "init-to-idle",
						Name:   "Initialize to Idle",
						Source: initialVertex,
						Target: &idleState.Vertex,
						Kind:   TransitionKindExternal,
					},
					{
						ID:     "idle-to-active",
						Name:   "Activate",
						Source: &idleState.Vertex,
						Target: &activeState.Vertex,
						Kind:   TransitionKindExternal,
					},
				},
			},
		},
		ConnectionPoints: []*Pseudostate{
			{
				Vertex: Vertex{
					ID:   "entry-point",
					Name: "Entry Point",
					Type: "pseudostate",
				},
				Kind: PseudostateKindEntryPoint,
			},
		},
	}
}

// demonstrateTraversal shows how to use the StateMachineTraverser
func demonstrateTraversal(sm *StateMachine) {
	traverser := NewStateMachineTraverser()

	fmt.Printf("Traversing state machine: %s\n", sm.Name)

	objectCount := 0
	maxDepth := 0

	err := traverser.TraverseStateMachine(sm, func(obj interface{}, path []string, depth int) error {
		objectCount++
		if depth > maxDepth {
			maxDepth = depth
		}

		// Print object information
		objType := getObjectTypeName(obj)
		objID := getObjectID(obj)
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}

		fmt.Printf("%s- %s (ID: %s) at depth %d\n", indent, objType, objID, depth)
		return nil
	})

	if err != nil {
		log.Printf("Traversal error: %v", err)
		return
	}

	fmt.Printf("Traversal complete: %d objects found, max depth: %d\n", objectCount, maxDepth)
}

// demonstrateValidationAggregation shows how to use ValidationResultAggregator
func demonstrateValidationAggregation(sm *StateMachine) {
	aggregator := NewValidationResultAggregator()

	fmt.Printf("Validating state machine: %s\n", sm.Name)

	// Validate the state machine and collect results
	err := sm.Validate()
	if err != nil {
		if validationErrors, ok := err.(*ValidationErrors); ok {
			aggregator.AddResult(sm.ID, validationErrors)
		} else {
			// Handle single error case
			aggregator.AddSingleError(sm.ID, ErrorTypeInvalid, "StateMachine", "General", err.Error(), []string{})
		}
	}

	// Validate individual regions
	for i, region := range sm.Regions {
		regionErr := region.Validate()
		if regionErr != nil {
			if validationErrors, ok := regionErr.(*ValidationErrors); ok {
				aggregator.AddResult(fmt.Sprintf("%s.Region[%d]", sm.ID, i), validationErrors)
			}
		}
	}

	// Display results
	if aggregator.HasErrors() {
		fmt.Printf("Validation found %d error(s) across %d object(s)\n",
			aggregator.GetTotalErrorCount(), len(aggregator.GetResults()))

		fmt.Println("\nSummary Report:")
		fmt.Println(aggregator.GetSummaryReport())
	} else {
		fmt.Println("✓ No validation errors found!")
	}
}

// demonstrateValidationDebugging shows how to use ValidationDebugger
func demonstrateValidationDebugging(sm *StateMachine) {
	debugger := NewValidationDebugger()

	fmt.Printf("Debugging state machine: %s\n", sm.Name)

	report, err := debugger.DebugStateMachine(sm)
	if err != nil {
		log.Printf("Debug error: %v", err)
		return
	}

	fmt.Printf("Debug report generated at: %s\n", report.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total objects analyzed: %d\n", report.TotalObjects)
	fmt.Printf("Total validation errors: %d\n", report.TotalErrors)

	// Display summary
	fmt.Println("\nDebug Summary:")
	fmt.Println(report.GetSummary())

	// Show some object details
	fmt.Println("\nObject Details (first 3):")
	count := 0
	for id, obj := range report.Objects {
		if count >= 3 {
			break
		}
		fmt.Printf("- %s (%s) at path: %s\n", id, obj.Type, obj.Path)
		fmt.Printf("  Properties: %v\n", obj.Properties)
		count++
	}
}

// demonstrateCommonPatterns shows how to use CommonValidationPatterns
func demonstrateCommonPatterns(sm *StateMachine) {
	patterns := NewCommonValidationPatterns()
	context := NewValidationContext()
	errors := &ValidationErrors{}

	fmt.Printf("Applying common validation patterns to: %s\n", sm.Name)

	// Validate state machine structure
	patterns.ValidateStateMachineStructure(sm, context, errors)
	fmt.Printf("State machine structure validation: %d error(s)\n", len(errors.Errors))

	// Validate each region structure
	for i, region := range sm.Regions {
		regionErrors := &ValidationErrors{}
		regionContext := context.WithPathIndex("Regions", i)
		patterns.ValidateRegionStructure(region, regionContext, regionErrors)
		fmt.Printf("Region[%d] structure validation: %d error(s)\n", i, len(regionErrors.Errors))
		errors.Merge(regionErrors)
	}

	// Validate object hierarchy
	hierarchyErrors := &ValidationErrors{}
	patterns.ValidateObjectHierarchy(sm, context, hierarchyErrors)
	fmt.Printf("Object hierarchy validation: %d error(s)\n", len(hierarchyErrors.Errors))
	errors.Merge(hierarchyErrors)

	// Display results
	if errors.HasErrors() {
		fmt.Printf("\nTotal validation errors found: %d\n", len(errors.Errors))
		fmt.Println("Detailed errors:")
		for i, err := range errors.Errors {
			fmt.Printf("  %d. %s\n", i+1, err.Error())
		}
	} else {
		fmt.Println("✓ All common validation patterns passed!")
	}
}

// demonstrateErrorScenarios shows validation utilities with intentionally invalid data
func demonstrateErrorScenarios() {
	fmt.Println("Creating intentionally invalid state machine for error demonstration...")

	// Create an invalid state machine
	invalidSM := &StateMachine{
		// Missing required fields: ID, Name, Version
		Regions: []*Region{
			{
				// Missing required fields: ID, Name
				States: []*State{
					{
						Vertex: Vertex{
							// Missing ID and Name
							Type: "state",
						},
					},
				},
				Transitions: []*Transition{
					{
						ID: "invalid-transition",
						// Missing Name, Source, Target
						Kind: TransitionKindInternal,
					},
				},
			},
		},
		IsMethod: true,
		ConnectionPoints: []*Pseudostate{
			{
				Vertex: Vertex{
					ID:   "invalid-cp",
					Name: "Invalid Connection Point",
					Type: "pseudostate",
				},
				Kind: PseudostateKindInitial, // Invalid for connection point
			},
		},
	}

	// Test with validation aggregator
	aggregator := NewValidationResultAggregator()

	err := invalidSM.Validate()
	if err != nil {
		if validationErrors, ok := err.(*ValidationErrors); ok {
			aggregator.AddResult("invalid-sm", validationErrors)
		}
	}

	fmt.Printf("Invalid state machine validation found %d error(s)\n", aggregator.GetTotalErrorCount())

	// Show detailed report
	fmt.Println("\nDetailed Error Report:")
	fmt.Println(aggregator.GetDetailedReport())

	// Test common patterns with invalid data
	patterns := NewCommonValidationPatterns()
	context := NewValidationContext()
	patternErrors := &ValidationErrors{}

	patterns.ValidateStateMachineStructure(invalidSM, context, patternErrors)

	fmt.Printf("\nCommon patterns validation found %d additional error(s)\n", len(patternErrors.Errors))

	// Test debugging with invalid data
	debugger := NewValidationDebugger()
	report, err := debugger.DebugStateMachine(invalidSM)
	if err != nil {
		fmt.Printf("Debug error (expected): %v\n", err)
	} else {
		fmt.Printf("Debug report for invalid SM: %d objects, %d errors\n",
			report.TotalObjects, report.TotalErrors)
	}
}

// Helper functions for demonstration

func getObjectTypeName(obj interface{}) string {
	switch obj.(type) {
	case *StateMachine:
		return "StateMachine"
	case *Region:
		return "Region"
	case *State:
		return "State"
	case *Vertex:
		return "Vertex"
	case *Pseudostate:
		return "Pseudostate"
	case *FinalState:
		return "FinalState"
	case *Transition:
		return "Transition"
	case *ConnectionPointReference:
		return "ConnectionPointReference"
	default:
		return "Unknown"
	}
}

func getObjectID(obj interface{}) string {
	switch v := obj.(type) {
	case *StateMachine:
		return v.ID
	case *Region:
		return v.ID
	case *State:
		return v.ID
	case *Vertex:
		return v.ID
	case *Pseudostate:
		return v.ID
	case *FinalState:
		return v.ID
	case *Transition:
		return v.ID
	case *ConnectionPointReference:
		return v.ID
	default:
		return ""
	}
}
