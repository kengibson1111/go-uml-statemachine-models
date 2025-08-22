package models

import (
	"strings"
	"testing"
)

func TestComplexPatternValidator_ValidateOrthogonalRegions(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil state",
			state:   nil,
			wantErr: true,
			errMsg:  "state cannot be nil",
		},
		{
			name: "single region - no validation needed",
			state: &State{
				Vertex:  Vertex{Name: "CompositeState", Type: "state"},
				Regions: []*Region{{Name: "Region1"}},
			},
			wantErr: false,
		},
		{
			name: "no regions - no validation needed",
			state: &State{
				Vertex:  Vertex{Name: "SimpleState", Type: "state"},
				Regions: []*Region{},
			},
			wantErr: false,
		},
		{
			name: "valid orthogonal regions",
			state: &State{
				Vertex: Vertex{Name: "OrthogonalState", Type: "state"},
				Regions: []*Region{
					{
						Name: "Region1",
						Vertices: []*Vertex{
							{Name: "State1", Type: "state"},
						},
					},
					{
						Name: "Region2",
						Vertices: []*Vertex{
							{Name: "State2", Type: "state"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil region in orthogonal regions",
			state: &State{
				Vertex: Vertex{Name: "OrthogonalState", Type: "state"},
				Regions: []*Region{
					{Name: "Region1"},
					nil,
				},
			},
			wantErr: true,
			errMsg:  "orthogonal region cannot be nil",
		},
		{
			name: "overlapping vertices in orthogonal regions",
			state: &State{
				Vertex: Vertex{Name: "OrthogonalState", Type: "state"},
				Regions: []*Region{
					{
						Name: "Region1",
						Vertices: []*Vertex{
							{Name: "SharedState", Type: "state"},
						},
					},
					{
						Name: "Region2",
						Vertices: []*Vertex{
							{Name: "SharedState", Type: "state"},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "vertex 'SharedState' appears in multiple orthogonal regions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &ValidationContext{}
			cpv := NewComplexPatternValidator(context)

			err := cpv.ValidateOrthogonalRegions(tt.state)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateOrthogonalRegions() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateOrthogonalRegions() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateOrthogonalRegions() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestComplexPatternValidator_ValidateConnectionPointReferences(t *testing.T) {
	tests := []struct {
		name         string
		stateMachine *StateMachine
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "nil state machine",
			stateMachine: nil,
			wantErr:      true,
			errMsg:       "state machine cannot be nil",
		},
		{
			name: "no connection points",
			stateMachine: &StateMachine{
				Name:             "SM1",
				ConnectionPoints: []*Pseudostate{},
				Regions:          []*Region{},
			},
			wantErr: false,
		},
		{
			name: "valid entry point",
			stateMachine: &StateMachine{
				Name: "SM1",
				ConnectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{ID: "ep1", Name: "EntryPoint1", Type: "pseudostate"},
						Kind:   PseudostateKindEntryPoint,
					},
				},
				Regions: []*Region{
					{
						Name: "Region1",
						Transitions: []*Transition{
							{
								Name:   "T1",
								Target: &Vertex{ID: "ep1", Name: "EntryPoint1", Type: "pseudostate"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil connection point",
			stateMachine: &StateMachine{
				Name:             "SM1",
				ConnectionPoints: []*Pseudostate{nil},
				Regions:          []*Region{},
			},
			wantErr: true,
			errMsg:  "connection point cannot be nil",
		},
		{
			name: "invalid connection point kind",
			stateMachine: &StateMachine{
				Name: "SM1",
				ConnectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{Name: "InvalidCP", Type: "pseudostate"},
						Kind:   PseudostateKindInitial,
					},
				},
				Regions: []*Region{},
			},
			wantErr: true,
			errMsg:  "connection point must be entry or exit point",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &ValidationContext{StateMachine: tt.stateMachine}
			cpv := NewComplexPatternValidator(context)

			err := cpv.ValidateConnectionPointReferences(tt.stateMachine)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateConnectionPointReferences() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateConnectionPointReferences() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConnectionPointReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestComplexPatternValidator_ValidateStateMachineInheritance(t *testing.T) {
	tests := []struct {
		name         string
		stateMachine *StateMachine
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "nil state machine",
			stateMachine: nil,
			wantErr:      true,
			errMsg:       "state machine cannot be nil",
		},
		{
			name: "simple state machine - no inheritance",
			stateMachine: &StateMachine{
				Name: "SimpleSM",
				Regions: []*Region{
					{Name: "Region1"},
				},
			},
			wantErr: false,
		},
		{
			name: "state machine with unique region names",
			stateMachine: &StateMachine{
				Name: "ExtendedSM",
				Regions: []*Region{
					{
						Name: "Region1",
						Vertices: []*Vertex{
							{Name: "State1", Type: "state"},
							{Name: "State2", Type: "state"},
						},
					},
					{
						Name: "Region2",
						Vertices: []*Vertex{
							{Name: "State3", Type: "state"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate region names in redefinition",
			stateMachine: &StateMachine{
				Name: "ConflictingSM",
				Regions: []*Region{
					{Name: "DuplicateRegion"},
					{Name: "DuplicateRegion"},
				},
			},
			wantErr: true,
			errMsg:  "duplicate region name 'DuplicateRegion' in redefinition",
		},
		{
			name: "duplicate vertex names in redefinition",
			stateMachine: &StateMachine{
				Name: "ConflictingSM",
				Regions: []*Region{
					{
						Name: "Region1",
						Vertices: []*Vertex{
							{Name: "DuplicateState", Type: "state"},
							{Name: "DuplicateState", Type: "state"},
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate vertex name 'DuplicateState' in redefinition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &ValidationContext{}
			cpv := NewComplexPatternValidator(context)

			err := cpv.ValidateStateMachineInheritance(tt.stateMachine)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateStateMachineInheritance() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateStateMachineInheritance() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateStateMachineInheritance() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestComplexPatternValidator_ValidateRegionSeparation(t *testing.T) {
	tests := []struct {
		name         string
		regions      []*Region
		currentIndex int
		wantErr      bool
		errMsg       string
	}{
		{
			name: "nil region at current index",
			regions: []*Region{
				{Name: "Region1"},
				nil,
			},
			currentIndex: 1,
			wantErr:      true,
			errMsg:       "region at index 1 is nil",
		},
		{
			name: "no vertex overlap",
			regions: []*Region{
				{
					Name: "Region1",
					Vertices: []*Vertex{
						{Name: "State1", Type: "state"},
					},
				},
				{
					Name: "Region2",
					Vertices: []*Vertex{
						{Name: "State2", Type: "state"},
					},
				},
			},
			currentIndex: 0,
			wantErr:      false,
		},
		{
			name: "vertex overlap between regions",
			regions: []*Region{
				{
					Name: "Region1",
					Vertices: []*Vertex{
						{Name: "SharedState", Type: "state"},
					},
				},
				{
					Name: "Region2",
					Vertices: []*Vertex{
						{Name: "SharedState", Type: "state"},
					},
				},
			},
			currentIndex: 0,
			wantErr:      true,
			errMsg:       "vertex 'SharedState' appears in multiple orthogonal regions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &ValidationContext{}
			cpv := NewComplexPatternValidator(context)

			err := cpv.validateRegionSeparation(tt.regions, tt.currentIndex)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateRegionSeparation() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateRegionSeparation() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateRegionSeparation() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestComplexPatternValidator_ValidateConnectionPointReferencesHelper(t *testing.T) {
	tests := []struct {
		name         string
		cp           *Pseudostate
		stateMachine *StateMachine
		wantErr      bool
		errMsg       string
	}{
		{
			name: "entry point as transition target",
			cp: &Pseudostate{
				Vertex: Vertex{ID: "ep1", Name: "EntryPoint1", Type: "pseudostate"},
				Kind:   PseudostateKindEntryPoint,
			},
			stateMachine: &StateMachine{
				Name: "SM1",
				Regions: []*Region{
					{
						Name: "Region1",
						Transitions: []*Transition{
							{
								Name:   "T1",
								Target: &Vertex{ID: "ep1", Name: "EntryPoint1", Type: "pseudostate"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "exit point as transition source",
			cp: &Pseudostate{
				Vertex: Vertex{ID: "ex1", Name: "ExitPoint1", Type: "pseudostate"},
				Kind:   PseudostateKindExitPoint,
			},
			stateMachine: &StateMachine{
				Name: "SM1",
				Regions: []*Region{
					{
						Name: "Region1",
						Transitions: []*Transition{
							{
								Name:   "T1",
								Source: &Vertex{ID: "ex1", Name: "ExitPoint1", Type: "pseudostate"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "unreferenced connection point",
			cp: &Pseudostate{
				Vertex: Vertex{ID: "unused", Name: "UnusedCP", Type: "pseudostate"},
				Kind:   PseudostateKindEntryPoint,
			},
			stateMachine: &StateMachine{
				Name:    "SM1",
				Regions: []*Region{{Name: "Region1", Transitions: []*Transition{}}},
			},
			wantErr: true,
			errMsg:  "connection point 'UnusedCP' is not referenced by any transitions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &ValidationContext{}
			cpv := NewComplexPatternValidator(context)

			err := cpv.validateConnectionPointReferences(tt.cp, tt.stateMachine)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateConnectionPointReferences() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateConnectionPointReferences() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateConnectionPointReferences() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestComplexPatternValidator_ValidateExtendedStateMachineCompatibility(t *testing.T) {
	tests := []struct {
		name         string
		stateMachine *StateMachine
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "nil state machine",
			stateMachine: nil,
			wantErr:      true,
			errMsg:       "cannot validate extension on nil state machine",
		},
		{
			name: "no regions",
			stateMachine: &StateMachine{
				Name:    "EmptySM",
				Regions: []*Region{},
			},
			wantErr: true,
			errMsg:  "extended state machine must have at least one region",
		},
		{
			name: "valid extended state machine",
			stateMachine: &StateMachine{
				Name: "ExtendedSM",
				Regions: []*Region{
					{Name: "Region1"},
				},
				ConnectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{Name: "EntryCP", Type: "pseudostate"},
						Kind:   PseudostateKindEntryPoint,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid connection point kind in extended state machine",
			stateMachine: &StateMachine{
				Name: "InvalidExtendedSM",
				Regions: []*Region{
					{Name: "Region1"},
				},
				ConnectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{Name: "InvalidCP", Type: "pseudostate"},
						Kind:   PseudostateKindInitial,
					},
				},
			},
			wantErr: true,
			errMsg:  "extended state machine connection point 'InvalidCP' must be entry or exit point",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &ValidationContext{}
			cpv := NewComplexPatternValidator(context)

			err := cpv.validateExtendedStateMachineCompatibility(tt.stateMachine)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateExtendedStateMachineCompatibility() expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateExtendedStateMachineCompatibility() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateExtendedStateMachineCompatibility() unexpected error = %v", err)
				}
			}
		})
	}
}

// Integration test for complex patterns
func TestComplexPatternValidator_Integration(t *testing.T) {
	// Create a complex state machine with orthogonal regions and connection points
	stateMachine := &StateMachine{
		Name: "ComplexSM",
		ConnectionPoints: []*Pseudostate{
			{
				Vertex: Vertex{ID: "entry1", Name: "Entry1", Type: "pseudostate"},
				Kind:   PseudostateKindEntryPoint,
			},
			{
				Vertex: Vertex{ID: "exit1", Name: "Exit1", Type: "pseudostate"},
				Kind:   PseudostateKindExitPoint,
			},
		},
		Regions: []*Region{
			{
				Name: "MainRegion",
				States: []*State{
					{
						Vertex: Vertex{Name: "CompositeState", Type: "state"},
						Regions: []*Region{
							{
								Name: "SubRegion1",
								Vertices: []*Vertex{
									{Name: "SubState1", Type: "state"},
								},
							},
							{
								Name: "SubRegion2",
								Vertices: []*Vertex{
									{Name: "SubState2", Type: "state"},
								},
							},
						},
					},
				},
				Transitions: []*Transition{
					{
						Name:   "EntryTransition",
						Target: &Vertex{ID: "entry1", Name: "Entry1", Type: "pseudostate"},
					},
					{
						Name:   "ExitTransition",
						Source: &Vertex{ID: "exit1", Name: "Exit1", Type: "pseudostate"},
					},
				},
			},
		},
	}

	context := &ValidationContext{StateMachine: stateMachine}
	cpv := NewComplexPatternValidator(context)

	// Test connection point validation
	err := cpv.ValidateConnectionPointReferences(stateMachine)
	if err != nil {
		t.Errorf("Integration test failed on connection points: %v", err)
	}

	// Test orthogonal regions validation
	compositeState := stateMachine.Regions[0].States[0]
	err = cpv.ValidateOrthogonalRegions(compositeState)
	if err != nil {
		t.Errorf("Integration test failed on orthogonal regions: %v", err)
	}

	// Test inheritance validation
	err = cpv.ValidateStateMachineInheritance(stateMachine)
	if err != nil {
		t.Errorf("Integration test failed on inheritance: %v", err)
	}
}
