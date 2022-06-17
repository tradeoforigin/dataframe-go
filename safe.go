package dataframe

import (
	"fmt"
	"strings"
)

func isNaN(f any) bool {
	switch v := f.(type) {
		case float32, float64: return v != v
	}
	return false
}

func formatType[T any]() string {
	return strings.Replace(fmt.Sprintf("%T", *new(T)), "<nil>", "any", 1)
}