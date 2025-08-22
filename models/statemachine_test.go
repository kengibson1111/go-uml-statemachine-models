package models

import (
	"testing"
	"time"
)

func TestStateMachine_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sm      *StateMachine
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid state machine",
			sm: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
					},
				},
				Entities:  make(map[string]string),
				Metadata:  make(map[string]interface{}),
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			sm: &StateMachine{
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
					},
				},
			},
			wantErr: true,
			errMsg:  "[Required] StateMachine.ID: field is required and cannot be empty",
		},
		{
			name: "empty Name",
			sm: &StateMachine{
				ID:      "sm1",
				Version: "1.0",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
					},
				},
			},
			wantErr: true,
			errMsg:  "[Required] StateMachine.Name: field is required and cannot be empty",
		},
		{
			name: "empty Version",
			sm: &StateMachine{
				ID:   "sm1",
				Name: "TestStateMachine",
				Regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
					},
				},
			},
			wantErr: true,
			errMsg:  "[Required] StateMachine.Version: field is required and cannot be empty",
		},
		{
			name: "invalid region",
			sm: &StateMachine{
				ID:      "sm1",
				Name:    "TestStateMachine",
				Version: "1.0",
				Regions: []*Region{
					{
						// Missing ID and Name
					},
				},
			},
			wantErr: true,
			errMsg:  "[Required] Region.ID: field is required and cannot be empty at Regions[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sm.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("StateMachine.Validate() expected error but got none")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("StateMachine.Validate() error = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("StateMachine.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestRegion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		region  *Region
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid region",
			region: &Region{
				ID:          "r1",
				Name:        "TestRegion",
				States:      []*State{},
				Transitions: []*Transition{},
				Vertices:    []*Vertex{},
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			region: &Region{
				Name: "TestRegion",
			},
			wantErr: true,
			errMsg:  "Region ID cannot be empty",
		},
		{
			name: "empty Name",
			region: &Region{
				ID: "r1",
			},
			wantErr: true,
			errMsg:  "Region Name cannot be empty",
		},
		{
			name: "invalid state",
			region: &Region{
				ID:   "r1",
				Name: "TestRegion",
				States: []*State{
					{
						Vertex: Vertex{
							// Missing required fields
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid state at index 0: invalid vertex in state: Vertex ID cannot be empty",
		},
		{
			name: "invalid transition",
			region: &Region{
				ID:   "r1",
				Name: "TestRegion",
				Transitions: []*Transition{
					{
						// Missing required fields
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid transition at index 0: Transition ID cannot be empty",
		},
		{
			name: "invalid vertex",
			region: &Region{
				ID:   "r1",
				Name: "TestRegion",
				Vertices: []*Vertex{
					{
						// Missing required fields
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid vertex at index 0: Vertex ID cannot be empty",
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
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Region.Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Region.Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestStateMachine_UMLConstraintValidation tests the UML-specific constraint validation methods
func TestStateMachine_UMLConstraintValidation(t *testing.T) {
	t.Run("validateConnectionPoints", func(t *testing.T) {
		tests := []struct {
			name             string
			connectionPoints []*Pseudostate
			wantErr          bool
			errMsg           string
		}{
			{
				name:             "no connection points",
				connectionPoints: []*Pseudostate{},
				wantErr:          false,
			},
			{
				name: "valid entry point connection point",
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "EntryPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindEntryPoint,
					},
				},
				wantErr: false,
			},
			{
				name: "valid exit point connection point",
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp2",
							Name: "ExitPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindExitPoint,
					},
				},
				wantErr: false,
			},
			{
				name: "multiple valid connection points",
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "EntryPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindEntryPoint,
					},
					{
						Vertex: Vertex{
							ID:   "cp2",
							Name: "ExitPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindExitPoint,
					},
				},
				wantErr: false,
			},
			{
				name: "invalid connection point kind - initial",
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "InvalidCP",
							Type: "pseudostate",
						},
						Kind: PseudostateKindInitial,
					},
				},
				wantErr: true,
				errMsg:  "connection point at index 0 must be an entry point or exit point pseudostate, got: initial",
			},
			{
				name: "invalid connection point kind - junction",
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "InvalidCP",
							Type: "pseudostate",
						},
						Kind: PseudostateKindJunction,
					},
				},
				wantErr: true,
				errMsg:  "connection point at index 0 must be an entry point or exit point pseudostate, got: junction",
			},
			{
				name: "invalid connection point type",
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "InvalidCP",
							Type: "state",
						},
						Kind: PseudostateKindEntryPoint,
					},
				},
				wantErr: true,
				errMsg:  "connection point at index 0 must have type 'pseudostate', got: state",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sm := &StateMachine{
					ID:               "sm1",
					Name:             "TestStateMachine",
					Version:          "1.0",
					ConnectionPoints: tt.connectionPoints,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "TestRegion",
						},
					},
				}

				err := sm.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("StateMachine.validateConnectionPoints() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("StateMachine.validateConnectionPoints() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("StateMachine.validateConnectionPoints() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateRegionMultiplicity", func(t *testing.T) {
		tests := []struct {
			name    string
			regions []*Region
			wantErr bool
			errMsg  string
		}{
			{
				name: "valid - one region",
				regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion",
					},
				},
				wantErr: false,
			},
			{
				name: "valid - multiple regions",
				regions: []*Region{
					{
						ID:   "r1",
						Name: "TestRegion1",
					},
					{
						ID:   "r2",
						Name: "TestRegion2",
					},
				},
				wantErr: false,
			},
			{
				name:    "invalid - no regions",
				regions: []*Region{},
				wantErr: true,
				errMsg:  "StateMachine must have at least one region (UML constraint)",
			},
			{
				name:    "invalid - nil regions",
				regions: nil,
				wantErr: true,
				errMsg:  "StateMachine must have at least one region (UML constraint)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sm := &StateMachine{
					ID:      "sm1",
					Name:    "TestStateMachine",
					Version: "1.0",
					Regions: tt.regions,
				}

				err := sm.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("StateMachine.validateRegionMultiplicity() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("StateMachine.validateRegionMultiplicity() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("StateMachine.validateRegionMultiplicity() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("validateMethodConstraints", func(t *testing.T) {
		tests := []struct {
			name             string
			isMethod         bool
			connectionPoints []*Pseudostate
			wantErr          bool
			errMsg           string
		}{
			{
				name:             "valid - not a method, no connection points",
				isMethod:         false,
				connectionPoints: []*Pseudostate{},
				wantErr:          false,
			},
			{
				name:     "valid - not a method, has connection points",
				isMethod: false,
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "EntryPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindEntryPoint,
					},
				},
				wantErr: false,
			},
			{
				name:             "valid - is a method, no connection points",
				isMethod:         true,
				connectionPoints: []*Pseudostate{},
				wantErr:          false,
			},
			{
				name:     "invalid - is a method, has connection points",
				isMethod: true,
				connectionPoints: []*Pseudostate{
					{
						Vertex: Vertex{
							ID:   "cp1",
							Name: "EntryPoint",
							Type: "pseudostate",
						},
						Kind: PseudostateKindEntryPoint,
					},
				},
				wantErr: true,
				errMsg:  "StateMachine used as method cannot have connection points (UML constraint)",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sm := &StateMachine{
					ID:               "sm1",
					Name:             "TestStateMachine",
					Version:          "1.0",
					IsMethod:         tt.isMethod,
					ConnectionPoints: tt.connectionPoints,
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "TestRegion",
						},
					},
				}

				err := sm.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("StateMachine.validateMethodConstraints() expected error but got none")
						return
					}
					if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
						t.Errorf("StateMachine.validateMethodConstraints() error = %v, want to contain %v", err.Error(), tt.errMsg)
					}
				} else {
					if err != nil {
						t.Errorf("StateMachine.validateMethodConstraints() unexpected error = %v", err)
					}
				}
			})
		}
	})

	t.Run("comprehensive UML constraint validation", func(t *testing.T) {
		tests := []struct {
			name    string
			sm      *StateMachine
			wantErr bool
			errMsgs []string // Multiple error messages to check for
		}{
			{
				name: "valid state machine with all UML constraints satisfied",
				sm: &StateMachine{
					ID:      "sm1",
					Name:    "ValidStateMachine",
					Version: "1.0",
					Regions: []*Region{
						{
							ID:   "r1",
							Name: "MainRegion",
						},
					},
					ConnectionPoints: []*Pseudostate{
						{
							Vertex: Vertex{
								ID:   "entry1",
								Name: "EntryPoint1",
								Type: "pseudostate",
							},
							Kind: PseudostateKindEntryPoint,
						},
						{
							Vertex: Vertex{
								ID:   "exit1",
								Name: "ExitPoint1",
								Type: "pseudostate",
							},
							Kind: PseudostateKindExitPoint,
						},
					},
					IsMethod: false,
				},
				wantErr: false,
			},
			{
				name: "multiple UML constraint violations",
				sm: &StateMachine{
					ID:      "sm2",
					Name:    "InvalidStateMachine",
					Version: "1.0",
					Regions: []*Region{}, // Violates region multiplicity
					ConnectionPoints: []*Pseudostate{
						{
							Vertex: Vertex{
								ID:   "cp1",
								Name: "InvalidCP",
								Type: "pseudostate",
							},
							Kind: PseudostateKindInitial, // Invalid connection point kind
						},
					},
					IsMethod: true, // Method with connection points - violation
				},
				wantErr: true,
				errMsgs: []string{
					"StateMachine must have at least one region (UML constraint)",
					"connection point at index 0 must be an entry point or exit point pseudostate, got: initial",
					"StateMachine used as method cannot have connection points (UML constraint)",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.sm.Validate()
				if tt.wantErr {
					if err == nil {
						t.Errorf("StateMachine.Validate() expected error but got none")
						return
					}
					for _, errMsg := range tt.errMsgs {
						if !contains(err.Error(), errMsg) {
							t.Errorf("StateMachine.Validate() error = %v, want to contain %v", err.Error(), errMsg)
						}
					}
				} else {
					if err != nil {
						t.Errorf("StateMachine.Validate() unexpected error = %v", err)
					}
				}
			})
		}
	})
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
