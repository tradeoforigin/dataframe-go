package csv

import (
	"github.com/tradeoforigin/dataframe-go"
)

type ConverterFn[T any] func (string) T

type ConverterAny interface {
	// function to instantiatiate series of type T
	series(string, *dataframe.SeriesInit) dataframe.SeriesAny
	// function for convert from string value to any type
	value(string) any
}

// Converter defines custom transformation from string to value of
// type T. To initialize new Converter use csv.NewConverter(ConverterFn[T])
type Converter[T any] struct {
	fn ConverterFn[T]
}

// Function to instantiate new Converter of type T. Converter transforms string
// into value of type T. There are predefined converters like csv.Float64, csv.Int,
// csv.Time, etc. 
// 
// Example:
//
//	var Float64 = NewConverter(
// 		func(s string) float64 {
// 			v, err := strconv.ParseFloat(s, 64)
// 			if err != nil {
// 				panic(err)
// 			}
// 			return v
// 		}
// 	)
func NewConverter[T any](fn ConverterFn[T]) Converter[T] {
	return Converter[T] { fn }
}

// Interface CovnverterAny function for auto instantiation series of type T
func (c Converter[T]) series(name string, init *dataframe.SeriesInit) dataframe.SeriesAny {
	return dataframe.NewSeries[T](name, init)
}

// Interface ConverterAny function to call converter.fn
func (c Converter[T]) value(s string) any {
	return c.fn(s)
}

