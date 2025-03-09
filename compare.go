package dataframe

import (
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/constraints"
)

// CompareFn is a function type used for comparing two values of the same type.
type CompareFn[T any] func(T, T) bool

// IsEqualFunc provides basic comparison for comparable types. It returns true if the two values are equal.
// Special handling is provided for NaN values to return true when both are NaN.
func IsEqualFunc[T comparable](f1, f2 T) bool {
	// Handle special case for NaN values
	if isNaN(f1) && isNaN(f2) {
		return true
	}
	
	// Compare using equality for comparable types
	return f1 == f2
}

// IsEqualDefaultFunc provides comparison for any type. It uses the `cmp.Equal` function
// for deep comparison of two values, and handles NaN values correctly by returning true when both are NaN.
func IsEqualDefaultFunc[T any](f1, f2 T) bool {
	// Handle special case for NaN values
	if isNaN(f1) && isNaN(f2) {
		return true
	}

	// Compare using cmp.Equal for general comparison
	return cmp.Equal(f1, f2)
}

// IsLessThanFunc provides a "less than" comparison for ordered types (types that support comparison using <).
// It returns true if f1 is less than f2, with special handling for NaN values.
func IsLessThanFunc[T constraints.Ordered](f1, f2 T) bool {
	// NaN is considered less than any value
	if isNaN(f1) {
		return true
	}

	// If f2 is NaN, then f1 is not less than f2
	if isNaN(f2) {
		return false
	}

	// Compare f1 and f2 using the less than operator
	return f1 < f2
}

// IsEqualPtrFunc provides comparison for pointers to comparable types.
// It returns true if the values pointed to by f1 and f2 are equal, taking into account nil pointers.
func IsEqualPtrFunc[T comparable](f1, f2 *T) bool {
	// If both pointers are nil, they are equal
	if f1 == nil {
		return f2 == nil
	}

	// If only one pointer is nil, they are not equal
	if f2 == nil {
		return false
	}

	// Compare the values pointed to by f1 and f2 using IsEqualFunc
	return IsEqualFunc(*f1, *f2)
}

// IsLessThanPtrFunc provides a "less than" comparison for pointers to ordered types.
// It returns true if the value pointed to by f1 is less than the value pointed to by f2, with special handling for nil pointers.
func IsLessThanPtrFunc[T constraints.Ordered](f1, f2 *T) bool {
	// If f1 is nil, it is considered less than any value
	if f1 == nil {
		return true
	}

	// If f2 is nil, f1 is not less than f2
	if f2 == nil {
		return false
	}

	// Compare the values pointed to by f1 and f2 using IsLessThanFunc
	return IsLessThanFunc(*f1, *f2)
}
