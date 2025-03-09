package dataframe

import (
	"fmt"
	"strings"
)

// isNaN checks if the given value is "Not-a-Number" (NaN).
// It works by checking if the value is of type float32 or float64,
// and if so, it checks whether the value is NaN using the property
// that NaN is not equal to itself (v != v).
//
// Parameters:
//   f (any): The value to check for NaN. It can be any type, but this function
//            only works for float32 or float64 types.
//
// Returns:
//   bool: Returns true if the value is NaN (for float32 or float64 types), 
//         otherwise returns false.
func isNaN(f any) bool {
	switch v := f.(type) {
		case float32, float64: return v != v
	}
	return false
}

// formatType returns a string representation of the type of the generic type T.
// The function is used to format the type of a value of type T, replacing the 
// "<nil>" type representation with "any" for clarity.
//
// Parameters:
//   None (T is the type passed into the function)
//
// Returns:
//   string: A string representing the type of T, where "<nil>" is replaced with "any".
func formatType[T any]() string {
	return strings.Replace(fmt.Sprintf("%T", *new(T)), "<nil>", "any", 1)
}
