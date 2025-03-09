package types

import (
	"time"

	"github.com/tradeoforigin/dataframe-go"
)

// DataFrame is a type alias for dataframe.DataFrame. It represents a data structure
// that holds tabular data with rows and columns, where columns can have different types.
type DataFrame = dataframe.DataFrame

// SeriesFloat64 is a type alias for dataframe.Series[float64].
// It represents a series of float64 values in a dataframe.
type SeriesFloat64 = dataframe.Series[float64]

// SeriesFloat32 is a type alias for dataframe.Series[float32].
// It represents a series of float32 values in a dataframe.
type SeriesFloat32 = dataframe.Series[float32]

// SeriesInt is a type alias for dataframe.Series[int].
// It represents a series of int values in a dataframe.
type SeriesInt = dataframe.Series[int]

// SeriesInt64 is a type alias for dataframe.Series[int64].
// It represents a series of int64 values in a dataframe.
type SeriesInt64 = dataframe.Series[int64]

// SeriesInt32 is a type alias for dataframe.Series[int32].
// It represents a series of int32 values in a dataframe.
type SeriesInt32 = dataframe.Series[int32]

// SeriesUInt is a type alias for dataframe.Series[uint].
// It represents a series of uint values in a dataframe.
type SeriesUInt = dataframe.Series[uint]

// SeriesUInt64 is a type alias for dataframe.Series[uint64].
// It represents a series of uint64 values in a dataframe.
type SeriesUInt64 = dataframe.Series[uint64]

// SeriesUInt32 is a type alias for dataframe.Series[uint32].
// It represents a series of uint32 values in a dataframe.
type SeriesUInt32 = dataframe.Series[uint32]

// SeriesComplex64 is a type alias for dataframe.Series[complex64].
// It represents a series of complex64 values in a dataframe.
type SeriesComplex64 = dataframe.Series[complex64]

// SeriesComplex128 is a type alias for dataframe.Series[complex128].
// It represents a series of complex128 values in a dataframe.
type SeriesComplex128 = dataframe.Series[complex128]

// SeriesString is a type alias for dataframe.Series[string].
// It represents a series of string values in a dataframe.
type SeriesString = dataframe.Series[string]

// SeriesTime is a type alias for dataframe.Series[time.Time].
// It represents a series of time.Time values in a dataframe.
type SeriesTime = dataframe.Series[time.Time]

// SeriesMixed is a type alias for dataframe.Series[any].
// It represents a series that can hold values of any type in a dataframe.
type SeriesMixed = dataframe.Series[any]
