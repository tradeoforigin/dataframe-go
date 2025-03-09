package csv

import (
	"github.com/tradeoforigin/dataframe-go"
)

// ConverterFn defines a function type that takes a string as input and returns
// a value of type T. This function is used for converting CSV string values
// into a specific type during the CSV import process.
type ConverterFn[T any] func(string) T

// ConverterAny is an interface for a converter that can handle the instantiation
// of series of type T and provide a method for converting string values to type T.
type ConverterAny interface {
	// series creates a new series of type T and initializes it with the provided
	// name and dataframe SeriesInit.
	series(string, *dataframe.SeriesInit) dataframe.SeriesAny

	// value converts a string to a value of type T.
	value(string) any
}

// Converter is a type that transforms a string into a value of type T. To use
// this converter, you instantiate it with a function (ConverterFn) that defines
// how to convert strings to values of type T. There are predefined converters
// for common types like csv.Float64, csv.Int, csv.Time, etc.
type Converter[T any] struct {
	fn ConverterFn[T]
}

// NewConverter initializes a new Converter for a specific type T. The function
// provided (ConverterFn) defines how to convert a string value to the type T.
//
// Example:
//
//	// Define a converter for float64 values
//	var Float64 = NewConverter(
//		func(s string) float64 {
//			v, err := strconv.ParseFloat(s, 64)
//			if err != nil {
//				panic(err)
//			}
//			return v
//		}
//	)
func NewConverter[T any](fn ConverterFn[T]) Converter[T] {
	return Converter[T]{fn}
}

// series is an implementation of the ConverterAny interface, which creates
// a new series of type T with the provided name and initialization options.
func (c Converter[T]) series(name string, init *dataframe.SeriesInit) dataframe.SeriesAny {
	return dataframe.NewSeries[T](name, init)
}

// value is an implementation of the ConverterAny interface, which uses the
// converter function (fn) to transform the string into a value of type T.
func (c Converter[T]) value(s string) any {
	return c.fn(s)
}
