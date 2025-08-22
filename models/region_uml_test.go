package models

import (
	"testing"
)

// TestRegion_UMLConstraintValidation tests the UML-specific constraint validation methods for Region
func TestRegion_UMLConstraintValidation(t *testing.T) {
	t.Run("validateInitialStates", func(t *testing.T) {
		tests := []struct {
			name    string
			region  *Region
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid - no initial pseudostates",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						{
							ID:   "s2",
							Name: "State2",
							Type: "state",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "valid - one initial pseudostate",
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
				name: "valid - one initial pseudostate with different naming",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "init_state",
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
				wantErr: false,
			},
			{
				name: "invalid - multiple initial pseudostates",
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
							Name: "initial",
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
				errMsg:  "Region can have at most one initial pseudostate, found 2 at indices: [0 1] (UML constraint)",
			},
			{
				name: "invalid - multiple initial pseudostates with mixed naming",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "start_ps",
							Name: "Start",
							Type: "pseudostate",
						},
						{
							ID:   "init_ps",
							Name: "INIT",
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
				errMsg:  "Region can have at most one initial pseudostate, found 2 at indices: [0 1] (UML constraint)",
			},
			{
				name: "valid - pseudostates that are not initial",
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
							ID:   "choice1",
							Name: "Choice",
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
				name: "valid - mixed vertices and states",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
						},
					},
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.region.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("Region.validateInitialStates() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("Region.validateInitialStates() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("Region.validateInitialStates() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateVertexContainment", func(t *testing.T) {
		tests := []struct {
			name    string
			region  *Region
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid - all states in vertices collection",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
						},
						{
							Vertex: Vertex{
								ID:   "s2",
								Name: "State2",
								Type: "state",
							},
						},
					},
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						{
							ID:   "s2",
							Name: "State2",
							Type: "state",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "valid - empty collections",
				region: &Region{
					ID:       "r1",
					Name:     "TestRegion",
					States:   []*State{},
					Vertices: []*Vertex{},
				},
				wantErr: false,
			},
			{
				name: "valid - vertices with valid types",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						{
							ID:   "ps1",
							Name: "Pseudostate1",
							Type: "pseudostate",
						},
						{
							ID:   "fs1",
							Name: "FinalState1",
							Type: "finalstate",
						},
					},
				},
				wantErr: false,
			},
			{
				name: "invalid - state not in vertices collection",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
						},
						{
							Vertex: Vertex{
								ID:   "s2",
								Name: "State2",
								Type: "state",
							},
						},
					},
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						// s2 is missing from vertices
					},
				},
				wantErr: true,
				errMsg:  "state at index 1 (ID: s2) is not contained in the region's vertices collection (UML constraint)",
			},
			{
				name: "invalid - vertex with empty ID",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "", // Empty ID
							Name: "State1",
							Type: "state",
						},
					},
				},
				wantErr: true,
				errMsg:  "vertex at index 0 must have a valid ID for proper containment (UML constraint)",
			},
			{
				name: "invalid - vertex with invalid type",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "v1",
							Name: "InvalidVertex",
							Type: "invalid_type",
						},
					},
				},
				wantErr: true,
				errMsg:  "vertex at index 0 has invalid type 'invalid_type' for region containment (UML constraint)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.region.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("Region.validateVertexContainment() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("Region.validateVertexContainment() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("Region.validateVertexContainment() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateTransitionScope", func(t *testing.T) {
		tests := []struct {
			name    string
			region  *Region
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid - transitions with vertices in region",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						{
							ID:   "s2",
							Name: "State2",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
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
					},
				},
				wantErr: false,
			},
			{
				name: "valid - internal transition with same source and target",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
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
					},
				},
				wantErr: false,
			},
			{
				name: "valid - empty transitions",
				region: &Region{
					ID:          "r1",
					Name:        "TestRegion",
					Vertices:    []*Vertex{},
					Transitions: []*Transition{},
				},
				wantErr: false,
			},
			{
				name: "invalid - source vertex not in region",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s2",
							Name: "State2",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "s1", // Not in vertices
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
					},
				},
				wantErr: true,
				errMsg:  "transition at index 0 has source vertex (ID: s1) that is not contained in this region (UML constraint)",
			},
			{
				name: "invalid - internal transition target not in region",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Target: &Vertex{
								ID:   "s2", // Not in vertices
								Name: "State2",
								Type: "state",
							},
							Kind: TransitionKindInternal,
						},
					},
				},
				wantErr: true,
				errMsg:  "transition at index 0 has target vertex (ID: s2) that is not contained in this region, but transition kind is internal (UML constraint)",
			},
			{
				name: "invalid - local transition target not in region",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Target: &Vertex{
								ID:   "s2", // Not in vertices
								Name: "State2",
								Type: "state",
							},
							Kind: TransitionKindLocal,
						},
					},
				},
				wantErr: true,
				errMsg:  "transition at index 0 has target vertex (ID: s2) that is not contained in this region, but transition kind is local (UML constraint)",
			},
			{
				name: "invalid - internal transition with different source and target",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						{
							ID:   "s2",
							Name: "State2",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Target: &Vertex{
								ID:   "s2", // Different from source
								Name: "State2",
								Type: "state",
							},
							Kind: TransitionKindInternal,
						},
					},
				},
				wantErr: true,
				errMsg:  "internal transition at index 0 must have the same source and target vertex (UML constraint)",
			},
			{
				name: "invalid - final state as source",
				region: &Region{
					ID:   "r1",
					Name: "TestRegion",
					Vertices: []*Vertex{
						{
							ID:   "fs1",
							Name: "FinalState1",
							Type: "finalstate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "fs1",
								Name: "FinalState1",
								Type: "finalstate",
							},
							Target: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Kind: TransitionKindExternal,
						},
					},
				},
				wantErr: true,
				errMsg:  "transition at index 0 has a final state as source, which is not allowed (UML constraint)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.region.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("Region.validateTransitionScope() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("Region.validateTransitionScope() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("Region.validateTransitionScope() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("comprehensive UML constraint validation", func(t *testing.T) {
		tests := []struct {
			name    string
			region  *Region
			wantErr bool
			errMsgs []string // Multiple error messages to check for
		}{
			{
				name: "valid region with all UML constraints satisfied",
				region: &Region{
					ID:   "r1",
					Name: "ValidRegion",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
						},
					},
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
						{
							ID:   "fs1",
							Name: "FinalState1",
							Type: "finalstate",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "initial1",
								Name: "Initial",
								Type: "pseudostate",
							},
							Target: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Kind: TransitionKindExternal,
						},
						{
							ID: "t2",
							Source: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Target: &Vertex{
								ID:   "fs1",
								Name: "FinalState1",
								Type: "finalstate",
							},
							Kind: TransitionKindExternal,
						},
					},
				},
				wantErr: false,
			},
			{
				name: "multiple UML constraint violations",
				region: &Region{
					ID:   "r2",
					Name: "InvalidRegion",
					States: []*State{
						{
							Vertex: Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
						},
						{
							Vertex: Vertex{
								ID:   "s3", // Not in vertices - containment violation
								Name: "State3",
								Type: "state",
							},
						},
					},
					Vertices: []*Vertex{
						{
							ID:   "initial1",
							Name: "Initial",
							Type: "pseudostate",
						},
						{
							ID:   "initial2", // Second initial - multiplicity violation
							Name: "initial",
							Type: "pseudostate",
						},
						{
							ID:   "s1",
							Name: "State1",
							Type: "state",
						},
						{
							ID:   "fs1",
							Name: "FinalState1",
							Type: "finalstate",
						},
						{
							ID:   "", // Empty ID - containment violation
							Name: "InvalidVertex",
							Type: "state",
						},
					},
					Transitions: []*Transition{
						{
							ID: "t1",
							Source: &Vertex{
								ID:   "fs1", // Final state as source - compatibility violation
								Name: "FinalState1",
								Type: "finalstate",
							},
							Target: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Kind: TransitionKindExternal,
						},
						{
							ID: "t2",
							Source: &Vertex{
								ID:   "s1",
								Name: "State1",
								Type: "state",
							},
							Target: &Vertex{
								ID:   "s2", // Not in region - scope violation
								Name: "State2",
								Type: "state",
							},
							Kind: TransitionKindInternal,
						},
					},
				},
				wantErr: true,
				errMsgs: []string{
					"Region can have at most one initial pseudostate, found 2 at indices: [0 1] (UML constraint)",
					"state at index 1 (ID: s3) is not contained in the region's vertices collection (UML constraint)",
					"vertex at index 4 must have a valid ID for proper containment (UML constraint)",
					"transition at index 0 has a final state as source, which is not allowed (UML constraint)",
					"transition at index 1 has target vertex (ID: s2) that is not contained in this region, but transition kind is internal (UML constraint)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.region.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("Region.Validate() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("Region.Validate() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("Region.Validate() unexpected error = %v", err)
					}
				}
			})
		}
	})
}
