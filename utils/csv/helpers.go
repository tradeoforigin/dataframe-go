package csv

import (
	"strconv"
	"time"
)

// CSV converters for various data types used for parsing CSV fields into specific Go types.
// These converters are used to convert CSV string values into the appropriate data type 
// for each column in the resulting dataframe.

var (
	// String is a converter that returns the string value as-is.
	// It is used for string type fields in the CSV.
	String = NewConverter(
		func(s string) string {
			return s
		},
	)

	// Float64 is a converter for parsing float64 values from CSV strings.
	// If parsing fails, it will panic.
	Float64 = NewConverter(
		func(s string) float64 {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(err)
			}
			return v
		},
	)

	// Float32 is a converter for parsing float32 values from CSV strings.
	// If parsing fails, it will panic.
	Float32 = NewConverter(
		func(s string) float32 {
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				panic(err)
			}
			return float32(v)
		},
	)

	// Int64 is a converter for parsing int64 values from CSV strings.
	// If parsing fails, it will panic.
	Int64 = NewConverter(
		func(s string) int64 {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(err)
			}
			return v
		},
	)

	// Int32 is a converter for parsing int32 values from CSV strings.
	// If parsing fails, it will panic.
	Int32 = NewConverter(
		func(s string) int32 {
			v, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				panic(err)
			}
			return int32(v)
		},
	)

	// Int is a converter for parsing int values from CSV strings.
	// If parsing fails, it will panic.
	Int = NewConverter(
		func(s string) int {
			v, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				panic(err)
			}
			return int(v)
		},
	)

	// UInt64 is a converter for parsing uint64 values from CSV strings.
	// If parsing fails, it will panic.
	UInt64 = NewConverter(
		func(s string) uint64 {
			v, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				panic(err)
			}
			return v
		},
	)

	// UInt32 is a converter for parsing uint32 values from CSV strings.
	// If parsing fails, it will panic.
	UInt32 = NewConverter(
		func(s string) uint32 {
			v, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				panic(err)
			}
			return uint32(v)
		},
	)

	// UInt is a converter for parsing uint values from CSV strings.
	// If parsing fails, it will panic.
	UInt = NewConverter(
		func(s string) uint {
			v, err := strconv.ParseUint(s, 10, 0)
			if err != nil {
				panic(err)
			}
			return uint(v)
		},
	)

	// Bool is a converter for parsing bool values from CSV strings.
	// If parsing fails, it will panic.
	Bool = NewConverter(
		func(s string) bool {
			v, err := strconv.ParseBool(s)
			if err != nil {
				panic(err)
			}
			return v
		},
	)

	// Complex128 is a converter for parsing complex128 values from CSV strings.
	// If parsing fails, it will panic.
	Complex128 = NewConverter(
		func(s string) complex128 {
			v, err := strconv.ParseComplex(s, 128)
			if err != nil {
				panic(err)
			}
			return v
		},
	)

	// Complex64 is a converter for parsing complex64 values from CSV strings.
	// If parsing fails, it will panic.
	Complex64 = NewConverter(
		func(s string) complex64 {
			v, err := strconv.ParseComplex(s, 64)
			if err != nil {
				panic(err)
			}
			return complex64(v)
		},
	)

	// Time is a converter for parsing time.Time values from CSV strings.
	// It first attempts to parse the string in RFC3339 format, and if it fails,
	// it tries to parse the string as a Unix timestamp.
	// If both attempts fail, it will panic.
	Time = NewConverter(
		func(s string) time.Time {
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
)
