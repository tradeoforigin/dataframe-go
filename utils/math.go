package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

var (
	// Epsilon is the smallest difference between two representable floating-point numbers.
	Epsilon = math.Nextafter(1, 2) - 1
)

// Number is a constraint that allows either Signed (integer) or Float types.
type Number interface {
	constraints.Signed | constraints.Float
}

// FibonacciParameters represents optional parameters for generating a Fibonacci sequence.
type FibonacciParameters struct {
	Zero, Weighted bool // Zero determines whether the sequence starts with 0 or 1, Weighted scales the sequence values.
}

// Fibonacci generates the first `n` Fibonacci numbers.
// If the `Zero` flag is set, the sequence starts with 0, else it starts with 1.
// If the `Weighted` flag is set, each number in the sequence is divided by the sum of all numbers.
func Fibonacci[T Number](n int, params *FibonacciParameters) []T {

	var zero, weighted = false, false
	if params != nil {
		zero = params.Zero
		weighted = params.Weighted
	}

	var a, b T = 1, 1

	// Adjust starting values based on the Zero flag
	if zero {
		a, b = 0, 1
	} else {
		n -= 1
	}

	result := make([]T, n+1)
	result[0] = a

	// Generate the Fibonacci sequence
	for i := 1; i < n+1; i++ {
		a, b = b, a + b
		result[i] = a
	}

	// If weighted, divide each value by the sum of the sequence
	if weighted {
		fibSum := Sum(result)
		if fibSum > 0 {
			for i := range result {
				result[i] /= fibSum
			}
		}
	}

	return result
}

// Factorial calculates the factorial of a number `n` (n!).
// It returns the result as type `T`.
func Factorial[T Number](n int) T {
	var f = 1
	if n > 1 {
		for i := 2; i <= n; i++ {
			f *= i
		}
	}
	return T(f)
}

// Combination calculates the binomial coefficient "n choose k" (n C k).
// It returns the result as type `T`. The result is rounded to the nearest integer.
func Combination[T Number](n, k int) T {
	var nck float64 = 1

	// If k > n, the combination is 0
	if k > n {
		return 0
	}

	// If k > n - k, swap values for efficiency
	if (n - k < k) {
		k = n - k
	}

	// Calculate the combination using the formula
	var _n, _k = float64(n), float64(k)
	for i := 1.; i <= _k; i++ {
		nck *= (_n + 1 - i) / i
	}

	// Round the result to avoid floating-point precision issues
	nck = math.Round(nck)

	return T(nck)
}

// PascalsTriangleParameters represents optional parameters for generating Pascal's triangle.
type PascalsTriangleParameters struct {
	Weighted, Inverse bool // Weighted scales the values and Inverse inverts them (1 - value).
}

// PascalsTriangle generates the nth row of Pascal's triangle.
// The `Weighted` flag divides each value by the sum of the row.
// The `Inverse` flag inverts each value in the row (1 - value).
func PascalsTriangle[T Number](n int, params *PascalsTriangleParameters) []T {
	var weighted, inverse = false, false
	if params != nil {
		weighted = params.Weighted
		inverse = params.Inverse
	}

	result := make([]T, n+1)
	for i := 0; i <= n; i++ {
		result[i] = Combination[T](n, i)
	}

	// If weighted, divide each value by the sum of the row
	if weighted {
		sum := Sum(result)
		for i := 0; i <= n; i++ {
			result[i] /= sum
			if inverse {
				result[i] = 1 - result[i]
			}
		}
	}

	return result
}

// SineWaveParameters represents optional parameters for generating a sine wave.
type SineWaveParameters struct {
	Weighted bool // Weighted scales the sine values by the sum of the values.
}

// SineWave generates a sine wave with `n` values.
// If the `Weighted` flag is set, the sine values are divided by the sum of all values.
func SineWave(n int, params *SineWaveParameters) []float64 {
	var weighted = false
	if params != nil {
		weighted = params.Weighted
	}

	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = math.Sin((float64(i) + 1) * math.Pi / (float64(n) + 1))
	}

	// If weighted, divide each value by the sum of the wave
	if weighted {
		sum := Sum(result)
		for i := 0; i < n; i++ {
			result[i] /= sum
		}
	}

	return result
}

