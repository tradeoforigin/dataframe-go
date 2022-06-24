package dataframe

// getSeriesAny helps return series for given name or id as 
// `SeriesAny` type
func (df *DataFrame) getSeriesAny(nameOrId any) SeriesAny {
	switch name := nameOrId.(type) {
	case string: return df.Series[df.MustNameToColumn(name)]
	}

	return df.Series[nameOrId.(int)]
}

// GetSeries helps get series of `DataFrame` as a series of concrete type
func GetSeries[T any, U int|string](df *DataFrame, name U) *Series[T] {
	return df.getSeriesAny(name).(*Series[T])
}

// DefaultOptions is helper function to resolve variadic options.
func DefaultOptions[T any](o ...T) T {
	if len(o) > 0 {
		return o[0]
	}
	return *new(T)
}

// Range is helper function for creating RangeOptions. 
// 
// Example:
//	r1 := Range(0, 10) // Equivalent to RangeOptions { Start: 0, End: &[]int { 10 }[0]}
//  r2 := Range(10) // Equivalent to RangeOptions { Start: 10 }
func Range(r ...int) RangeOptions {
	if len(r) > 1 {
		return RangeOptions { r[0], &r[1] }
	}

	if len(r) == 1 {
		return RangeOptions { Start: r[0] }
	}

	return RangeOptions{}
}