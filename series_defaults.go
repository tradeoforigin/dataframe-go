package dataframe

import (
	"fmt"
)

// DefaultValueFormatter returns a string representation of the data in a particular row.
// If the value is nil, it returns "NaN". Otherwise, it uses the default string formatting
// for the value based on its type.

// This function is commonly used to format the values in a series or data frame when converting
// them to string for display purposes.

// Parameters:
//
//	v (interface{}): The value to be formatted as a string. It can be of any type.
//
// Returns:
//
//	string: A string representation of the value. If the value is nil, it returns "NaN".
//	        Otherwise, it returns the formatted string using `fmt.Sprintf`.
func DefaultValueFormatter(v any) string {
	if v == nil {
		return "NaN"
	}
	return fmt.Sprintf("%v", v)
}
