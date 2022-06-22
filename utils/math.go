package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

var (
	Epsilon = math.Nextafter(1, 2) - 1
)

type Number interface {
	constraints.Signed | constraints.Float
}

type FibonacciParameters struct {
	Zero, Weighted bool
}
func Fibonacci[T Number](n int, params *FibonacciParameters) []T {

	var zero, weighted = false, false
	if params != nil {
		zero = params.Zero
		weighted = params.Weighted
	}

	var a, b T = 1, 1

	if zero {
		a, b = 0, 1
	} else {
		n -= 1
	}

	result := make([]T, n + 1)
	result[0] = a

	for i := 1; i < n + 1; i++ {
		a, b = b, a + b
		result[i] = a
	}

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

// // Slower version 
// func Factorial2[T Number](n int) T {
// 	if n > 0 {
// 		return T(n) * Factorial[T](n - 1)
// 	}
// 	return 1
// }

func Factorial[T Number](n int) T {
	var f = 1
	if n > 1 {
		for i := 2; i <= n; i++ {
			f *= i
		}
	}
	return T(f)
}

// // Slower version
// func Combination[T Number](n, k int) int {
//     if k == 0 {
// 		return T(1)
// 	}
//     return T((n * Choose(n - 1, k - 1)) / k)
// }

func Combination[T Number](n, k int) T {
	var nck float64 = 1

	if k > n {
		return 0
	} 
	
	if (n - k < k) {
		k = n - k
	}

	// we need to use float numberse becasue of division
	var _n, _k = float64(n), float64(k)

	for i := 1.; i <= _k; i++ {
		nck *= (_n + 1 - i) / i
	}

	// we must round result
	nck = math.Round(nck)

	return T(nck)
}

type PascalsTriangleParameters struct {
	Weighted, Inverse bool
}

func PascalsTriangle[T Number](n int, params *PascalsTriangleParameters) []T {
	var weighted, inverse = false, false
	if params != nil {
		weighted = params.Weighted
		inverse = params.Inverse
	}

	result := make([]T, n + 1)
	for i := 0; i <= n; i++ {
		result[i] = Combination[T](n, i)
	}

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

type SineWaveParameters struct {
	Weighted bool
}

func SineWave(n int, params *SineWaveParameters) []float64 {
	var weighted = false
	if params != nil {
		weighted = params.Weighted
	}

	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = math.Sin((float64(i) + 1) * math.Pi / (float64(n) + 1))
	}

	if weighted {
		sum := Sum(result)
		for i := 0; i < n; i++ {
			result[i] /= sum
		}
	}

	return result
}

type SymmetricTriangleParameters struct {
	Weighted bool
}

func SymmetricTriangle[T Number](n int, params *SymmetricTriangleParameters) []T {
	if n < 2 { n = 2 }

	var weighted = false
	if params != nil {
		weighted = params.Weighted
	}

	var triangle []T
	if n > 2 {
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

	if weighted {
		sum := Sum(triangle)
		for i := 0; i < len(triangle); i++ {
			triangle[i] /= sum
		}
	}

	return triangle
}

func Reverse[T any](values []T) {
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}
}

func Sum[T Number](values []T) T {
	var m T = 0
	for i := range values {
		m += values[i]
	}
	return m
}

func Arange[T Number](from T, to T, step T) []T {
	result := make([]T, int(math.Ceil(float64((to - from) / step))))
	for i := 0; from < to; i++ {
		result[i] = from; from += step
	}
	return result
}

func Sign[T Number](val T) T {
	if val < 0 {
		return -1
	} else if val > 0 {
		return 1
	}
	return 0
}

func Abs[T Number](val T) T {
	return Sign(val) * val
}

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

func ClipLower[T Number](value T, lower T) T {
	if value < lower {
		return lower
	}
	return value
}

func ClipUpper[T Number](value T, upper T) T {
	if value > upper {
		return upper
	}
	return value
}

func Clip[T Number](value T, lower T, upper T) T {
	return ClipUpper(ClipLower(value, lower), upper)
}