package models

import (
	"fmt"
)

// DemoEnhancedValidation demonstrates the enhanced validation error handling capabilities
func DemoEnhancedValidation() {
	fmt.Println("=== Enhanced Validation Error Handling Demo ===")
	fmt.Println()

	// Create a state machine with multiple validation issues
	sm := &StateMachine{
		ID:      "", // Missing required field
		Name:    "", // Missing required field
		Version: "", // Missing required field
		Regions: []*Region{
			{
				ID:   "", // Missing required field
				Name: "", // Missing required field
				States: []*State{
					{
						Vertex: Vertex{
							ID:   "", // Missing required field
							Name: "", // Missing required field
							Type: "", // Missing required field
						},
						Entry: &Behavior{
							ID:            "", // Missing required field
							Specification: "", // Missing required field
						},
					},
				},
				Transitions: []*Transition{
					{
						ID:     "",  // Missing required field
						Source: nil, // Missing required reference
						Target: nil, // Missing required reference
						Kind:   "",  // Invalid kind
					},
				},
			},
			// Duplicate region to show structural integrity validation
			{
				ID:   "region1",
				Name: "Region 1",
			},
			{
				ID:   "region1", // Duplicate ID
				Name: "Region 1 Duplicate",
			},
		},
	}

	fmt.Println("1. Comprehensive Error Collection:")
	fmt.Println("   - Validates entire object hierarchy")
	fmt.Println("   - Collects ALL errors instead of stopping at first error")
	fmt.Println("   - Provides detailed path tracking")
	fmt.Println()

	// Validate and collect all errors
	err := sm.Validate()
	if err != nil {
		if validationErrors, ok := err.(*ValidationErrors); ok {
			fmt.Printf("   Found %d validation errors:\n\n", validationErrors.Count())

			// Show error summary by type
			fmt.Println("2. Error Summary by Type:")
			summary := validationErrors.GetSummary()
			for errorType, count := range summary {
				fmt.Printf("   - %s: %d errors\n", errorType.String(), count)
			}
			fmt.Println()

			// Show path tracking
			fmt.Println("3. Detailed Path Tracking:")
			for i, validationError := range validationErrors.Errors {
				if i >= 5 { // Show first 5 errors for demo
					fmt.Printf("   ... and %d more errors\n", len(validationErrors.Errors)-5)
					break
				}
				pathStr := "root"
				if len(validationError.Path) > 0 {
					pathStr = fmt.Sprintf("root.%s", fmt.Sprintf("%v", validationError.Path))
				}
				fmt.Printf("   %d. [%s] %s.%s at %s\n",
					i+1,
					validationError.Type.String(),
					validationError.Object,
					validationError.Field,
					pathStr)
			}
			fmt.Println()

			// Show filtering capabilities
			fmt.Println("4. Error Filtering and Querying:")
			requiredErrors := validationErrors.GetErrorsByType(ErrorTypeRequired)
			fmt.Printf("   - Required field errors: %d\n", len(requiredErrors))

			constraintErrors := validationErrors.GetErrorsByType(ErrorTypeConstraint)
			fmt.Printf("   - Constraint violations: %d\n", len(constraintErrors))

			regionErrors := validationErrors.GetErrorsByPath("Regions")
			fmt.Printf("   - Errors in Regions: %d\n", len(regionErrors))
			fmt.Println()

			// Show detailed report
			fmt.Println("5. Detailed Validation Report:")
			fmt.Println("   (Showing first 500 characters of full report)")
			report := validationErrors.GetDetailedReport()
			if len(report) > 500 {
				fmt.Printf("   %s...\n", report[:500])
			} else {
				fmt.Printf("   %s\n", report)
			}
		}
	}

	fmt.Println("\n6. Helper Methods for Common Patterns:")

	// Demonstrate helper methods
	context := NewValidationContext()
	errors := &ValidationErrors{}
	helper := NewValidationHelper()

	// Unique ID validation
	regions := []interface{}{
		&Region{ID: "r1", Name: "Region 1"},
		&Region{ID: "r1", Name: "Region 2"}, // Duplicate
		&Region{ID: "r2", Name: "Region 3"},
	}

	helper.ValidateUniqueIDs(regions, "Regions", "StateMachine", context, errors, func(obj interface{}) string {
		if region, ok := obj.(*Region); ok {
			return region.ID
		}
		return ""
	})

	fmt.Printf("   - Unique ID validation found %d duplicate(s)\n", errors.Count())

	// Conditional required validation
	errors.Clear()
	helper.ValidateConditionalRequired("", "ConnectionPoints", "StateMachine", true, "state machine is used as method", context, errors)
	fmt.Printf("   - Conditional required validation: %d error(s)\n", errors.Count())

	// Collection size validation
	errors.Clear()
	emptyRegions := []*Region{}
	helper.ValidateCollectionSize(emptyRegions, "Regions", "StateMachine", 1, 0, context, errors)
	fmt.Printf("   - Collection size validation: %d error(s)\n", errors.Count())

	fmt.Println("\n7. Context Information and Metadata:")

	// Demonstrate context enhancements
	enhancedContext := NewValidationContext().
		WithStateMachine(&StateMachine{ID: "demo-sm", Name: "Demo StateMachine"}).
		WithRegion(&Region{ID: "demo-region", Name: "Demo Region"}).
		WithPath("States").
		WithPathIndex("Transitions", 0).
		WithMetadata("validationPhase", "UMLConstraints").
		WithMetadata("validationLevel", "comprehensive")

	fmt.Printf("   - Full path: %s\n", enhancedContext.GetFullPath())

	contextInfo := enhancedContext.GetContextInfo()
	fmt.Printf("   - Context includes: StateMachine, Region, Path, and %d metadata entries\n", len(contextInfo["metadata"].(map[string]interface{})))

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nKey Enhancements Implemented:")
	fmt.Println("✓ Comprehensive error collection (multiple errors vs. first error)")
	fmt.Println("✓ Detailed path tracking through object hierarchy")
	fmt.Println("✓ Context information for better debugging")
	fmt.Println("✓ Helper methods for common validation patterns")
	fmt.Println("✓ Error filtering, querying, and reporting capabilities")
	fmt.Println("✓ Enhanced context with metadata and full path generation")
	fmt.Println("✓ Structural integrity validation with duplicate detection")
}
