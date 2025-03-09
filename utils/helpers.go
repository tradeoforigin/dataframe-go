package utils

// MakeSlice creates and returns a slice of type `T` with a specified size and capacity.
// The slice is filled with the provided `fill` value. The `size` and `capacity` parameters
// are optional and can be passed just like in the standard `make` function. If only one size 
// is provided, the capacity will default to that size. If both size and capacity are provided, 
// the slice will be created with that specific size and capacity. If the capacity is less than 
// the size, it will be adjusted to match the size.
func MakeSlice[T any](fill T, size ...int) []T {
	var s, c int

	// Set size and capacity based on the provided arguments.
	if len(size) > 1 {
		s, c = size[0], size[1]
	} else if len(size) == 1 {
		s = size[0]
	}

	// Ensure that the capacity is at least the size of the slice.
	if c < s {
		c = s
	}

	// Create the slice with the specified size and capacity.
	slice := make([]T, s, c)

	// Fill the slice with the provided value.
	for i := range slice {
		slice[i] = fill
	}
	return slice
}
