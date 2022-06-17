package dataframe

func (df *DataFrame) getSeriesAny(nameOrId any) SeriesAny {
	switch name := nameOrId.(type) {
	case string: return df.Series[df.MustNameToColumn(name)]
	}

	return df.Series[nameOrId.(int)]
}

func GetSeries[T any, U int|string](df *DataFrame, name U) *Series[T] {
	return df.getSeriesAny(name).(*Series[T])
}

func DefaultOptions[T any](o ...T) T {
	if len(o) > 0 {
		return o[0]
	}
	return *new(T)
}

func Range(r ...int) RangeOptions {
	if len(r) > 1 {
		return RangeOptions { r[0], &r[1] }
	}

	if len(r) == 1 {
		return RangeOptions { Start: r[0] }
	}

	return RangeOptions{}
}