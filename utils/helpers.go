package utils

func MakeSlice[T any](fill T, size ...int) []T {
	var s, c int

	if len(size) > 1 {
		s, c = size[0], size[1]
	} else if len(size) == 1 {
		s = size[0]
	}

	if c < s {
		c = s
	}

	slice := make([]T, s, c)
	for i := range slice {
		slice[i] = fill
	}
	return slice
}