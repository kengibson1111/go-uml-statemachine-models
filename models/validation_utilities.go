package models

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

// ValidationUtilities provides a comprehensive set of utility functions for validation operations
type ValidationUtilities struct {
	helper *ValidationHelper
}

// NewValidationUtilities creates a new validation utilities instance
func NewValidationUtilities() *ValidationUtilities {
	return &ValidationUtilities{
		helper: NewValidationHelper(),
	}
}

// StateMachineTraverser provides utilities for traversing state machine hierarchies
type StateMachineTraverser struct {
	visited map[string]bool
}

// NewStateMachineTraverser creates a new state machine traverser
func NewStateMachineTraverser() *StateMachineTraverser {
	return &StateMachineTraverser{
		visited: make(map[string]bool),
	}
}

// TraversalCallback defines the callback function for traversal operations
type TraversalCallback func(obj interface{}, path []string, depth int) error

// TraverseStateMachine performs a depth-first traversal of a state machine hierarchy
func (smt *StateMachineTraverser) TraverseStateMachine(sm *StateMachine, callback TraversalCallback) error {
	if sm == nil {
		return fmt.Errorf("state machine cannot be nil")
	}

	// Reset visited map for new traversal
	smt.visited = make(map[string]bool)

	return smt.traverseObject(sm, []string{"StateMachine"}, 0, callback)
}

// TraverseRegion performs a depth-first traversal of a region hierarchy
func (smt *StateMachineTraverser) TraverseRegion(region *Region, callback TraversalCallback) error {
	if region == nil {
		return fmt.Errorf("region cannot be nil")
	}

	// Reset visited map for new traversal
	smt.visited = make(map[string]bool)

	return smt.traverseObject(region, []string{"Region"}, 0, callback)
}

// traverseObject recursively traverses an object hierarchy
func (smt *StateMachineTraverser) traverseObject(obj interface{}, path []string, depth int, callback TraversalCallback) error {
	if obj == nil {
		return nil
	}

	// Get object ID to prevent infinite loops
	objID := smt.getObjectID(obj)
	if objID != "" {
		if smt.visited[objID] {
			return nil // Already visited this object
		}
		smt.visited[objID] = true
	}

	// Call the callback for this object
	if err := callback(obj, path, depth); err != nil {
		return err
	}

	// Traverse child objects based on type
	switch v := obj.(type) {
	case *StateMachine:
		return smt.traverseStateMachine(v, path, depth, callback)
	case *Region:
		return smt.traverseRegion(v, path, depth, callback)
	case *State:
		return smt.traverseState(v, path, depth, callback)
	case *Transition:
		return smt.traverseTransition(v, path, depth, callback)
	case *Pseudostate:
		return smt.traversePseudostate(v, path, depth, callback)
	case *FinalState:
		return smt.traverseFinalState(v, path, depth, callback)
	case *ConnectionPointReference:
		return smt.traverseConnectionPointReference(v, path, depth, callback)
	}

	return nil
}

