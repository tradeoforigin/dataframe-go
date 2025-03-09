package dataframe

import "context"

// FilterAction represents the possible actions to be taken during a filter operation.
// It is used by `FilterSeriesFn` and `FilterDataFrameFn` to decide whether to KEEP, DROP, or CHOOSE a row.
type FilterAction int

const (
	// DROP signifies that the row should be removed from the *Series or dataframe.
	DROP FilterAction = 0

	// KEEP signifies that the row should be retained in the *Series or dataframe.
	KEEP FilterAction = 1

	// CHOOSE is synonymous with KEEP and signifies that the row should be retained.
	CHOOSE FilterAction = 1
)

// FilterSeriesFn defines a filter function for a *Series. It is used in the `FilterSeries` function.
// This function is called for each row in the *Series, where `val` is the value of the current row,
// `row` is the row index, and `nRows` is the total number of rows in the *Series.
// The function should return a `FilterAction` indicating whether the row should be kept or dropped, and possibly an error.
type FilterSeriesFn[T any] func(val T, row, nRows int) (FilterAction, error)

// FilterDataFrameFn defines a filter function for a DataFrame. It is used in the `FilterDataFrame` function.
// This function is called for each row in the dataframe, where `vals` is a map of values for the current row,
// `row` is the row index, and `nRows` is the total number of rows in the dataframe.
// The function should return a `FilterAction` indicating whether the row should be kept or dropped, and possibly an error.
type FilterDataFrameFn func(vals map[string]any, row, nRows int) (FilterAction, error)

// FilterSeries applies the given filter function to a *Series. If `FilterOptions` is set with `InPlace: true`,
// the *Series will be modified directly. Otherwise, a new *Series is returned.
// The filter function is applied to each row, and rows that return `DROP` will be removed.
// Rows that return `KEEP` or `CHOOSE` will be kept.
func FilterSeries[T any](ctx context.Context, s *Series[T], fn FilterSeriesFn[T], options ...FilterOptions) (*Series[T], error) {

	if fn == nil {
		panic("fn is required")
	}

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	transfer := []int{}

	iterator := s.Iterator(IteratorOptions{InitialRow: 0, Step: 1, DontLock: true})

	for iterator.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		fa, err := fn(iterator.Value, iterator.Index, iterator.Total)
		if err != nil {
			return nil, err
		}

		if fa == DROP {
			if opts.InPlace {
				transfer = append(transfer, iterator.Index)
			}
		} else if fa == KEEP || fa == CHOOSE {
			if !opts.InPlace {
				transfer = append(transfer, iterator.Index)
			}
		} else {
			panic("unrecognized FilterAction returned by fn")
		}
	}

	if !opts.InPlace {
		ns := NewSeries[T](s.Name(dontLock), &SeriesInit{Capacity: len(transfer)})
		for _, rowToTransfer := range transfer {
			val := s.Value(rowToTransfer, dontLock)
			ns.Append([]T{val}, dontLock)
		}
		return ns, nil
	}

	// Remove rows that need to be removed
	for idx := len(transfer) - 1; idx >= 0; idx-- {
		s.Remove(transfer[idx], dontLock)
	}

	return s, nil
}

// Filter applies the given filter function to the *Series. This is a method of the `Series[T]` type.
// If `FilterOptions` is set with `InPlace: true`, the *Series will be modified directly. Otherwise, a new *Series is returned.
func (s *Series[T]) Filter(ctx context.Context, fn FilterSeriesFn[T], options ...FilterOptions) (*Series[T], error) {
	return FilterSeries(ctx, s, fn, options...)
}

// FilterDataFrame applies the given filter function to a dataframe. If `FilterOptions` is set with `InPlace: true`,
// the dataframe will be modified directly. Otherwise, a new dataframe is returned.
// The filter function is applied to each row, and rows that return `DROP` will be removed.
// Rows that return `KEEP` or `CHOOSE` will be kept.
func FilterDataFrame(ctx context.Context, df *DataFrame, fn FilterDataFrameFn, options ...FilterOptions) (*DataFrame, error) {

	if fn == nil {
		panic("fn is required")
	}

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.Lock()
		defer df.Unlock()
	}

	transfer := []int{}

	iterator := df.Iterator(IteratorOptions{InitialRow: 0, Step: 1, DontLock: true})

	for iterator.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		fa, err := fn(iterator.Value, iterator.Index, iterator.Total)
		if err != nil {
			return nil, err
		}

		if fa == DROP {
			if opts.InPlace {
				transfer = append(transfer, iterator.Index)
			}
		} else if fa == KEEP || fa == CHOOSE {
			if !opts.InPlace {
				transfer = append(transfer, iterator.Index)
			}
		} else {
			panic("unrecognized FilterAction returned by fn")
		}
	}

	if !opts.InPlace {
		// Create all *Series
		seriess := []SeriesAny{}
		for _, s := range df.Series() {
			seriess = append(seriess, s.cloneAsEmpty(len(transfer), len(transfer)))
		}

		// Create a new dataframe
		ndf := NewDataFrame(seriess...)

		for _, rowToTransfer := range transfer {
			vals := df.Row(rowToTransfer, dontLock)
			ndf.Append(vals, dontLock)
		}
		return ndf, nil
	}

	// Remove rows that need to be removed
	for idx := len(transfer) - 1; idx >= 0; idx-- {
		df.Remove(transfer[idx], dontLock)
	}

	return df, nil
}

// Filter applies the given filter function to the dataframe. This is a method of the `DataFrame` type.
// If `FilterOptions` is set with `InPlace: true`, the dataframe will be modified directly. Otherwise, a new dataframe is returned.
func (df *DataFrame) Filter(ctx context.Context, fn FilterDataFrameFn, options ...FilterOptions) (*DataFrame, error) {
	return FilterDataFrame(ctx, df, fn, options...)
}
