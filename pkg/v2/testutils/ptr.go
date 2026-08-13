// Package testutils provides test helper utilities.
package testutils

// Ptr returns a pointer to the given value.
func Ptr[T any](value T) *T {
	return &value
}
