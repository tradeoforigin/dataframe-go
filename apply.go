package dataframe

import (
	"context"
)

// ApplyDataFrameFn is used by the Apply function when used with DataFrames.
// vals contains the values for the current row. The keys contain ints (index of Series) and strings (name of Series).
// The returned map must only contain what values you intend to update. The key can be a string (name of Series) or int (index of Series).
// If nil is returned, the existing values for the row are unchanged.
type ApplyDataFrameFn func(vals map[string]any, row, nRows int) map[string]any

// ApplySeriesFn is used by the Apply function when used with Series.
// val contains the value of the current row. The returned value is the updated value.
type ApplySeriesFn[T any] func(val T, row, nRows int) T

// ApplyDataFrame applies function to DataFrame. If ApplyOptions are set as `ApplyOptions { InPlace: true }`
// then dataframe is modified, otherwise new dataframe is returned.
func ApplyDataFrame(ctx context.Context, df *DataFrame, fn ApplyDataFrameFn, options ...ApplyOptions) (*DataFrame, error) {

	if fn == nil {
		panic("fn is required")
	}

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.Lock()
		defer df.Unlock()
	}

	var ndf *DataFrame

	if !opts.InPlace {
		// Create all series
		seriess := []SeriesAny{}
		for _, s := range df.Series {
			seriess = append(seriess, s.cloneAsEmpty())
		}

		// Create a new dataframe
		ndf = NewDataFrame(seriess...)
	}

	iterator := df.Iterator(IteratorOptions { InitialRow: 0, Step: 1, DontLock: true })

	for iterator.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		newVals := fn(iterator.Value, iterator.Index, iterator.Total)

		if opts.InPlace {
			df.UpdateRow(iterator.Index, newVals, dontLock)
		} else {
			ndf.Append(newVals, dontLock)
		}
	}

	if !opts.InPlace {
		return ndf, nil
	}

	return df, nil
}

// Apply applies function to DataFrame. If ApplyOptions are set as `ApplyOptions { InPlace: true }`
// then dataframe is modified, otherwise new dataframe is returned.
func (df *DataFrame) Apply(ctx context.Context, fn ApplyDataFrameFn, options ...ApplyOptions) (*DataFrame, error) {
	return ApplyDataFrame(ctx, df, fn, options...)
}

// ApplySeries applies filter function to series. If ApplyOptions are set as `ApplyOptions { InPlace: true }`
// then series is modified, otherwise new series is returned.
func ApplySeries[T any](ctx context.Context, s *Series[T], fn ApplySeriesFn[T], options ...ApplyOptions) (*Series[T], error) {

	if fn == nil {
		panic("fn is required")
	}

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	nRows := s.NRows(dontLock)

	var ns *Series[T]

	if !opts.InPlace {
		// Create a New Series
		ns = NewSeries[T](s.Name(dontLock), &SeriesInit{ Capacity: nRows })
	}

	iterator := s.Iterator(IteratorOptions { InitialRow: 0, Step: 1, DontLock: true })

	for iterator.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		newVal := fn(iterator.Value, iterator.Index, iterator.Total)

		if opts.InPlace {
			s.Update(iterator.Index, newVal, dontLock)
		} else {
			ns.Append([]T { newVal }, dontLock)
		}
	}

	if !opts.InPlace {
		return ns, nil
	}

	return s, nil
}

// Apply applies filter function to series. If ApplyOptions are set as `ApplyOptions { InPlace: true }`
// then series is modified, otherwise new series is returned.
func (s *Series[T]) Apply(ctx context.Context, fn ApplySeriesFn[T], options ...ApplyOptions) (*Series[T], error) {
	return ApplySeries(ctx, s, fn, options...)
}
