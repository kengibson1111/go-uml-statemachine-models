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
			errMsg:  "Vertex ID cannot be empty",
		},
		{
			name: "empty Name",
			vertex: &Vertex{
				ID:   "v1",
				Type: "state",
			},
			wantErr: true,
			errMsg:  "Vertex Name cannot be empty",
		},
		{
			name: "empty Type",
			vertex: &Vertex{
				ID:   "v1",
				Name: "TestVertex",
			},
			wantErr: true,
			errMsg:  "Vertex Type cannot be empty",
		},
		{
			name: "invalid Type",
			vertex: &Vertex{
				ID:   "v1",
				Name: "TestVertex",
				Type: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid Vertex Type: invalid, must be one of: state, pseudostate, finalstate",
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
			errMsg:  "State must have type 'state', got: pseudostate",
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
			errMsg:  "invalid entry behavior: Behavior ID cannot be empty",
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
					Name: "TestPseudostate",
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
			errMsg:  "Pseudostate must have type 'pseudostate', got: state",
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
			errMsg:  "invalid PseudostateKind: invalid",
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
			errMsg:  "FinalState must have type 'finalstate', got: state",
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
			errMsg:  "invalid entry pseudostate at index 0: invalid vertex in pseudostate: Vertex ID cannot be empty",
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
				wantErr: true,
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
				wantErr: true,
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
