package dataframe

import (
	"math"
	"math/rand"
	"time"
)

type RandFn[T any] func() T

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