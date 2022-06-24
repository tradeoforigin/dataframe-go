package dataframe

import "context"

// FilterAction is the return value of FilterSeriesFn and FilterDataFrameFn.
type FilterAction int

const (
	// DROP is used to signify that a row must be dropped.
	DROP FilterAction = 0

	// KEEP is used to signify that a row must be kept.
	KEEP FilterAction = 1

	// CHOOSE is used to signify that a row must be kept.
	CHOOSE FilterAction = 1
)

// FilterSeriesFn is used by the Filter function to determine which rows are selected.
// val contains the value of the current row.
// If the function returns DROP, then the row is removed. If KEEP or CHOOSE is chosen, the row is kept.
type FilterSeriesFn[T any] func(val T, row, nRows int) (FilterAction, error)

// FilterDataFrameFn is used by the Filter function to determine which rows are selected.
// vals contains the values for the current row. The keys contain ints (index of Series) and strings (name of Series).
// If the function returns DROP, then the row is removed. If KEEP or CHOOSE is chosen, the row is kept.
type FilterDataFrameFn func(vals map[string]any, row, nRows int) (FilterAction, error)

// FilterSeries applies filter function to series. If FilterOptions are set as `FilterOptions { InPlace: true }`
// then series is modified, otherwise new series is returned.
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

	iterator := s.Iterator(IteratorOptions { InitialRow: 0, Step: 1, DontLock: true })

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
		ns := NewSeries[T](s.Name(dontLock), &SeriesInit{ Capacity: len(transfer) })
		for _, rowToTransfer := range transfer {
			val := s.Value(rowToTransfer, dontLock)
			ns.Append([] T { val }, dontLock)
		}
		return ns, nil
	}

	// Remove rows that need to be removed
	for idx := len(transfer) - 1; idx >= 0; idx-- {
		s.Remove(transfer[idx], dontLock)
	}

	return s, nil
}

// Filter applies filter function to series. If FilterOptions are set as `FilterOptions { InPlace: true }`
// then series is modified, otherwise new series is returned.
func (s *Series[T]) Filter(ctx context.Context, fn FilterSeriesFn[T], options ...FilterOptions) (*Series[T], error) {
	return FilterSeries(ctx, s, fn, options...)
}

// FilterDataFrame applies filter function to DataFrame. If FilterOptions are set as `FilterOptions { InPlace: true }`
// then dataframe is modified, otherwise new dataframe is returned.
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

	iterator := df.Iterator(IteratorOptions { InitialRow: 0, Step: 1, DontLock: true })

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
		// Create all series
		seriess := []SeriesAny{}
		for _, s := range df.Series {
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

// Filter applies filter function to DataFrame. If FilterOptions are set as `FilterOptions { InPlace: true }`
// then dataframe is modified, otherwise new dataframe is returned.
func (df *DataFrame) Filter(ctx context.Context, fn FilterDataFrameFn, options ...FilterOptions) (*DataFrame, error) {
	return FilterDataFrame(ctx, df, fn, options...)
}
