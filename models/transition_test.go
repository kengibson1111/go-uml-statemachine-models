package models

import "testing"

func TestTransitionKind_IsValid(t *testing.T) {
	tests := []struct {
		name string
		tk   TransitionKind
		want bool
	}{
		{"internal", TransitionKindInternal, true},
		{"local", TransitionKindLocal, true},
		{"external", TransitionKindExternal, true},
		{"invalid", TransitionKind("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tk.IsValid(); got != tt.want {
				t.Errorf("TransitionKind.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransition_Validate(t *testing.T) {
	validSource := &Vertex{
		ID:   "source1",
		Name: "SourceState",
		Type: "state",
	}
	validTarget := &Vertex{
		ID:   "target1",
		Name: "TargetState",
		Type: "state",
	}

	tests := []struct {
		name       string
		transition *Transition
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid transition",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Target: validTarget,
				Kind:   TransitionKindExternal,
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			transition: &Transition{
				Name:   "TestTransition",
				Source: validSource,
				Target: validTarget,
				Kind:   TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "[Required] Transition.ID: field is required and cannot be empty",
		},
		{
			name: "nil source",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Target: validTarget,
				Kind:   TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "[Required] Transition.Source: required reference cannot be nil",
		},
		{
			name: "nil target",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Kind:   TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "[Required] Transition.Target: required reference cannot be nil",
		},
		{
			name: "invalid source vertex",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: &Vertex{
					// Missing required fields
				},
				Target: validTarget,
				Kind:   TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
		{
			name: "invalid target vertex",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Target: &Vertex{
					// Missing required fields
				},
				Kind: TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
		{
			name: "invalid kind",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Target: validTarget,
				Kind:   TransitionKind("invalid"),
			},
			wantErr: true,
			errMsg:  "[Invalid] Transition.Kind: invalid TransitionKind: invalid",
		},
		{
			name: "invalid trigger",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Target: validTarget,
				Kind:   TransitionKindExternal,
				Triggers: []*Trigger{
					{
						// Missing required fields
					},
				},
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
		{
			name: "invalid guard",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Target: validTarget,
				Kind:   TransitionKindExternal,
				Guard:  &Constraint{
					// Missing required fields
				},
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
		{
			name: "invalid effect",
			transition: &Transition{
				ID:     "t1",
				Name:   "TestTransition",
				Source: validSource,
				Target: validTarget,
				Kind:   TransitionKindExternal,
				Effect: &Behavior{
					// Missing required fields
				},
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transition.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Transition.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Transition.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Transition.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestTransition_ValidateSourceTarget(t *testing.T) {
	tests := []struct {
		name       string
		transition *Transition
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid state to state transition",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Target: &Vertex{
					ID:   "s2",
					Name: "State2",
					Type: "state",
				},
				Kind: TransitionKindExternal,
			},
			wantErr: false,
		},
		{
			name: "final state as source",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "f1",
					Name: "FinalState",
					Type: "finalstate",
				},
				Target: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Kind: TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "final state cannot be the source of a transition (UML constraint)",
		},
		{
			name: "initial pseudostate as target",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Target: &Vertex{
					ID:   "initial",
					Name: "Initial",
					Type: "pseudostate",
				},
				Kind: TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "initial pseudostate cannot be the target of a transition within the same region (UML constraint)",
		},
		{
			name: "terminate pseudostate as source",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "terminate",
					Name: "Terminate",
					Type: "pseudostate",
				},
				Target: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Kind: TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "terminate pseudostate cannot have outgoing transitions (UML constraint)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := NewValidationContext()
			errors := &ValidationErrors{}
			tt.transition.validateSourceTarget(context, errors)

			if tt.wantErr {
				if !errors.HasErrors() {
					t.Errorf("validateSourceTarget() expected error but got none")
					return
				}
				if tt.errMsg != "" {
					found := false
					for _, err := range errors.Errors {
						if err.Message == tt.errMsg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("validateSourceTarget() error messages = %v, want to contain %v", errors.Error(), tt.errMsg)
					}
				}
			} else {
				if errors.HasErrors() {
					t.Errorf("validateSourceTarget() unexpected error = %v", errors.Error())
				}
			}
		})
	}
}

func TestTransition_ValidateKindConstraints(t *testing.T) {
	tests := []struct {
		name       string
		transition *Transition
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid internal transition",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Target: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Kind: TransitionKindInternal,
			},
			wantErr: false,
		},
		{
			name: "internal transition with different source and target",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Target: &Vertex{
					ID:   "s2",
					Name: "State2",
					Type: "state",
				},
				Kind: TransitionKindInternal,
			},
			wantErr: true,
			errMsg:  "internal transition must have the same source and target vertex (UML constraint)",
		},
		{
			name: "internal transition with pseudostate source",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "junction",
					Name: "Junction",
					Type: "pseudostate",
				},
				Target: &Vertex{
					ID:   "junction",
					Name: "Junction",
					Type: "pseudostate",
				},
				Kind: TransitionKindInternal,
			},
			wantErr: true,
			errMsg:  "internal transition source should be a state (UML constraint)",
		},
		{
			name: "valid external transition",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Target: &Vertex{
					ID:   "s2",
					Name: "State2",
					Type: "state",
				},
				Kind: TransitionKindExternal,
			},
			wantErr: false,
		},
		{
			name: "external self-transition",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Target: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Kind: TransitionKindExternal,
			},
			wantErr: true,
			errMsg:  "self-transition with external kind may cause exit/entry actions to be executed (UML design consideration)",
		},
		{
			name: "local transition with connection point source",
			transition: &Transition{
				ID: "t1",
				Source: &Vertex{
					ID:   "entryPoint",
					Name: "EntryPoint",
					Type: "pseudostate",
				},
				Target: &Vertex{
					ID:   "s1",
					Name: "State1",
					Type: "state",
				},
				Kind: TransitionKindLocal,
			},
			wantErr: true,
			errMsg:  "local transition should not originate from connection points (UML constraint)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := NewValidationContext()
			errors := &ValidationErrors{}
			tt.transition.validateKindConstraints(context, errors)

			if tt.wantErr {
				if !errors.HasErrors() {
					t.Errorf("validateKindConstraints() expected error but got none")
					return
				}
				if tt.errMsg != "" {
					found := false
					for _, err := range errors.Errors {
						if err.Message == tt.errMsg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("validateKindConstraints() error messages = %v, want to contain %v", errors.Error(), tt.errMsg)
					}
				}
			} else {
				if errors.HasErrors() {
					t.Errorf("validateKindConstraints() unexpected error = %v", errors.Error())
				}
			}
		})
	}
}

