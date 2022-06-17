package dataframe

import (
	"fmt"
)

// DefaultValueFormatter will return a string representation
// of the data in a particular row.
func DefaultValueFormatter(v interface{}) string {
	if v == nil {
		return "NaN"
	}
	return fmt.Sprintf("%v", v)
}