// traverseStateMachine traverses StateMachine children
func (smt *StateMachineTraverser) traverseStateMachine(sm *StateMachine, path []string, depth int, callback TraversalCallback) error {
	// Traverse regions
	for i, region := range sm.Regions {
		if region != nil {
			childPath := append(path, fmt.Sprintf("Regions[%d]", i))
			if err := smt.traverseObject(region, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	// Traverse connection points
	for i, cp := range sm.ConnectionPoints {
		if cp != nil {
			childPath := append(path, fmt.Sprintf("ConnectionPoints[%d]", i))
			if err := smt.traverseObject(cp, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	return nil
}

// traverseRegion traverses Region children
func (smt *StateMachineTraverser) traverseRegion(region *Region, path []string, depth int, callback TraversalCallback) error {
	// Traverse states
	for i, state := range region.States {
		if state != nil {
			childPath := append(path, fmt.Sprintf("States[%d]", i))
			if err := smt.traverseObject(state, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	// Traverse vertices
	for i, vertex := range region.Vertices {
		if vertex != nil {
			childPath := append(path, fmt.Sprintf("Vertices[%d]", i))
			if err := smt.traverseObject(vertex, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	// Traverse transitions
	for i, transition := range region.Transitions {
		if transition != nil {
			childPath := append(path, fmt.Sprintf("Transitions[%d]", i))
			if err := smt.traverseObject(transition, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	return nil
}

// traverseState traverses State children
func (smt *StateMachineTraverser) traverseState(state *State, path []string, depth int, callback TraversalCallback) error {
	// Traverse regions in composite states
	for i, region := range state.Regions {
		if region != nil {
			childPath := append(path, fmt.Sprintf("Regions[%d]", i))
			if err := smt.traverseObject(region, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	// Traverse submachine if present
	if state.Submachine != nil {
		childPath := append(path, "Submachine")
		if err := smt.traverseObject(state.Submachine, childPath, depth+1, callback); err != nil {
			return err
		}
	}

	// Traverse connection point references
	for i, conn := range state.Connections {
		if conn != nil {
			childPath := append(path, fmt.Sprintf("Connections[%d]", i))
			if err := smt.traverseObject(conn, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	return nil
}

// traverseTransition traverses Transition children
func (smt *StateMachineTraverser) traverseTransition(transition *Transition, path []string, depth int, callback TraversalCallback) error {
	// Note: We don't traverse source/target vertices to avoid cycles
	// Those are handled by reference validation

	// Traverse triggers
	for i, trigger := range transition.Triggers {
		if trigger != nil {
			childPath := append(path, fmt.Sprintf("Triggers[%d]", i))
			if err := smt.traverseObject(trigger, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	return nil
}

// traversePseudostate traverses Pseudostate children (usually none)
func (smt *StateMachineTraverser) traversePseudostate(ps *Pseudostate, path []string, depth int, callback TraversalCallback) error {
	// Pseudostates typically don't have child objects
	return nil
}

// traverseFinalState traverses FinalState children (usually none)
func (smt *StateMachineTraverser) traverseFinalState(fs *FinalState, path []string, depth int, callback TraversalCallback) error {
	// Final states typically don't have child objects
	return nil
}

// traverseConnectionPointReference traverses ConnectionPointReference children
func (smt *StateMachineTraverser) traverseConnectionPointReference(cpr *ConnectionPointReference, path []string, depth int, callback TraversalCallback) error {
	// Traverse entry pseudostates
	for i, entry := range cpr.Entry {
		if entry != nil {
			childPath := append(path, fmt.Sprintf("Entry[%d]", i))
			if err := smt.traverseObject(entry, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	// Traverse exit pseudostates
	for i, exit := range cpr.Exit {
		if exit != nil {
			childPath := append(path, fmt.Sprintf("Exit[%d]", i))
			if err := smt.traverseObject(exit, childPath, depth+1, callback); err != nil {
				return err
			}
		}
	}

	return nil
}

// getObjectID extracts the ID from an object using reflection
func (smt *StateMachineTraverser) getObjectID(obj interface{}) string {
	if obj == nil {
		return ""
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	idField := v.FieldByName("ID")
	if !idField.IsValid() || idField.Kind() != reflect.String {
		return ""
	}

	return idField.String()
}

// ValidationResultAggregator provides utilities for aggregating and reporting validation results
type ValidationResultAggregator struct {
	results map[string]*ValidationErrors
}

// NewValidationResultAggregator creates a new validation result aggregator
func NewValidationResultAggregator() *ValidationResultAggregator {
	return &ValidationResultAggregator{
		results: make(map[string]*ValidationErrors),
	}
}

// AddResult adds a validation result for a specific object
func (vra *ValidationResultAggregator) AddResult(objectID string, errors *ValidationErrors) {
	if errors != nil && errors.HasErrors() {
		vra.results[objectID] = errors
	}
}

// AddSingleError adds a single validation error for a specific object
func (vra *ValidationResultAggregator) AddSingleError(objectID string, errorType ValidationErrorType, object, field, message string, path []string) {
	if _, exists := vra.results[objectID]; !exists {
		vra.results[objectID] = &ValidationErrors{}
	}
	vra.results[objectID].AddError(errorType, object, field, message, path)
}

// GetResults returns all validation results
func (vra *ValidationResultAggregator) GetResults() map[string]*ValidationErrors {
	return vra.results
}

// GetTotalErrorCount returns the total number of validation errors across all objects
func (vra *ValidationResultAggregator) GetTotalErrorCount() int {
	total := 0
	for _, errors := range vra.results {
		total += len(errors.Errors)
	}
	return total
}

// HasErrors returns true if there are any validation errors
func (vra *ValidationResultAggregator) HasErrors() bool {
	return len(vra.results) > 0
}

// GetSummaryReport returns a summary report of all validation results
func (vra *ValidationResultAggregator) GetSummaryReport() string {
	if !vra.HasErrors() {
		return "No validation errors found."
	}

	var report strings.Builder
	totalErrors := vra.GetTotalErrorCount()

	report.WriteString(fmt.Sprintf("Validation Summary: %d error(s) found across %d object(s)\n", totalErrors, len(vra.results)))
	report.WriteString(strings.Repeat("=", 60) + "\n\n")

	// Sort object IDs for consistent output
	objectIDs := make([]string, 0, len(vra.results))
	for objectID := range vra.results {
		objectIDs = append(objectIDs, objectID)
	}
	sort.Strings(objectIDs)

	for _, objectID := range objectIDs {
		errors := vra.results[objectID]
		report.WriteString(fmt.Sprintf("Object: %s (%d error(s))\n", objectID, len(errors.Errors)))
		report.WriteString(strings.Repeat("-", 40) + "\n")

		for i, err := range errors.Errors {
			report.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
		}
		report.WriteString("\n")
	}

	return report.String()
}

// GetDetailedReport returns a detailed report of all validation results
func (vra *ValidationResultAggregator) GetDetailedReport() string {
	if !vra.HasErrors() {
		return "No validation errors found."
	}

	var report strings.Builder
	totalErrors := vra.GetTotalErrorCount()

	report.WriteString(fmt.Sprintf("Detailed Validation Report\n"))
	report.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("Total Errors: %d across %d object(s)\n", totalErrors, len(vra.results)))
	report.WriteString(strings.Repeat("=", 80) + "\n\n")

	// Group errors by type across all objects
	errorsByType := make(map[ValidationErrorType][]*ValidationError)
	for _, errors := range vra.results {
		for _, err := range errors.Errors {
			errorsByType[err.Type] = append(errorsByType[err.Type], err)
		}
	}

	// Report by error type
	for errorType, errors := range errorsByType {
		report.WriteString(fmt.Sprintf("%s Errors (%d)\n", errorType.String(), len(errors)))
		report.WriteString(strings.Repeat("-", 50) + "\n")

		for i, err := range errors {
			report.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
			if len(err.Context) > 0 {
				report.WriteString("     Context: ")
				for k, v := range err.Context {
					report.WriteString(fmt.Sprintf("%s=%v ", k, v))
				}
				report.WriteString("\n")
			}
		}
		report.WriteString("\n")
	}

	// Report by object
	report.WriteString("Errors by Object\n")
	report.WriteString(strings.Repeat("-", 50) + "\n")

	objectIDs := make([]string, 0, len(vra.results))
	for objectID := range vra.results {
		objectIDs = append(objectIDs, objectID)
	}
	sort.Strings(objectIDs)

	for _, objectID := range objectIDs {
		errors := vra.results[objectID]
		report.WriteString(fmt.Sprintf("\n%s (%d error(s))\n", objectID, len(errors.Errors)))

		for i, err := range errors.Errors {
			report.WriteString(fmt.Sprintf("  %d. [%s] %s.%s: %s\n", i+1, err.Type.String(), err.Object, err.Field, err.Message))
			if len(err.Path) > 0 {
				report.WriteString(fmt.Sprintf("     Path: %s\n", strings.Join(err.Path, ".")))
			}
		}
	}

	return report.String()
}

// Merge merges another aggregator's results into this one
func (vra *ValidationResultAggregator) Merge(other *ValidationResultAggregator) {
	if other == nil {
		return
	}

	for objectID, errors := range other.results {
		if existingErrors, exists := vra.results[objectID]; exists {
			existingErrors.Merge(errors)
		} else {
			vra.results[objectID] = errors
		}
	}
}

// Clear removes all validation results
func (vra *ValidationResultAggregator) Clear() {
	vra.results = make(map[string]*ValidationErrors)
}

// ValidationDebugger provides debugging utilities for validation troubleshooting
type ValidationDebugger struct {
	traverser  *StateMachineTraverser
	aggregator *ValidationResultAggregator
}

// NewValidationDebugger creates a new validation debugger
func NewValidationDebugger() *ValidationDebugger {
	return &ValidationDebugger{
		traverser:  NewStateMachineTraverser(),
		aggregator: NewValidationResultAggregator(),
	}
}

// DebugStateMachine performs comprehensive debugging of a state machine
func (vd *ValidationDebugger) DebugStateMachine(sm *StateMachine) (*ValidationDebugReport, error) {
	if sm == nil {
		return nil, fmt.Errorf("state machine cannot be nil")
	}

	report := &ValidationDebugReport{
		StateMachineID: sm.ID,
		Timestamp:      time.Now(),
		Objects:        make(map[string]*ObjectDebugInfo),
	}

	// Traverse the state machine and collect debug information
	err := vd.traverser.TraverseStateMachine(sm, func(obj interface{}, path []string, depth int) error {
		debugInfo := vd.analyzeObject(obj, path, depth)
		if debugInfo != nil {
			report.Objects[debugInfo.ID] = debugInfo
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error during traversal: %w", err)
	}

	// Perform validation and collect errors
	vd.aggregator.Clear()
	validationErr := sm.Validate()
	if validationErr != nil {
		if validationErrors, ok := validationErr.(*ValidationErrors); ok {
			vd.aggregator.AddResult(sm.ID, validationErrors)
		}
	}

	report.ValidationResults = vd.aggregator.GetResults()
	report.TotalObjects = len(report.Objects)
	report.TotalErrors = vd.aggregator.GetTotalErrorCount()

	return report, nil
}

// analyzeObject analyzes a single object and returns debug information
func (vd *ValidationDebugger) analyzeObject(obj interface{}, path []string, depth int) *ObjectDebugInfo {
	if obj == nil {
		return nil
	}

	info := &ObjectDebugInfo{
		Path:       strings.Join(path, "."),
		Depth:      depth,
		Properties: make(map[string]interface{}),
	}

	// Use reflection to extract object information
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			info.Type = "nil"
			return info
		}
		v = v.Elem()
	}

	t := v.Type()
	info.Type = t.Name()

	// Extract basic properties
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		fieldName := fieldType.Name

		// Extract field value based on type
		switch field.Kind() {
		case reflect.String:
			info.Properties[fieldName] = field.String()
			if fieldName == "ID" {
				info.ID = field.String()
			}
		case reflect.Bool:
			info.Properties[fieldName] = field.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			info.Properties[fieldName] = field.Int()
		case reflect.Slice:
			info.Properties[fieldName+"_count"] = field.Len()
		case reflect.Ptr:
			if field.IsNil() {
				info.Properties[fieldName] = nil
			} else {
				info.Properties[fieldName+"_present"] = true
			}
		case reflect.Map:
			info.Properties[fieldName+"_count"] = field.Len()
		}
	}

	// If no ID was found, generate one based on type and path
	if info.ID == "" {
		info.ID = fmt.Sprintf("%s_%s", info.Type, strings.ReplaceAll(info.Path, ".", "_"))
	}

	return info
}

// ValidationDebugReport contains comprehensive debugging information
type ValidationDebugReport struct {
	StateMachineID    string                       `json:"state_machine_id"`
	Timestamp         time.Time                    `json:"timestamp"`
	TotalObjects      int                          `json:"total_objects"`
	TotalErrors       int                          `json:"total_errors"`
	Objects           map[string]*ObjectDebugInfo  `json:"objects"`
	ValidationResults map[string]*ValidationErrors `json:"validation_results"`
}

// ObjectDebugInfo contains debugging information for a single object
type ObjectDebugInfo struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Path       string                 `json:"path"`
	Depth      int                    `json:"depth"`
	Properties map[string]interface{} `json:"properties"`
}

// GetSummary returns a summary of the debug report
func (vdr *ValidationDebugReport) GetSummary() string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Debug Report for StateMachine: %s\n", vdr.StateMachineID))
	summary.WriteString(fmt.Sprintf("Generated: %s\n", vdr.Timestamp.Format("2006-01-02 15:04:05")))
	summary.WriteString(fmt.Sprintf("Total Objects: %d\n", vdr.TotalObjects))
	summary.WriteString(fmt.Sprintf("Total Errors: %d\n", vdr.TotalErrors))
	summary.WriteString(strings.Repeat("-", 50) + "\n")

	// Object type distribution
	typeCount := make(map[string]int)
	for _, obj := range vdr.Objects {
		typeCount[obj.Type]++
	}

	summary.WriteString("Object Distribution:\n")
	for objType, count := range typeCount {
		summary.WriteString(fmt.Sprintf("  %s: %d\n", objType, count))
	}

	if vdr.TotalErrors > 0 {
		summary.WriteString("\nValidation Issues Found:\n")
		for objectID, errors := range vdr.ValidationResults {
			summary.WriteString(fmt.Sprintf("  %s: %d error(s)\n", objectID, len(errors.Errors)))
		}
	}

	return summary.String()
}

// CommonValidationPatterns provides utilities for common validation patterns
type CommonValidationPatterns struct {
	helper *ValidationHelper
}

// NewCommonValidationPatterns creates a new common validation patterns utility
func NewCommonValidationPatterns() *CommonValidationPatterns {
	return &CommonValidationPatterns{
		helper: NewValidationHelper(),
	}
}

// ValidateStateMachineStructure validates the overall structure of a state machine
func (cvp *CommonValidationPatterns) ValidateStateMachineStructure(sm *StateMachine, context *ValidationContext, errors *ValidationErrors) {
	if sm == nil {
		errors.AddError(ErrorTypeRequired, "StateMachine", "Structure", "state machine cannot be nil", context.Path)
		return
	}

	// Validate basic structure requirements
	cvp.helper.ValidateRequired(sm.ID, "ID", "StateMachine", context, errors)
	cvp.helper.ValidateRequired(sm.Name, "Name", "StateMachine", context, errors)
	cvp.helper.ValidateRequired(sm.Version, "Version", "StateMachine", context, errors)

	// Validate region multiplicity (UML constraint)
	cvp.helper.ValidateCollectionSize(sm.Regions, "Regions", "StateMachine", 1, 0, context, errors)

	// Validate connection points are appropriate types
	for i, cp := range sm.ConnectionPoints {
		if cp != nil {
			cpContext := context.WithPathIndex("ConnectionPoints", i)
			if cp.Kind != PseudostateKindEntryPoint && cp.Kind != PseudostateKindExitPoint {
				errors.AddError(
					ErrorTypeConstraint,
					"StateMachine",
					"ConnectionPoints",
					fmt.Sprintf("connection point must be entry or exit point, got: %s", cp.Kind),
					cpContext.Path,
				)
			}
		}
	}

	// Validate method constraints
	if sm.IsMethod && len(sm.ConnectionPoints) > 0 {
		errors.AddError(
			ErrorTypeConstraint,
			"StateMachine",
			"IsMethod",
			"state machine used as method cannot have connection points",
			context.Path,
		)
	}
}

// ValidateRegionStructure validates the structure of a region
func (cvp *CommonValidationPatterns) ValidateRegionStructure(region *Region, context *ValidationContext, errors *ValidationErrors) {
	if region == nil {
		errors.AddError(ErrorTypeRequired, "Region", "Structure", "region cannot be nil", context.Path)
		return
	}

	// Validate basic structure requirements
	cvp.helper.ValidateRequired(region.ID, "ID", "Region", context, errors)
	cvp.helper.ValidateRequired(region.Name, "Name", "Region", context, errors)

	// Validate initial state multiplicity (at most one initial pseudostate)
	initialCount := 0
	for i, vertex := range region.Vertices {
		if vertex != nil && vertex.Type == "pseudostate" {
			// Check if this is an initial pseudostate using naming conventions
			if cvp.isInitialPseudostate(vertex) {
				initialCount++
				if initialCount > 1 {
					errors.AddError(
						ErrorTypeMultiplicity,
						"Region",
						"Vertices",
						fmt.Sprintf("region can have at most one initial pseudostate, found multiple at index %d", i),
						context.WithPathIndex("Vertices", i).Path,
					)
				}
			}
		}
	}
}

// ValidateTransitionStructure validates the structure of a transition
func (cvp *CommonValidationPatterns) ValidateTransitionStructure(transition *Transition, context *ValidationContext, errors *ValidationErrors) {
	if transition == nil {
		errors.AddError(ErrorTypeRequired, "Transition", "Structure", "transition cannot be nil", context.Path)
		return
	}

	// Validate basic structure requirements
	cvp.helper.ValidateRequired(transition.ID, "ID", "Transition", context, errors)
	cvp.helper.ValidateRequiredPointer(transition.Source, "Source", "Transition", context, errors)
	cvp.helper.ValidateRequiredPointer(transition.Target, "Target", "Transition", context, errors)

	// Validate transition kind constraints
	if transition.Source != nil && transition.Target != nil {
		switch transition.Kind {
		case TransitionKindInternal:
			if transition.Source.ID != transition.Target.ID {
				errors.AddError(
					ErrorTypeConstraint,
					"Transition",
					"Kind",
					"internal transition must have same source and target",
					context.Path,
				)
			}
		case TransitionKindLocal:
			// Local transitions have specific containment rules
			// Implementation depends on containment validation
		case TransitionKindExternal:
			// External transitions can cross region boundaries
			// No additional constraints here
		}
	}
}

// isInitialPseudostate checks if a vertex is an initial pseudostate using naming conventions
func (cvp *CommonValidationPatterns) isInitialPseudostate(vertex *Vertex) bool {
	if vertex == nil || vertex.Type != "pseudostate" {
		return false
	}

	// Check common naming patterns for initial pseudostates
	name := strings.ToLower(vertex.Name)
	id := strings.ToLower(vertex.ID)

	initialPatterns := []string{"initial", "init", "start"}

	for _, pattern := range initialPatterns {
		if name == pattern || id == pattern || strings.Contains(name, pattern) || strings.Contains(id, pattern) {
			return true
		}
	}

	return false
}

// ValidateObjectHierarchy validates the hierarchical structure of objects
func (cvp *CommonValidationPatterns) ValidateObjectHierarchy(obj interface{}, context *ValidationContext, errors *ValidationErrors) {
	if obj == nil {
		return
	}

	// Use traverser to validate hierarchy
	traverser := NewStateMachineTraverser()

	err := traverser.traverseObject(obj, context.Path, 0, func(currentObj interface{}, path []string, depth int) error {
		// Validate depth limits to prevent excessive nesting
		if depth > 10 {
			errors.AddError(
				ErrorTypeConstraint,
				"Hierarchy",
				"Depth",
				fmt.Sprintf("object hierarchy depth exceeds maximum allowed (10), current depth: %d", depth),
				path,
			)
		}

		// Validate object consistency at each level
		cvp.validateObjectConsistency(currentObj, path, errors)

		return nil
	})

	if err != nil {
		errors.AddError(
			ErrorTypeReference,
			"Hierarchy",
			"Traversal",
			fmt.Sprintf("error during hierarchy validation: %s", err.Error()),
			context.Path,
		)
	}
}

// validateObjectConsistency validates consistency of a single object
func (cvp *CommonValidationPatterns) validateObjectConsistency(obj interface{}, path []string, errors *ValidationErrors) {
	if obj == nil {
		return
	}

	// Basic consistency checks using reflection
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	objType := t.Name()

	// Check for required ID field
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.Kind() == reflect.String {
		if idField.String() == "" {
			errors.AddError(
				ErrorTypeRequired,
				objType,
				"ID",
				"object ID cannot be empty",
				path,
			)
		}
	}

	// Check for required Name field (if present)
	nameField := v.FieldByName("Name")
	if nameField.IsValid() && nameField.Kind() == reflect.String {
		if nameField.String() == "" {
			errors.AddError(
				ErrorTypeRequired,
				objType,
				"Name",
				"object name cannot be empty",
				path,
			)
		}
	}
}