// SymmetricTriangleParameters represents optional parameters for generating a symmetric triangle.
type SymmetricTriangleParameters struct {
	Weighted bool // Weighted scales the triangle values by the sum of the values.
}

// SymmetricTriangle generates a symmetric triangle with `n` values.
// If `n` is even, the triangle is mirrored symmetrically. If `n` is odd, the peak value is repeated.
// The `Weighted` flag scales the values by the sum of the values.
func SymmetricTriangle[T Number](n int, params *SymmetricTriangleParameters) []T {
	if n < 2 { n = 2 }

	var weighted = false
	if params != nil {
		weighted = params.Weighted
	}

	var triangle []T
	if n > 2 {
		// Handle the case when `n` is even or odd
		if n % 2 == 0 {
			_t := make([]T, n/ 2)
			for i := 0; i < n / 2; i++ {
				_t[i] = T(i + 1)
			}
			triangle = append(triangle, _t...)
			Reverse(_t)
			triangle = append(triangle, _t...)
		} else {
			var _t = make([]T, int(0.5 * float64(n + 1)))
			for i := 0; i < int(0.5 * float64(n + 1)); i++ {
				_t[i] = T(i + 1)
			}
			triangle = append(triangle, _t...)
			_t = _t[:int(0.5 * float64(n + 1)) - 1]
			Reverse(_t)
			triangle = append(triangle, _t...)
		}
	} else {
		triangle = []T{ 1, 1 }
	}

	// If weighted, divide each value by the sum of the triangle
	if weighted {
		sum := Sum(triangle)
		for i := 0; i < len(triangle); i++ {
			triangle[i] /= sum
		}
	}

	return triangle
}

// Reverse reverses the order of the elements in the slice.
func Reverse[T any](values []T) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

// Sum calculates the sum of all values in a slice and returns the result as type `T`.
func Sum[T Number](values []T) T {
	var m T = 0
	for i := range values {
		m += values[i]
	}
	return m
}

// Arange generates a slice of numbers starting at `from` and ending at `to` with a specified `step`.
// It returns a slice of values of type `T`.
func Arange[T Number](from T, to T, step T) []T {
	result := make([]T, int(math.Ceil(float64((to - from) / step))))
	for i := 0; from < to; i++ {
		result[i] = from
		from += step
	}
	return result
}

// Sign returns the sign of the value (-1 for negative, 1 for positive, 0 for zero).
func Sign[T Number](val T) T {
	if val < 0 {
		return -1
	} else if val > 0 {
		return 1
	}
	return 0
}

// Abs returns the absolute value of `val`.
func Abs[T Number](val T) T {
	return Sign(val) * val
}

// Min returns the minimum value from a list of values.
func Min[T Number](values ...T) T {
	if len(values) == 0 {
		panic("No values provided")
	}

	var m T = 0
	for i, v := range values {
		if i == 0 || v < m {
			m = v
		}
	}
	return m
}

// Max returns the maximum value from a list of values.
func Max[T Number](values ...T) T {
	if len(values) == 0 {
		panic("No values provided")
	}
	var m T = 0
	for i, v := range values {
		if i == 0 || v > m {
			m = v
		}
	}
	return m
}

// ArgMin returns the index of the minimum value from a list of values.
func ArgMin[T Number](values ...T) int {
	var a int
	var m T
	for i, v := range values {
		if i == 0 || v < m {
			m = v
			a = i
		}
	}
	return a
}

// ArgMax returns the index of the maximum value from a list of values.
func ArgMax[T Number](values ...T) int {
	var a int
	var m T
	for i, v := range values {
		if i == 0 || v > m {
			m = v
			a = i
		}
	}
	return a
}

// ClipLower returns the value clipped to be no lower than `lower`.
func ClipLower[T Number](value T, lower T) T {
	if value < lower {
		return lower
	}
	return value
}

// ClipUpper returns the value clipped to be no higher than `upper`.
func ClipUpper[T Number](value T, upper T) T {
	if value > upper {
		return upper
	}
	return value
}

// Clip returns the value clipped to be between `lower` and `upper`.
func Clip[T Number](value T, lower T, upper T) T {
	return ClipUpper(ClipLower(value, lower), upper)
}
