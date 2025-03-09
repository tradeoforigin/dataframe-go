package dataframe

// getSeriesAny returns the series of type SeriesAny for a given name or id.
// It accepts either a string (name) or an integer (id) and returns the corresponding series.
func (df *DataFrame) getSeriesAny(nameOrId any) SeriesAny {
	switch name := nameOrId.(type) {
	case string:
		// If the input is a string, it resolves the column by its name
		return df.series[df.MustNameToColumn(name)]
	}

	// If the input is an integer (id), it resolves the column by its index
	return df.series[nameOrId.(int)]
}

// GetSeries is a generic function that returns a typed series for the given `name`
// in the data frame `df`. The type `T` represents the type of the values in the series,
// and `U` can either be an integer (column index) or a string (column name).
func GetSeries[T any, U int | string](df *DataFrame, name U) *Series[T] {
	// Retrieves the series and casts it to the expected concrete type *series[T]
	return df.getSeriesAny(name).(*Series[T])
}

// DefaultOptions is a helper function to resolve variadic options. If options are provided,
// it returns the first element; otherwise, it returns the zero value of type T.
func DefaultOptions[T any](o ...T) T {
	if len(o) > 0 {
		return o[0] // Returns the first option if provided
	}
	return *new(T) // Returns the zero value of type T if no options are provided
}

// Range is a helper function for creating RangeOptions. It can create a range from a start index
// to an end index. If only one argument is provided, it is treated as the start index, and the
// end index defaults to `nil`.
//
// Example:
// 		r1 := Range(0, 10) // Equivalent to RangeOptions { Start: 0, End: &[]int { 10 }[0]}
// 		r2 := Range(10) // Equivalent to RangeOptions { Start: 10 }
func Range(r ...int) RangeOptions {
	if len(r) > 1 {
		// If two arguments are provided, it creates a RangeOptions with both Start and End
		return RangeOptions{r[0], &r[1]}
	}

	if len(r) == 1 {
		// If only one argument is provided, it sets the Start and leaves End as nil
		return RangeOptions{Start: r[0]}
	}

	// If no arguments are provided, it returns a default RangeOptions with Start: 0 and End: nil
	return RangeOptions{}
}
