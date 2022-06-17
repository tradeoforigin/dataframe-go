package types

import (
	"dataframe-go"
	"time"
)

type Series[T any] 		  dataframe.Series[T]

type DataFrame 			= dataframe.DataFrame
type SeriesAny 			= dataframe.SeriesAny

type SeriesFloat64 		= dataframe.Series[float64]
type SeriesFloat32 		= dataframe.Series[float32]

type SeriesInt 			= dataframe.Series[int]
type SeriesInt64 		= dataframe.Series[int64]
type SeriesInt32 		= dataframe.Series[int32]

type SeriesUInt 		= dataframe.Series[uint]
type SeriesUInt64 		= dataframe.Series[uint64]
type SeriesUInt32 		= dataframe.Series[uint32]

type SeriesComplex64 	= dataframe.Series[complex64]
type SeriesComplex128 	= dataframe.Series[complex128]

type SeriesString 		= dataframe.Series[string]

type SeriesTime 		= dataframe.Series[time.Time]

type SeriesMixed 		= dataframe.Series[any]