func TestTransition_ValidateContainment(t *testing.T) {
	// Create test vertices
	sourceVertex := &Vertex{
		ID:   "s1",
		Name: "State1",
		Type: "state",
	}
	targetVertex := &Vertex{
		ID:   "s2",
		Name: "State2",
		Type: "state",
	}
	externalVertex := &Vertex{
		ID:   "s3",
		Name: "State3",
		Type: "state",
	}

	// Create test region with vertices
	region := &Region{
		ID:       "r1",
		Name:     "Region1",
		Vertices: []*Vertex{sourceVertex, targetVertex},
	}

	tests := []struct {
		name       string
		transition *Transition
		region     *Region
		wantErr    bool
		errMsg     string
	}{
		{
			name: "valid containment",
			transition: &Transition{
				ID:     "t1",
				Source: sourceVertex,
				Target: targetVertex,
				Kind:   TransitionKindExternal,
			},
			region:  region,
			wantErr: false,
		},
		{
			name: "no region context",
			transition: &Transition{
				ID:     "t1",
				Source: sourceVertex,
				Target: targetVertex,
				Kind:   TransitionKindExternal,
			},
			region:  nil,
			wantErr: false, // No error when no region context - validation is lenient
		},
		{
			name: "source not in region",
			transition: &Transition{
				ID:     "t1",
				Source: externalVertex,
				Target: targetVertex,
				Kind:   TransitionKindExternal,
			},
			region:  region,
			wantErr: true,
			errMsg:  "Source vertex (ID: s3) is not contained in the transition's region (UML constraint)",
		},
		{
			name: "internal transition target not in region",
			transition: &Transition{
				ID:     "t1",
				Source: sourceVertex,
				Target: externalVertex,
				Kind:   TransitionKindInternal,
			},
			region:  region,
			wantErr: true,
			errMsg:  "Target vertex (ID: s3) is not contained in the transition's region (UML constraint)",
		},
		{
			name: "external transition target not in region (no state machine context)",
			transition: &Transition{
				ID:     "t1",
				Source: sourceVertex,
				Target: externalVertex,
				Kind:   TransitionKindExternal,
			},
			region:  region,
			wantErr: false, // No error when no state machine context
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := NewValidationContext()
			if tt.region != nil {
				context = context.WithRegion(tt.region)
			}
			errors := &ValidationErrors{}
			tt.transition.validateContainment(context, errors)

			if tt.wantErr {
				if !errors.HasErrors() {
					t.Errorf("validateContainment() expected error but got none")
					return
				}
				if tt.errMsg != "" {
					found := false
					for _, err := range errors.Errors {
						if err.Message == tt.errMsg {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("validateContainment() error messages = %v, want to contain %v", errors.Error(), tt.errMsg)
					}
				}
			} else {
				if errors.HasErrors() {
					t.Errorf("validateContainment() unexpected error = %v", errors.Error())
				}
			}
		})
	}
}

