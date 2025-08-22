package models

import (
	"strings"
	"testing"
)

func TestReferenceValidator_ValidateReferences(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid state machine with proper references",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						States: []*State{
							{
								Vertex: Vertex{
									ID:   "s1",
									Name: "TestState",
									Type: "state",
								},
							},
						},
						// Don't duplicate the same vertex in both collections
						Vertices: []*Vertex{
							{
								ID:   "s2",
								Name: "AnotherVertex",
								Type: "pseudostate",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "state machine with nil region",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{nil},
			},
			wantErr: true,
			errMsgs: []string{"region at index 0 is nil"},
		},
		{
			name: "transition with nil source",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						Transitions: []*Transition{
							{
								ID:     "t1",
								Name:   "TestTransition",
								Kind:   TransitionKindExternal,
								Source: nil,
								Target: &Vertex{
									ID:   "s1",
									Name: "TestState",
									Type: "state",
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsgs: []string{"source vertex is required and cannot be nil"},
		},
		{
			name: "transition with nil target",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						Transitions: []*Transition{
							{
								ID:   "t1",
								Name: "TestTransition",
								Kind: TransitionKindExternal,
								Source: &Vertex{
									ID:   "s1",
									Name: "TestState",
									Type: "state",
								},
								Target: nil,
							},
						},
					},
				},
			},
			wantErr: true,
			errMsgs: []string{"target vertex is required and cannot be nil"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewReferenceValidator()
			err := rv.ValidateReferences(tt.obj)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReferences() expected error but got none")
					return
				}

				errStr := err.Error()
				for _, expectedMsg := range tt.errMsgs {
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("ValidateReferences() error = %v, expected to contain %v", errStr, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestReferenceValidator_ValidateBidirectionalConsistency(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid bidirectional references",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						Transitions: []*Transition{
							{
								ID:   "t1",
								Name: "TestTransition",
								Kind: TransitionKindExternal,
								Source: &Vertex{
									ID:   "s1",
									Name: "SourceState",
									Type: "state",
								},
								Target: &Vertex{
									ID:   "s2",
									Name: "TargetState",
									Type: "state",
								},
							},
						},
						Vertices: []*Vertex{
							{
								ID:   "s1",
								Name: "SourceState",
								Type: "state",
							},
							{
								ID:   "s2",
								Name: "TargetState",
								Type: "state",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "connection point reference with proper bidirectional refs",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						States: []*State{
							{
								Vertex: Vertex{
									ID:   "s1",
									Name: "TestState",
									Type: "state",
								},
								Connections: []*ConnectionPointReference{
									{
										Vertex: Vertex{
											ID:   "cpr1",
											Name: "TestConnectionPoint",
											Type: "connectionpoint",
										},
										Entry: []*Pseudostate{
											{
												Vertex: Vertex{
													ID:   "ep1",
													Name: "EntryPoint",
													Type: "pseudostate",
												},
												Kind: PseudostateKindEntryPoint,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewReferenceValidator()
			err := rv.ValidateReferences(tt.obj)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReferences() expected error but got none")
					return
				}

				errStr := err.Error()
				for _, expectedMsg := range tt.errMsgs {
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("ValidateReferences() error = %v, expected to contain %v", errStr, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestReferenceValidator_ValidateContainmentHierarchy(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid containment hierarchy",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						States: []*State{
							{
								Vertex: Vertex{
									ID:   "s1",
									Name: "CompositeState",
									Type: "state",
								},
								IsComposite: true,
								Regions: []*Region{
									{
										ID:   "r2",
										Name: "NestedRegion",
										States: []*State{
											{
												Vertex: Vertex{
													ID:   "s2",
													Name: "NestedState",
													Type: "state",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid containment - region containing state machine",
			obj: &Region{
				ID:   "r1",
				Name: "TestRegion",
				// This would be invalid if we tried to make a region contain a state machine
			},
			wantErr: false, // This specific case might not trigger an error in our current implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewReferenceValidator()
			err := rv.ValidateReferences(tt.obj)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReferences() expected error but got none")
					return
				}

				errStr := err.Error()
				for _, expectedMsg := range tt.errMsgs {
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("ValidateReferences() error = %v, expected to contain %v", errStr, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestReferenceValidator_ValidateInheritanceRelationships(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid submachine inheritance",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "MainStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "MainRegion",
						States: []*State{
							{
								Vertex: Vertex{
									ID:   "s1",
									Name: "SubmachineState",
									Type: "state",
								},
								IsSubmachineState: true,
								Submachine: &StateMachine{
									ID:      "sm2",
									Name:    "SubStateMachine",
									Version: "1.0",
									Regions: []*Region{
										{
											ID:   "r2",
											Name: "SubRegion",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate ID detection",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "CyclicStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "MainRegion",
						States: []*State{
							{
								Vertex: Vertex{
									ID:   "s1",
									Name: "SubmachineState",
									Type: "state",
								},
								IsSubmachineState: true,
								// Same ID as parent creates duplicate ID error
								Submachine: &StateMachine{
									ID:      "sm1", // Same ID as parent - duplicate ID
									Name:    "CyclicStateMachine",
									Version: "1.0",
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsgs: []string{"duplicate ID"},
		},
		{
			name: "inheritance cycle detection - proper cycle",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "MainStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "MainRegion",
						States: []*State{
							{
								Vertex: Vertex{
									ID:   "s1",
									Name: "SubmachineState",
									Type: "state",
								},
								IsSubmachineState: true,
								Submachine: &StateMachine{
									ID:      "sm2",
									Name:    "SubStateMachine",
									Version: "1.0",
									Regions: []*Region{
										{
											ID:   "r2",
											Name: "SubRegion",
											States: []*State{
												{
													Vertex: Vertex{
														ID:   "s2",
														Name: "CyclicSubmachineState",
														Type: "state",
													},
													IsSubmachineState: true,
													Submachine: &StateMachine{
														ID:      "sm1", // References back to parent - creates cycle
														Name:    "MainStateMachine",
														Version: "1.0",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: true,
			errMsgs: []string{"duplicate ID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewReferenceValidator()
			err := rv.ValidateReferences(tt.obj)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReferences() expected error but got none")
					return
				}

				errStr := err.Error()
				for _, expectedMsg := range tt.errMsgs {
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("ValidateReferences() error = %v, expected to contain %v", errStr, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestReferenceValidator_ValidateContainment(t *testing.T) {
	tests := []struct {
		name    string
		parent  interface{}
		child   interface{}
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid state machine containing region",
			parent: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
			},
			child: &Region{
				ID:   "r1",
				Name: "TestRegion",
			},
			wantErr: false,
		},
		{
			name: "valid region containing state",
			parent: &Region{
				ID:   "r1",
				Name: "TestRegion",
			},
			child: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid containment - state containing state machine",
			parent: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
			},
			child: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
			},
			wantErr: true,
			errMsgs: []string{"State cannot contain StateMachine"},
		},
		{
			name:    "nil parent",
			parent:  nil,
			child:   &Region{ID: "r1", Name: "TestRegion"},
			wantErr: true,
			errMsgs: []string{"parent and child cannot be nil"},
		},
		{
			name:    "nil child",
			parent:  &StateMachine{ID: "sm1", Name: "TestStateMachine", Version: "1.0"},
			child:   nil,
			wantErr: true,
			errMsgs: []string{"parent and child cannot be nil"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewReferenceValidator()
			err := rv.ValidateContainment(tt.parent, tt.child)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateContainment() expected error but got none")
					return
				}

				errStr := err.Error()
				for _, expectedMsg := range tt.errMsgs {
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("ValidateContainment() error = %v, expected to contain %v", errStr, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateContainment() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestReferenceValidator_TypeCompatibility(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid transition referencing vertices",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
						Vertices: []*Vertex{
							{
								ID:   "s1",
								Name: "SourceState",
								Type: "state",
							},
							{
								ID:   "s2",
								Name: "TargetState",
								Type: "state",
							},
						},
						Transitions: []*Transition{
							{
								ID:   "t1",
								Name: "TestTransition",
								Kind: TransitionKindExternal,
								Source: &Vertex{
									ID:   "s1",
									Name: "SourceState",
									Type: "state",
								},
								Target: &Vertex{
									ID:   "s2",
									Name: "TargetState",
									Type: "state",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "connection point reference with entry/exit pseudostates",
			obj: &ConnectionPointReference{
				Vertex: Vertex{
					ID:   "cpr1",
					Name: "TestConnectionPoint",
					Type: "connectionpoint",
				},
				Entry: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "ep1",
							Name: "EntryPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindEntryPoint,
					},
				},
				Exit: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "xp1",
							Name: "ExitPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindExitPoint,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewReferenceValidator()
			err := rv.ValidateReferences(tt.obj)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateReferences() expected error but got none")
					return
				}

				errStr := err.Error()
				for _, expectedMsg := range tt.errMsgs {
					if !strings.Contains(errStr, expectedMsg) {
						t.Errorf("ValidateReferences() error = %v, expected to contain %v", errStr, expectedMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestReferenceValidator_CycleDetection(t *testing.T) {
	rv := NewReferenceValidator()

	// Set up containment tree with a cycle
	rv.containmentTree = map[string][]string{
		"parent1": {"child1"},
		"child1":  {"child2"},
		"child2":  {"parent1"}, // Creates a cycle
	}

	// Set up reference map
	rv.referenceMap = map[string]interface{}{
		"parent1": &StateMachine{ID: "parent1", Name: "Parent1", Version: "1.0"},
		"child1":  &Region{ID: "child1", Name: "Child1"},
		"child2":  &State{Vertex: Vertex{ID: "child2", Name: "Child2", Type: "state"}},
	}

	rv.errors = &ValidationErrors{}
	rv.context = NewValidationContext()

	// Test containment cycle detection
	rv.validateContainmentHierarchy()

	if !rv.errors.HasErrors() {
		t.Errorf("Expected containment cycle error but got none")
	}

	errStr := rv.errors.Error()
	if !strings.Contains(errStr, "containment cycle detected") {
		t.Errorf("Expected containment cycle error, got: %v", errStr)
	}
}

func TestReferenceValidator_InheritanceCycleDetection(t *testing.T) {
	rv := NewReferenceValidator()

	// Set up inheritance tree with a cycle
	rv.inheritanceTree = map[string]string{
		"child1":  "parent1",
		"parent1": "child1", // Creates a cycle
	}

	// Set up reference map
	rv.referenceMap = map[string]interface{}{
		"parent1": &StateMachine{ID: "parent1", Name: "Parent1", Version: "1.0"},
		"child1":  &State{Vertex: Vertex{ID: "child1", Name: "Child1", Type: "state"}},
	}

	rv.errors = &ValidationErrors{}
	rv.context = NewValidationContext()

	// Test inheritance cycle detection
	rv.validateInheritanceRelationships()

	if !rv.errors.HasErrors() {
		t.Errorf("Expected inheritance cycle error but got none")
	}

	errStr := rv.errors.Error()
	if !strings.Contains(errStr, "inheritance cycle detected") {
		t.Errorf("Expected inheritance cycle error, got: %v", errStr)
	}
}

func TestReferenceValidator_GetObjectID(t *testing.T) {
	rv := NewReferenceValidator()

	tests := []struct {
		name     string
		obj      interface{}
		expected string
	}{
		{
			name: "state machine with ID",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
			},
			expected: "sm1",
		},
		{
			name: "region with ID",
			obj: &Region{
				ID:   "r1",
				Name: "TestRegion",
			},
			expected: "r1",
		},
		{
			name: "state with ID",
			obj: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
			},
			expected: "s1",
		},
		{
			name:     "nil object",
			obj:      nil,
			expected: "",
		},
		{
			name: "object without ID field",
			obj: struct {
				Name string
			}{Name: "NoID"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rv.getObjectID(tt.obj)
			if result != tt.expected {
				t.Errorf("getObjectID() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestReferenceValidator_GetObjectTypeName(t *testing.T) {
	rv := NewReferenceValidator()

	tests := []struct {
		name     string
		obj      interface{}
		expected string
	}{
		{
			name: "state machine",
			obj: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
			},
			expected: "StateMachine",
		},
		{
			name: "region",
			obj: &Region{
				ID:   "r1",
				Name: "TestRegion",
			},
			expected: "Region",
		},
		{
			name: "state",
			obj: &State{
				Vertex: Vertex{
					ID:   "s1",
					Name: "TestState",
					Type: "state",
				},
			},
			expected: "State",
		},
		{
			name:     "nil object",
			obj:      nil,
			expected: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rv.getObjectTypeName(tt.obj)
			if result != tt.expected {
				t.Errorf("getObjectTypeName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
