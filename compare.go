package dataframe

import (
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/constraints"
)

// CompareFn type for compare function for comparision values of same type
type CompareFn[T any] func(T, T) bool

// IsEqualFunc provides basic comparison for comparable types
func IsEqualFunc[T comparable] (f1, f2 T) bool {
	if isNaN(f1) && isNaN(f2) {
		return true
	}
	
	return f1 == f2
}

// IsEqualDefaultFunc provides comparaision for any type
func IsEqualDefaultFunc[T any] (f1, f2 T) bool {
	if isNaN(f1) && isNaN(f2) {
		return true
	}

	return cmp.Equal(f1, f2)
}

// IsLessThanFunc provides (less than) comparision for Ordered types
func IsLessThanFunc[T constraints.Ordered] (f1, f2 T) bool {
	if isNaN(f1) {
		return true
	}

	if isNaN(f2) {
		return false
	}

	return f1 < f2
}

// IsEqualPtrFunc provides comparision for pointers of comparable types
func IsEqualPtrFunc[T comparable] (f1, f2 *T) bool {
	if f1 == nil {
		return f2 == nil
	}

	if f2 == nil {
		return false
	}

	return IsEqualFunc(*f1, *f2)
}

// IsLessThanPtrFunc provides (less than) comparision for pointers of Ordered types
func IsLessThanPtrFunc[T constraints.Ordered] (f1, f2 *T) bool {
	if f1 == nil {
		return true
	}

	if f2 == nil {
		return false
	}

	return IsLessThanFunc(*f1, *f2)
}