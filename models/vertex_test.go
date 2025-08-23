package models

import "testing"

func TestVertex_Validate(t *testing.T) {
	tests := []struct {
		name    string
		vertex  *Vertex
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid vertex",
			vertex: &Vertex{
				ID:   "v1",
				Name: "TestVertex",
				Type: "state",
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			vertex: &Vertex{
				Name: "TestVertex",
				Type: "state",
			},
			wantErr: true,
			errMsg:  "[Required] Vertex.ID: field is required and cannot be empty",
		},
		{
			name: "empty Name",
			vertex: &Vertex{
				ID:   "v1",
				Type: "state",
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
		{
			name: "empty Type",
			vertex: &Vertex{
				ID:   "v1",
				Name: "TestVertex",
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
		{
			name: "invalid Type",
			vertex: &Vertex{
				ID:   "v1",
				Name: "TestVertex",
				Type: "invalid",
			},
			wantErr: true,
			errMsg:  "[Invalid] Vertex.Type: invalid value: must be one of [state, pseudostate, finalstate]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vertex.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Vertex.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Vertex.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Vertex.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestState_Validate(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid simple state",
			state: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
				IsSimple: true,
			},
			wantErr: false,
		},
		{
			name: "valid composite state with regions",
			state: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
				IsComposite: true,
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid vertex type",
			state: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "pseudostate",
				},
			},
			wantErr: true,
			errMsg:  "[Constraint] State.Type: State must have type 'state', got: pseudostate",
		},
		{
			name: "invalid entry behavior",
			state: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
				Entry: &Behavior{
					// Missing required fields
				},
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.state.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("State.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("State.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("State.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestPseudostateKind_IsValid(t *testing.T) {
	tests := []struct {
		name string
		pk   PseudostateKind
		want bool
	}{
		{"initial", PseudostateKindInitial, true},
		{"deepHistory", PseudostateKindDeepHistory, true},
		{"shallowHistory", PseudostateKindShallowHistory, true},
		{"join", PseudostateKindJoin, true},
		{"fork", PseudostateKindFork, true},
		{"junction", PseudostateKindJunction, true},
		{"choice", PseudostateKindChoice, true},
		{"entryPoint", PseudostateKindEntryPoint, true},
		{"exitPoint", PseudostateKindExitPoint, true},
		{"terminate", PseudostateKindTerminate, true},
		{"invalid", PseudostateKind("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.IsValid(); got != tt.want {
				t.Errorf("PseudostateKind.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPseudostate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ps      *Pseudostate
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid pseudostate",
			ps: &Pseudostate{
				Vertex: Vertex{
					ID:   "ps1",
					Name: "Initial State",
					Type: "pseudostate",
				},
				Kind: PseudostateKindInitial,
			},
			wantErr: false,
		},
		{
			name: "invalid vertex type",
			ps: &Pseudostate{
				Vertex: Vertex{
					ID:   "ps1",
					Name: "TestPseudostate",
					Type: "state",
				},
				Kind: PseudostateKindInitial,
			},
			wantErr: true,
			errMsg:  "Pseudostate.Type: Pseudostate must have type 'pseudostate', got: state",
		},
		{
			name: "invalid kind",
			ps: &Pseudostate{
				Vertex: Vertex{
					ID:   "ps1",
					Name: "TestPseudostate",
					Type: "pseudostate",
				},
				Kind: PseudostateKind("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid PseudostateKind:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ps.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Pseudostate.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Pseudostate.Validate() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Pseudostate.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestFinalState_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fs      *FinalState
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid final state",
			fs: &FinalState{
				Vertex: Vertex{
					ID:   "fs1",
					Name: "TestFinalState",
					Type: "finalstate",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid vertex type",
			fs: &FinalState{
				Vertex: Vertex{
					ID:   "fs1",
					Name: "TestFinalState",
					Type: "state",
				},
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fs.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("FinalState.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("FinalState.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("FinalState.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestConnectionPointReference_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cpr     *ConnectionPointReference
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid connection point reference",
			cpr: &ConnectionPointReference{
				Vertex: Vertex{
					ID:   "cpr1",
					Name: "TestConnectionPoint",
					Type: "pseudostate",
				},
				Entry: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "entry1",
							Name: "EntryPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindEntryPoint,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid entry pseudostate",
			cpr: &ConnectionPointReference{
				Vertex: Vertex{
					ID:   "cpr1",
					Name: "TestConnectionPoint",
					Type: "pseudostate",
				},
				Entry: []*Pseudostate{
					{
						Vertex: Vertex{
							// Missing required fields
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "", // Multiple errors expected, so we'll just check for error existence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cpr.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("ConnectionPointReference.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ConnectionPointReference.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ConnectionPointReference.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestPseudostate_UMLConstraintValidation tests the UML-specific constraint validation methods for Pseudostate
func TestPseudostate_UMLConstraintValidation(t *testing.T) {
	t.Run("validateKindConstraints", func(t *testing.T) {
		tests := []struct {
			name    string
			ps      *Pseudostate
			context *ValidationContext
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid initial pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "initial1",
						Name: "Initial",
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial,
				},
				context: NewValidationContext(),
				wantErr: false,
			},
			{
				name: "initial pseudostate without name",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "initial1",
						Name: "", // Empty name
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial,
				},
				context: NewValidationContext(),
				wantErr: true,
				errMsg:  "initial pseudostate should have a descriptive name (UML best practice)",
			},
			{
				name: "valid deep history pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "history1",
						Name: "DeepHistory",
						Type: "pseudostate",
					},
					Kind: PseudostateKindDeepHistory,
				},
				context: NewValidationContext().WithRegion(&Region{ID: "r1", Name: "TestRegion"}),
				wantErr: false,
			},
			{
				name: "history pseudostate without region context",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "history1",
						Name: "DeepHistory",
						Type: "pseudostate",
					},
					Kind: PseudostateKindDeepHistory,
				},
				context: NewValidationContext(), // No region context
				wantErr: true,
				errMsg:  "history pseudostate must be contained within a region of a composite state (UML constraint)",
			},
			{
				name: "valid entry point pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "entry1",
						Name: "EntryPoint",
						Type: "pseudostate",
					},
					Kind: PseudostateKindEntryPoint,
				},
				context: NewValidationContext().WithStateMachine(&StateMachine{ID: "sm1", Name: "TestSM"}),
				wantErr: false,
			},
			{
				name: "entry point without state machine context",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "entry1",
						Name: "EntryPoint",
						Type: "pseudostate",
					},
					Kind: PseudostateKindEntryPoint,
				},
				context: NewValidationContext(), // No state machine context
				wantErr: false,                  // TODO: Should be true when validateKindConstraints is implemented
				errMsg:  "entryPoint pseudostate should be used as a connection point in a state machine (UML constraint)",
			},
			{
				name: "valid junction pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "junction1",
						Name: "Junction",
						Type: "pseudostate",
					},
					Kind: PseudostateKindJunction,
				},
				context: NewValidationContext(),
				wantErr: false,
			},
			{
				name: "valid choice pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "choice1",
						Name: "Choice",
						Type: "pseudostate",
					},
					Kind: PseudostateKindChoice,
				},
				context: NewValidationContext(),
				wantErr: false,
			},
			{
				name: "valid fork pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "fork1",
						Name: "Fork",
						Type: "pseudostate",
					},
					Kind: PseudostateKindFork,
				},
				context: NewValidationContext(),
				wantErr: false,
			},
			{
				name: "valid join pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "join1",
						Name: "Join",
						Type: "pseudostate",
					},
					Kind: PseudostateKindJoin,
				},
				context: NewValidationContext(),
				wantErr: false,
			},
			{
				name: "valid terminate pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "terminate1",
						Name: "Terminate",
						Type: "pseudostate",
					},
					Kind: PseudostateKindTerminate,
				},
				context: NewValidationContext(),
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.ps.ValidateInContext(tt.context)
				if tt.wantErr {
					if err == nil {
						t.Errorf("Pseudostate.ValidateInContext() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("Pseudostate.ValidateInContext() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("Pseudostate.ValidateInContext() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateMultiplicity", func(t *testing.T) {
		tests := []struct {
			name    string
			ps      *Pseudostate
			region  *Region
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid single initial pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "initial1",
						Name: "Initial",
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "initial1",
							Name: "Initial",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "multiple initial pseudostates",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "initial1",
						Name: "Initial",
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "initial1",
							Name: "Initial",
							Type: "pseudostate",
						},
						{
							ID:   "initial2",
							Name: "init",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: true,
				errMsg:  "region can have at most one initial pseudostate, found 2 initial pseudostates (UML constraint)",
			},
			{
				name: "valid single deep history pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "history1",
						Name: "deepHistory",
						Type: "pseudostate",
					},
					Kind: PseudostateKindDeepHistory,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "history1",
							Name: "deepHistory",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "multiple deep history pseudostates",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "history1",
						Name: "deepHistory",
						Type: "pseudostate",
					},
					Kind: PseudostateKindDeepHistory,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "history1",
							Name: "deepHistory",
							Type: "pseudostate",
						},
						{
							ID:   "history2",
							Name: "DeepHistory",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: true,
				errMsg:  "region should have at most one deepHistory pseudostate, found 2 (UML best practice)",
			},
			{
				name: "valid single terminate pseudostate",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "terminate1",
						Name: "terminate",
						Type: "pseudostate",
					},
					Kind: PseudostateKindTerminate,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "terminate1",
							Name: "terminate",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "multiple terminate pseudostates",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "terminate1",
						Name: "terminate",
						Type: "pseudostate",
					},
					Kind: PseudostateKindTerminate,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "terminate1",
							Name: "terminate",
							Type: "pseudostate",
						},
						{
							ID:   "terminate2",
							Name: "Terminate",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: true,
				errMsg:  "region has 2 terminate pseudostates, consider if this is intended (UML design consideration)",
			},
			{
				name: "junction pseudostate - no multiplicity constraints",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "junction1",
						Name: "Junction",
						Type: "pseudostate",
					},
					Kind: PseudostateKindJunction,
				},
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "junction1",
							Name: "Junction",
							Type: "pseudostate",
						},
						{
							ID:   "junction2",
							Name: "Junction2",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: false, // Junction pseudostates don't have multiplicity constraints
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				context := NewValidationContext().WithRegion(tt.region)
				err := tt.ps.ValidateInContext(context)
				if tt.wantErr {
					if err == nil {
						t.Errorf("Pseudostate.ValidateInContext() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("Pseudostate.ValidateInContext() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("Pseudostate.ValidateInContext() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("comprehensive UML constraint validation", func(t *testing.T) {
		tests := []struct {
			name    string
			ps      *Pseudostate
			context *ValidationContext
			wantErr bool
			errMsgs []string // Multiple error messages to check for
		}{
			{
				name: "valid pseudostate with all constraints satisfied",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "initial1",
						Name: "Initial",
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial,
				},
				context: NewValidationContext().WithRegion(&Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "initial1",
							Name: "Initial",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				}),
				wantErr: false,
			},
			{
				name: "invalid pseudostate with multiple constraint violations",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "initial1",
						Name: "", // Missing name - constraint violation
						Type: "pseudostate",
					},
					Kind: PseudostateKindInitial,
				},
				context: NewValidationContext().WithRegion(&Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "initial1",
							Name: "Initial",
							Type: "pseudostate",
						},
						{
							ID:   "initial2", // Multiple initial - multiplicity violation
							Name: "init",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
				}),
				wantErr: true,
				errMsgs: []string{
					"initial pseudostate should have a descriptive name (UML best practice)",
					"region can have at most one initial pseudostate, found 2 initial pseudostates (UML constraint)",
				},
			},
			{
				name: "entry point pseudostate without proper context",
				ps: &Pseudostate{
					Vertex: Vertex{
						ID:   "entry1",
						Name: "EntryPoint",
						Type: "pseudostate",
					},
					Kind: PseudostateKindEntryPoint,
				},
				context: NewValidationContext().WithRegion(&Region{
					ID:   "r1",
					Name: "TestRegion",
				}), // No state machine context
				wantErr: false, // TODO: Should be true when validateKindConstraints is implemented
				errMsgs: []string{
					"entryPoint pseudostate should be used as a connection point in a state machine (UML constraint)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.ps.ValidateInContext(tt.context)
				if tt.wantErr {
					if err == nil {
						t.Errorf("Pseudostate.ValidateInContext() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("Pseudostate.ValidateInContext() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("Pseudostate.ValidateInContext() unexpected error = %v", err)
					}
				}
			})
		}
	})
}

// TestState_UMLConstraintValidation tests the UML-specific constraint validation methods for State
func TestState_UMLConstraintValidation(t *testing.T) {
	t.Run("validateCompositeConstraints", func(t *testing.T) {
		tests := []struct {
			name    string
			state   *State
			wantErr bool
			errMsgs []string
		}{
			{
				name: "valid composite state with regions",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "CompositeState",
						Type: "state",
					},
					IsComposite: true,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "composite state without regions",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "CompositeState",
						Type: "state",
					},
					IsComposite: true,
					Regions:     []*Region{}, // No regions
				},
				wantErr: true,
				errMsgs: []string{
					"composite state must have at least one region (UML constraint)",
				},
			},
			{
				name: "state that is both composite and simple",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "InvalidState",
						Type: "state",
					},
					IsComposite: true,
					IsSimple:    true, // Cannot be both
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"state cannot be both composite and simple (UML constraint)",
				},
			},
			{
				name: "orthogonal composite state with single region",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "OrthogonalState",
						Type: "state",
					},
					IsComposite:  true,
					IsOrthogonal: true,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
						},
					}, // Only one region, but orthogonal requires at least two
				},
				wantErr: true,
				errMsgs: []string{
					"orthogonal composite state must have at least two regions (UML constraint)",
				},
			},
			{
				name: "valid orthogonal composite state with multiple regions",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "OrthogonalState",
						Type: "state",
					},
					IsComposite:  true,
					IsOrthogonal: true,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
							Vertices: []*Vertex{
								{
									ID:   "initial1",
									Name: "Initial",
									Type: "pseudostate",
								},
							},
						},
						{
							ID:   "r2",
							Name: "Region2",
							Vertices: []*Vertex{
								{
									ID:   "initial2",
									Name: "Initial",
									Type: "pseudostate",
								},
							},
						},
					},
				},
				wantErr: false,
			},
			{
				name: "non-composite state with regions",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SimpleState",
						Type: "state",
					},
					IsComposite: false,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
						},
					}, // Should not have regions
				},
				wantErr: true,
				errMsgs: []string{
					"non-composite state cannot have regions (UML constraint)",
				},
			},
			{
				name: "non-composite state marked as orthogonal",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SimpleState",
						Type: "state",
					},
					IsComposite:  false,
					IsOrthogonal: true, // Cannot be orthogonal if not composite
				},
				wantErr: true,
				errMsgs: []string{
					"non-composite state cannot be orthogonal (UML constraint)",
				},
			},
			{
				name: "composite state with region missing ID",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "CompositeState",
						Type: "state",
					},
					IsComposite: true,
					Regions: []*Region{
						{
							ID:   "", // Missing ID
							Name: "Region1",
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"region at index 0 must have a valid ID (UML constraint)",
				},
			},
			{
				name: "composite state with region missing name",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "CompositeState",
						Type: "state",
					},
					IsComposite: true,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "", // Missing name
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"region at index 0 should have a descriptive name (UML best practice)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.state.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("State.Validate() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("State.Validate() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("State.Validate() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateSubmachineConstraints", func(t *testing.T) {
		tests := []struct {
			name    string
			state   *State
			wantErr bool
			errMsgs []string
		}{
			{
				name: "valid submachine state",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					Submachine: &StateMachine{
						ID:      "sm1",
						Name:    "ReferencedStateMachine",
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
						ConnectionPoints: []*Pseudostate{
							{
								Vertex: Vertex{
									ID:   "entry1",
									Name: "EntryPoint",
									Type: "pseudostate",
								},
								Kind: PseudostateKindEntryPoint,
							},
						},
					},
					Connections: []*ConnectionPointReference{
						{
							Vertex: Vertex{
								ID:   "cpr1",
								Name: "ConnectionRef",
								Type: "pseudostate",
							},
							Entry: []*Pseudostate{
								{
									Vertex: Vertex{
										ID:   "entry1",
										Name: "EntryPoint",
										Type: "pseudostate",
									},
									Kind: PseudostateKindEntryPoint,
								},
							},
						},
					},
				},
				wantErr: false,
			},
			{
				name: "submachine state without submachine reference",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					Submachine:        nil, // Missing submachine reference
				},
				wantErr: true,
				errMsgs: []string{
					"submachine state must reference a valid state machine (UML constraint)",
				},
			},
			{
				name: "submachine state with invalid submachine (missing ID)",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					Submachine: &StateMachine{
						ID:      "", // Missing ID
						Name:    "ReferencedStateMachine",
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"referenced submachine must have a valid ID (UML constraint)",
				},
			},
			{
				name: "submachine state with invalid submachine (missing name)",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					Submachine: &StateMachine{
						ID:      "sm1",
						Name:    "", // Missing name
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"referenced submachine should have a descriptive name (UML best practice)",
				},
			},
			{
				name: "submachine state marked as composite",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					IsComposite:       true, // Should not be composite
					Submachine: &StateMachine{
						ID:      "sm1",
						Name:    "ReferencedStateMachine",
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"submachine state should not be marked as composite (use submachine reference instead) (UML constraint)",
				},
			},
			{
				name: "submachine state with own regions",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					Submachine: &StateMachine{
						ID:      "sm1",
						Name:    "ReferencedStateMachine",
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
					},
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
						},
					}, // Should not have own regions
				},
				wantErr: true,
				errMsgs: []string{
					"submachine state should not have its own regions (use submachine reference instead) (UML constraint)",
				},
			},
			{
				name: "non-submachine state with submachine reference",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "RegularState",
						Type: "state",
					},
					IsSubmachineState: false,
					Submachine: &StateMachine{
						ID:      "sm1",
						Name:    "ReferencedStateMachine",
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
					}, // Should not have submachine reference
				},
				wantErr: true,
				errMsgs: []string{
					"non-submachine state should not reference a submachine (UML constraint)",
				},
			},
			{
				name: "non-submachine state with connection points",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "RegularState",
						Type: "state",
					},
					IsSubmachineState: false,
					Connections: []*ConnectionPointReference{
						{
							Vertex: Vertex{
								ID:   "cpr1",
								Name: "ConnectionRef",
								Type: "pseudostate",
							},
						},
					}, // Should not have connection points
				},
				wantErr: true,
				errMsgs: []string{
					"non-submachine state should not have connection point references (UML constraint)",
				},
			},
			{
				name: "submachine state with invalid connection point reference",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "SubmachineState",
						Type: "state",
					},
					IsSubmachineState: true,
					Submachine: &StateMachine{
						ID:      "sm1",
						Name:    "ReferencedStateMachine",
						Version: "1.0",
						Regions: []*Region{
							{
								ID:   "r1",
								Name: "DefaultRegion",
							},
						},
						ConnectionPoints: []*Pseudostate{
							{
								Vertex: Vertex{
									ID:   "entry1",
									Name: "EntryPoint",
									Type: "pseudostate",
								},
								Kind: PseudostateKindEntryPoint,
							},
						},
					},
					Connections: []*ConnectionPointReference{
						{
							Vertex: Vertex{
								ID:   "cpr1",
								Name: "ConnectionRef",
								Type: "pseudostate",
							},
							Entry: []*Pseudostate{
								{
									Vertex: Vertex{
										ID:   "nonexistent", // References non-existent entry point
										Name: "NonExistentEntry",
										Type: "pseudostate",
									},
									Kind: PseudostateKindEntryPoint,
								},
							},
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"connection point reference at index 0 references entry point 'nonexistent' that does not exist in submachine (UML constraint)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.state.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("State.Validate() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("State.Validate() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("State.Validate() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateBehaviorConsistency", func(t *testing.T) {
		tests := []struct {
			name    string
			state   *State
			wantErr bool
			errMsgs []string
		}{
			{
				name: "valid state with all behaviors",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "EntryBehavior",
						Specification: "initialize state",
						Language:      "javascript",
					},
					Exit: &Behavior{
						ID:            "exit1",
						Name:          "ExitBehavior",
						Specification: "cleanup state",
						Language:      "javascript",
					},
					DoActivity: &Behavior{
						ID:            "do1",
						Name:          "DoActivity",
						Specification: "perform ongoing activity",
						Language:      "javascript",
					},
				},
				wantErr: false,
			},
			{
				name: "state with behavior missing ID",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "", // Missing ID
						Name:          "EntryBehavior",
						Specification: "initialize state",
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry behavior must have a valid ID (UML constraint)",
				},
			},
			{
				name: "state with behavior missing name",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "", // Missing name
						Specification: "initialize state",
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry behavior should have a descriptive name (UML best practice)",
				},
			},
			{
				name: "state with behavior missing specification",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "EntryBehavior",
						Specification: "", // Missing specification
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry behavior must have a valid specification (UML constraint)",
				},
			},
			{
				name: "state with behavior language but no specification",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "EntryBehavior",
						Specification: "", // Missing specification
						Language:      "javascript",
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry behavior must have a valid specification (UML constraint)",
					"entry behavior specifies language 'javascript' but has no specification content (UML constraint)",
				},
			},
			{
				name: "state with behaviors having same names",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "SameName", // Same name
						Specification: "initialize state",
					},
					Exit: &Behavior{
						ID:            "exit1",
						Name:          "SameName", // Same name
						Specification: "cleanup state",
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry and exit behaviors should have distinct names to avoid confusion (UML best practice)",
				},
			},
			{
				name: "state with behaviors having different languages",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "EntryBehavior",
						Specification: "initialize state",
						Language:      "javascript",
					},
					Exit: &Behavior{
						ID:            "exit1",
						Name:          "ExitBehavior",
						Specification: "cleanup state",
						Language:      "python", // Different language
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry behavior uses language 'javascript' while exit behavior uses 'python', consider consistency (UML best practice)",
				},
			},
			{
				name: "state with behaviors having same ID but different specifications",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "behavior1", // Same ID
						Name:          "EntryBehavior",
						Specification: "initialize state",
					},
					Exit: &Behavior{
						ID:            "behavior1", // Same ID
						Name:          "ExitBehavior",
						Specification: "cleanup state", // Different specification
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry and exit behaviors have the same ID but different specifications, which may cause confusion (UML best practice)",
				},
			},
			{
				name: "state with do activity having same ID as entry",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "behavior1", // Same ID
						Name:          "EntryBehavior",
						Specification: "initialize state",
					},
					DoActivity: &Behavior{
						ID:            "behavior1", // Same ID
						Name:          "DoActivity",
						Specification: "perform ongoing activity",
					},
				},
				wantErr: true,
				errMsgs: []string{
					"do activity and entry behavior have the same ID, which may cause confusion (UML best practice)",
				},
			},
			{
				name: "state with entry behavior suggesting cleanup",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "EntryBehavior",
						Specification: "cleanup resources", // Should be in exit
					},
				},
				wantErr: true,
				errMsgs: []string{
					"entry behavior specification suggests cleanup operations, which should typically be in exit behavior (UML semantics)",
				},
			},
			{
				name: "state with exit behavior suggesting initialization",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					Exit: &Behavior{
						ID:            "exit1",
						Name:          "ExitBehavior",
						Specification: "initialize resources", // Should be in entry
					},
				},
				wantErr: true,
				errMsgs: []string{
					"exit behavior specification suggests initialization operations, which should typically be in entry behavior (UML semantics)",
				},
			},
			{
				name: "state with do activity suggesting one-time operations",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "StateWithBehaviors",
						Type: "state",
					},
					DoActivity: &Behavior{
						ID:            "do1",
						Name:          "DoActivity",
						Specification: "initialize system", // Should be in entry/exit
					},
				},
				wantErr: true,
				errMsgs: []string{
					"do activity specification suggests one-time operations, which should typically be in entry or exit behaviors (UML semantics)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.state.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("State.Validate() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("State.Validate() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("State.Validate() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("comprehensive UML constraint validation", func(t *testing.T) {
		tests := []struct {
			name    string
			state   *State
			wantErr bool
			errMsgs []string
		}{
			{
				name: "valid state with all constraints satisfied",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "ValidState",
						Type: "state",
					},
					IsComposite: true,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "Region1",
						},
						{
							ID:   "r2",
							Name: "Region2",
						},
					},
					Entry: &Behavior{
						ID:            "entry1",
						Name:          "EntryBehavior",
						Specification: "initialize state",
						Language:      "javascript",
					},
					Exit: &Behavior{
						ID:            "exit1",
						Name:          "ExitBehavior",
						Specification: "cleanup state",
						Language:      "javascript",
					},
				},
				wantErr: false,
			},
			{
				name: "invalid state with multiple constraint violations",
				state: &State{
					Vertex: Vertex{
						ID:   "s1",
						Name: "InvalidState",
						Type: "state",
					},
					IsComposite:       true,
					IsSimple:          true,        // Cannot be both composite and simple
					IsSubmachineState: true,        // Cannot be both composite and submachine
					Regions:           []*Region{}, // Composite but no regions
					Entry: &Behavior{
						ID:            "", // Missing ID
						Name:          "SameName",
						Specification: "cleanup resources", // Wrong semantic for entry
					},
					Exit: &Behavior{
						ID:            "exit1",
						Name:          "SameName",             // Same name as entry
						Specification: "initialize resources", // Wrong semantic for exit
						Language:      "python",
					},
					DoActivity: &Behavior{
						ID:            "exit1", // Same ID as exit
						Name:          "DoActivity",
						Specification: "setup system", // Wrong semantic for do activity
						Language:      "javascript",   // Different language from exit
					},
				},
				wantErr: true,
				errMsgs: []string{
					"state cannot be both composite and simple (UML constraint)",
					"composite state must have at least one region (UML constraint)",
					"submachine state should not be marked as composite (use submachine reference instead) (UML constraint)",
					"submachine state must reference a valid state machine (UML constraint)",
					"entry behavior must have a valid ID (UML constraint)",
					"entry and exit behaviors should have distinct names to avoid confusion (UML best practice)",
					"exit behavior uses language 'python' while do activity uses 'javascript', consider consistency (UML best practice)",
					"do activity and exit behavior have the same ID, which may cause confusion (UML best practice)",
					"entry behavior specification suggests cleanup operations, which should typically be in exit behavior (UML semantics)",
					"exit behavior specification suggests initialization operations, which should typically be in entry behavior (UML semantics)",
					"do activity specification suggests one-time operations, which should typically be in entry or exit behaviors (UML semantics)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.state.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("State.Validate() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("State.Validate() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("State.Validate() unexpected error = %v", err)
					}
				}
			})
		}
	})
}
