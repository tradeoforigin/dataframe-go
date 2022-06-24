package dataframe

import (
	"math"
	"math/rand"
	"time"
)

type RandFn[T any] func() T

// RandFillerFloat64 is helper function to fill data of *Series[float64]
// randomly. probNil is an optional parameter which indicates probability
// of NaN value as a return.
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