package csv

import (
	"strconv"
	"time"
)

var String = NewConverter(
	func(s string) string {
		return s
	},
)

var Float64 = NewConverter(
	func(s string) float64 {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
)

var Float32 = NewConverter(
	func(s string) float32 {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			panic(err)
		}
		return float32(v)
	},
)

var Int64 = NewConverter(
	func(s string) int64 {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
)

var Int32 = NewConverter(
	func(s string) int32 {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			panic(err)
		}
		return int32(v)
	},
)

var Int = NewConverter(
	func(s string) int {
		v, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			panic(err)
		}
		return int(v)
	},
)

var UInt64 = NewConverter(
	func(s string) uint64 {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	},
)

var UInt32 = NewConverter(
	func(s string) uint32 {
		v, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			panic(err)
		}
		return uint32(v)
	},
)

var UInt = NewConverter(
	func(s string) uint {
		v, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			panic(err)
		}
		return uint(v)
	},
)

var Bool = NewConverter(
	func(s string) bool {
		v, err := strconv.ParseBool(s)
		if err != nil {
			panic(err)
		}
		return v
	},
)

var Complex128 = NewConverter(
	func(s string) complex128 {
		v, err := strconv.ParseComplex(s, 128)
		if err != nil {
			panic(err)
		}
		return v
	},
)

var Complex64 = NewConverter(
	func(s string) complex64 {
		v, err := strconv.ParseComplex(s, 64)
		if err != nil {
			panic(err)
		}
		return complex64(v)
	},
)

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

