package models

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}
