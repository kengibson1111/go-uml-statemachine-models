# Go UML State Machine Models

A comprehensive Go library providing UML-compliant state machine model definitions with robust validation capabilities. This library implements the core data structures and validation logic for UML state machines, serving as the foundation for state machine processing tools and applications.

## Overview

This library provides:

- **Complete UML State Machine Models**: Full implementation of UML 2.5.1 state machine concepts including states, transitions, pseudostates, regions, and behaviors
- **Comprehensive Validation**: Built-in validation engine that enforces UML constraints and structural integrity
- **Type Safety**: Strongly-typed Go structs with proper validation tags and enum types
- **Extensible Architecture**: Clean interfaces and modular design for easy extension and integration

## Key Features

### Core Models

- **StateMachine**: Top-level container with regions, connection points, and metadata
- **Region**: Contains states, transitions, and vertices with UML constraint validation
- **State**: Simple, composite, orthogonal, and submachine states with proper lifecycle behaviors
- **Transition**: Internal, local, and external transitions with triggers, guards, and effects
- **Vertex**: Base type for states and pseudostates with type-safe implementations
- **Pseudostate**: All UML pseudostate kinds (initial, choice, junction, fork, join, history, etc.)
- **Behavior & Constraint**: Action specifications and guard conditions with language support

### Validation Engine

The library includes a sophisticated validation system that enforces:

- **UML Constraints**: All standard UML 2.5.1 state machine constraints
- **Structural Integrity**: Reference consistency and containment validation  
- **Multiplicity Rules**: Proper cardinality enforcement (e.g., at most one initial state per region)
- **Type Safety**: Enum validation and required field checking
- **Cross-Reference Validation**: Ensures transitions reference valid vertices within appropriate scopes

### Validation Features

- **Contextual Validation**: Path-aware error reporting with precise location information
- **Multiple Error Collection**: Comprehensive error reporting that doesn't stop at first failure
- **Hierarchical Validation**: Validates nested structures with proper context propagation
- **Extensible Error Types**: Categorized error types (Required, Invalid, Constraint, Reference, Multiplicity)

## Installation

```bash
go get github.com/kengibson1111/go-uml-statemachine-models
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/kengibson1111/go-uml-statemachine-models/models"
)

func main() {
    // Create a simple state machine
    sm := &models.StateMachine{
        ID:      "traffic-light",
        Name:    "Traffic Light Controller",
        Version: "1.0",
        Regions: []*models.Region{
            {
                ID:   "main-region",
                Name: "Main Region",
                States: []*models.State{
                    {
                        Vertex: models.Vertex{
                            ID:   "red",
                            Name: "Red Light",
                            Type: "state",
                        },
                    },
                    {
                        Vertex: models.Vertex{
                            ID:   "green",
                            Name: "Green Light", 
                            Type: "state",
                        },
                    },
                },
                Transitions: []*models.Transition{
                    {
                        ID:     "red-to-green",
                        Source: &models.Vertex{ID: "red", Name: "Red Light", Type: "state"},
                        Target: &models.Vertex{ID: "green", Name: "Green Light", Type: "state"},
                        Kind:   models.TransitionKindExternal,
                    },
                },
            },
        },
    }

    // Validate the state machine
    if err := sm.Validate(); err != nil {
        log.Printf("Validation failed: %v", err)
    } else {
        fmt.Println("State machine is valid!")
    }
}
```

## Validation Examples

The library provides detailed validation with precise error reporting:

```go
// Invalid state machine - missing required fields
invalidSM := &models.StateMachine{
    // Missing ID, Name, Version
    Regions: []*models.Region{}, // Empty regions violates UML constraint
}

err := invalidSM.Validate()
if err != nil {
    fmt.Println(err.Error())
    // Output: multiple validation errors:
    // [Required] StateMachine.ID: field is required and cannot be empty at 
    // [Required] StateMachine.Name: field is required and cannot be empty at 
    // [Required] StateMachine.Version: field is required and cannot be empty at 
    // [Multiplicity] StateMachine.Regions: StateMachine must have at least one region (UML constraint) at 
}
```

## UML Compliance

This library implements UML 2.5.1 state machine semantics including:

- **Region Constraints**: At least one region per state machine, at most one initial pseudostate per region
- **Connection Points**: Entry/exit points for submachine states with proper validation
- **Transition Kinds**: Internal (no exit/entry), local (within composite state), external (full exit/entry)
- **Pseudostate Rules**: Proper multiplicity and transition constraints for each pseudostate kind
- **Composite States**: Orthogonal regions, submachine references, and proper containment
- **Method Constraints**: State machines used as methods cannot have connection points

## Architecture

The library is organized into several key components:

- **Core Models** (`models/statemachine.go`, `models/vertex.go`, `models/transition.go`): Primary UML entities
- **Validation Framework** (`models/validation.go`): Extensible validation infrastructure  
- **Behavioral Models** (`models/behavior.go`, `models/trigger.go`): Actions, guards, and events
- **Reference Validation** (`models/reference_validator.go`): Cross-reference integrity checking
- **Comprehensive Tests**: Extensive test coverage for all validation scenarios

## Use Cases

This library serves as the foundation for:

- **State Machine Editors**: Providing validated model persistence and manipulation
- **Code Generators**: Ensuring generated code comes from valid UML models
- **Model Transformations**: Converting between different state machine representations
- **Analysis Tools**: Static analysis of state machine properties and behaviors
- **Runtime Engines**: Execution engines that require validated state machine definitions

## Contributing

Contributions are welcome! Please ensure:

1. All new features include comprehensive tests
2. UML compliance is maintained for any model changes  
3. Validation logic includes proper error reporting with context
4. Documentation is updated for new functionality

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