func TestTransition_PseudostateHelpers(t *testing.T) {
	transition := &Transition{}

	tests := []struct {
		name     string
		vertex   *Vertex
		helper   func(*Vertex) bool
		expected bool
	}{
		{
			name: "initial pseudostate by name",
			vertex: &Vertex{
				ID:   "ps1",
				Name: "Initial",
				Type: "pseudostate",
			},
			helper:   transition.isInitialPseudostate,
			expected: true,
		},
		{
			name: "initial pseudostate by id",
			vertex: &Vertex{
				ID:   "initial",
				Name: "StartState",
				Type: "pseudostate",
			},
			helper:   transition.isInitialPseudostate,
			expected: true,
		},
		{
			name: "not initial pseudostate",
			vertex: &Vertex{
				ID:   "ps1",
				Name: "SomeState",
				Type: "pseudostate",
			},
			helper:   transition.isInitialPseudostate,
			expected: false,
		},
		{
			name: "terminate pseudostate",
			vertex: &Vertex{
				ID:   "terminate",
				Name: "Terminate",
				Type: "pseudostate",
			},
			helper:   transition.isTerminatePseudostate,
			expected: true,
		},
		{
			name: "history pseudostate",
			vertex: &Vertex{
				ID:   "h1",
				Name: "History",
				Type: "pseudostate",
			},
			helper:   transition.isHistoryPseudostate,
			expected: true,
		},
		{
			name: "junction pseudostate",
			vertex: &Vertex{
				ID:   "j1",
				Name: "Junction",
				Type: "pseudostate",
			},
			helper:   transition.isJunctionOrChoice,
			expected: true,
		},
		{
			name: "choice pseudostate",
			vertex: &Vertex{
				ID:   "c1",
				Name: "Choice",
				Type: "pseudostate",
			},
			helper:   transition.isJunctionOrChoice,
			expected: true,
		},
		{
			name: "entry point",
			vertex: &Vertex{
				ID:   "ep1",
				Name: "EntryPoint",
				Type: "pseudostate",
			},
			helper:   transition.isConnectionPoint,
			expected: true,
		},
		{
			name: "exit point",
			vertex: &Vertex{
				ID:   "exit",
				Name: "ExitPoint",
				Type: "pseudostate",
			},
			helper:   transition.isConnectionPoint,
			expected: true,
		},
		{
			name: "regular state",
			vertex: &Vertex{
				ID:   "s1",
				Name: "State1",
				Type: "state",
			},
			helper:   transition.isInitialPseudostate,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.helper(tt.vertex)
			if result != tt.expected {
				t.Errorf("helper function returned %v, want %v for vertex %+v", result, tt.expected, tt.vertex)
			}
		})
	}
}
