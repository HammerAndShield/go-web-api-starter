package validator

import (
	"regexp"
	"slices"
)

var (
	// EmailRX is a regular expression for validating email addresses.
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator holds a map of error messages associated with validation failures.
// The map keys represent the fields where validation failed.
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the error map is empty, indicating no validation errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the Validator's error map if an error for the given key does not already exist.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check evaluates a boolean condition and adds an error message if the condition is false.
// This is a convenience method to streamline common validation checks.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if a given value exists in a list of allowed values.
// This is useful for validating if a value is among a predefined set of acceptable values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches checks if a given string value matches a specific regular expression.
// This function is often used to validate formats such as emails or phone numbers.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all elements in a slice are unique.
// This function can be used to validate that there are no duplicate values in a slice.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
