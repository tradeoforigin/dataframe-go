package dataframe

import (
	"math"
	"math/rand"
	"time"
)

type RandFn[T any] func() T

// RandFillerFloat64 is a helper function to generate random values for 
// filling data in a Series[float64]. It returns a RandFn[float64] function 
// that generates random `float64` values. The function can optionally return 
// NaN values based on the specified probability.
//
// Parameters:
//   probNil (optional, float64): The probability of returning NaN (Not a Number) 
//   values. If provided, it must be a value between 0 and 1. The higher the value, 
//   the more likely NaN will be returned. If not provided, NaN will never be returned.
//
// Returns:
//   RandFn[float64]: A function that generates random float64 values. When called, 
//   it will either return a random float64 or NaN based on the provided probability.
//
// Example:
//   randFiller := RandFillerFloat64(0.1) // 10% chance of NaN
//   value := randFiller()                // Call to generate a random value
func RandFillerFloat64(probNil ...float64) RandFn[float64] {
	var pNil float64

	if len(probNil) > 0 {
		pNil = probNil[0]
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	return func() float64 {
		if rnd.Float64() < pNil {
			return math.NaN()
		}

		return rnd.Float64()
	}
}
