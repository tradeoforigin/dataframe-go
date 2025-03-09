package dataframe

import (
	"context"
)

// ApplyDataFrameFn is a function type used by the Apply function when working with DataFrames.
// It takes a map containing the values of the current row, where the keys are either ints (index of series) or strings (name of series).
// The returned map must only contain the values you intend to update, using the same keys (either string or int).
// If nil is returned, the existing values for the row are not changed.
type ApplyDataFrameFn func(vals map[string]any, row, nRows int) map[string]any

// ApplySeriesFn is a function type used by the Apply function when working with series.
// It takes a single value from the series and returns the updated value for that row.
type ApplySeriesFn[T any] func(val T, row, nRows int) T

// ApplyDataFrame applies the provided function to each row of the dataFrame.
// If ApplyOptions are set with `ApplyOptions { InPlace: true }`, the dataframe is modified in place; otherwise, a new dataframe is returned.
func ApplyDataFrame(ctx context.Context, df *dataFrame, fn ApplyDataFrameFn, options ...ApplyOptions) (DataFrame, error) {

	if fn == nil {
		panic("fn is required")
	}

	opts := DefaultOptions(options...)

	// Lock the dataframe if necessary
	if !opts.DontLock {
		df.Lock()
		defer df.Unlock()
	}

	var ndf DataFrame

	if !opts.InPlace {
		// Create a new dataframe if InPlace is false
		seriess := []SeriesAny{}
		for _, s := range df.series {
			seriess = append(seriess, s.cloneAsEmpty())
		}

		// Create a new dataframe
		ndf = NewDataFrame(seriess...)
	}

	// Iterate over the rows in the dataframe
	iterator := df.Iterator(IteratorOptions{InitialRow: 0, Step: 1, DontLock: true})

	// Apply the function to each row
	for iterator.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		newVals := fn(iterator.Value, iterator.Index, iterator.Total)

		// Update the dataframe in place or create a new dataframe
		if opts.InPlace {
			df.UpdateRow(iterator.Index, newVals, dontLock)
		} else {
			ndf.Append(newVals, dontLock)
		}
	}

	// Return the updated dataframe
	if !opts.InPlace {
		return ndf, nil
	}

	return df, nil
}

// Apply applies a function to a dataframe. If ApplyOptions are set with `ApplyOptions { InPlace: true }`, 
// the dataframe is modified in place; otherwise, a new dataframe is returned.
func (df *dataFrame) Apply(ctx context.Context, fn ApplyDataFrameFn, options ...ApplyOptions) (DataFrame, error) {
	return ApplyDataFrame(ctx, df, fn, options...)
}

// ApplySeries applies the provided function to each value in the series. 
// If ApplyOptions are set with `ApplyOptions { InPlace: true }`, the series is modified in place; 
// otherwise, a new series is returned.
func ApplySeries[T any](ctx context.Context, s Series[T], fn ApplySeriesFn[T], options ...ApplyOptions) (Series[T], error) {

	if fn == nil {
		panic("fn is required")
	}

	opts := DefaultOptions(options...)

	// Lock the series if necessary
	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	nRows := s.NRows(dontLock)

	var ns Series[T]

	// Create a new series if InPlace is false
	if !opts.InPlace {
		ns = NewSeries[T](s.Name(dontLock), &SeriesInit{Capacity: nRows})
	}

	// Iterate over the rows in the series
	iterator := s.Iterator(IteratorOptions{InitialRow: 0, Step: 1, DontLock: true})

	// Apply the function to each row
	for iterator.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		newVal := fn(iterator.Value, iterator.Index, iterator.Total)

		// Update the series in place or create a new series
		if opts.InPlace {
			s.Update(iterator.Index, newVal, dontLock)
		} else {
			ns.Append([]T{newVal}, dontLock)
		}
	}

	// Return the updated series
	if !opts.InPlace {
		return ns, nil
	}

	return s, nil
}

// Apply applies the provided function to a series. If ApplyOptions are set with `ApplyOptions { InPlace: true }`, 
// the series is modified in place; otherwise, a new series is returned.
func (s *series[T]) Apply(ctx context.Context, fn ApplySeriesFn[T], options ...ApplyOptions) (Series[T], error) {
	return ApplySeries(ctx, Series[T](s), fn, options...)
}
