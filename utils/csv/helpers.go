package csv

import (
	"strconv"
	"time"
)

// CSV converted for string types
var String = NewConverter(
	func(s string) string {
		return s
	},
)

// CSV converted for float64 types
var Float64 = NewConverter(
	func(s string) float64 {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
)

// CSV converted for float32 types
var Float32 = NewConverter(
	func(s string) float32 {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			panic(err)
		}
		return float32(v)
	},
)

// CSV converted for int64 types
var Int64 = NewConverter(
	func(s string) int64 {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
)

// CSV converted for int32 types
var Int32 = NewConverter(
	func(s string) int32 {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			panic(err)
		}
		return int32(v)
	},
)

// CSV converted for int types
var Int = NewConverter(
	func(s string) int {
		v, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			panic(err)
		}
		return int(v)
	},
)

// CSV converted for uint64 types
var UInt64 = NewConverter(
	func(s string) uint64 {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
)

// CSV converted for uint32 types
var UInt32 = NewConverter(
	func(s string) uint32 {
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			panic(err)
		}
		return uint32(v)
	},
)

// CSV converted for uint types
var UInt = NewConverter(
	func(s string) uint {
		v, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			panic(err)
		}
		return uint(v)
	},
)

// CSV converted for bool types
var Bool = NewConverter(
	func(s string) bool {
		v, err := strconv.ParseBool(s)
		if err != nil {
			panic(err)
		}
		return v
	},
)

// CSV converted for complex128 types
var Complex128 = NewConverter(
	func(s string) complex128 {
		v, err := strconv.ParseComplex(s, 128)
		if err != nil {
			panic(err)
		}
		return v
	},
)

// CSV converted for complex64 types
var Complex64 = NewConverter(
	func(s string) complex64 {
		v, err := strconv.ParseComplex(s, 64)
		if err != nil {
			panic(err)
		}
		return complex64(v)
	},
)

// CSV converted for time.Time types
var Time = NewConverter(
	func (s string) time.Time {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			sec, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(err)
			}
			return time.Unix(sec, 0)
		}
		return t
	},
)

