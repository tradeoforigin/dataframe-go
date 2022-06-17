package dataframe

import (
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/constraints"
)

type CompareFn[T any] func(T, T) bool

func IsEqualFunc[T comparable] (f1, f2 T) bool {
	if isNaN(f1) && isNaN(f2) {
		return true
	}
	
	return f1 == f2
}

func IsEqualDefaultFunc[T any] (f1, f2 T) bool {
	if isNaN(f1) && isNaN(f2) {
		return true
	}

	return cmp.Equal(f1, f2)
}

func IsLessThanFunc[T constraints.Ordered] (f1, f2 T) bool {
	if isNaN(f1) {
		return true
	}

	if isNaN(f2) {
		return false
	}

	return f1 < f2
}

func IsEqualPtrFunc[T comparable] (f1, f2 *T) bool {
	if f1 == nil {
		return f2 == nil
	}

	if f2 == nil {
		return false
	}

	return IsEqualFunc(*f1, *f2)
}

func IsLessThanPtrFunc[T constraints.Ordered] (f1, f2 *T) bool {
	if f1 == nil {
		return true
	}

	if f2 == nil {
		return false
	}

	return IsLessThanFunc(*f1, *f2)
}