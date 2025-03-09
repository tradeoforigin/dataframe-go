package dataframe

import (
	"context"
	"sort"
)

// SortKey represents a key used to sort a DataFrame. It can either be an integer (position of series)
// or a string (name of the series) that determines the sorting criterion.
type SortKey struct {

	// Key can be an int (position of series) or string (name of series).
	Key any

	// Desc is a boolean that indicates whether the sort should be in descending order.
	Desc bool

	// seriesIndex is the internal index for accessing the series in the DataFrame.
	seriesIndex int
}

// sorter is an internal type used for sorting the DataFrame. It implements the sort.Interface
// to allow sorting based on the provided `SortKey` values.
type sorter struct {
	keys []SortKey
	df   *DataFrame
	ctx  context.Context
}

// Len returns the number of elements in the DataFrame, implementing sort.Interface.
func (s *sorter) Len() int {
	return s.df.n
}

// Less compares the elements at indices i and j in the DataFrame. It uses the keys to determine
// the sorting order for each row. It supports sorting in both ascending and descending order.
func (s *sorter) Less(i, j int) bool {

	// Check for context cancellation.
	if err := s.ctx.Err(); err != nil {
		panic(err)
	}

	// Loop through the keys and compare the corresponding series values for rows i and j.
	for _, key := range s.keys {
		series := s.df.series[key.seriesIndex]

		left := series.valueAny(i)
		right := series.valueAny(j)

		// Check if the values at i and j are not equal
		if !series.isEqualAnyFunc(left, right) {
			// If Desc is true, sort in descending order
			if key.Desc {
				return !series.isLessThanAnyFunc(left, right)
			}
			// Default: sort in ascending order
			return series.isLessThanAnyFunc(left, right)
		}
	}

	// If all values are equal, maintain the current order.
	return false
}

// Swap swaps the elements with indices i and j in the DataFrame, implementing sort.Interface.
func (s *sorter) Swap(i, j int) {
	s.df.Swap(i, j, DontLock)
}

// Sort is used to sort the DataFrame based on the provided keys. It will sort the DataFrame's rows
// according to the values in the series specified by the keys.
// If the context is canceled, it returns false to indicate that the sorting was not completed.
func (df *DataFrame) Sort(ctx context.Context, keys []SortKey, options ...SortOptions) (completed bool) {
	if len(keys) == 0 {
		return true
	}

	// Recover from any panic caused by context cancellation.
	defer func() {
		if x := recover(); x != nil {
			if x == context.Canceled || x == context.DeadlineExceeded {
				completed = false
			} else {
				panic(x)
			}
		}
	}()

	// Get options with default values if not provided
	opts := DefaultOptions(options...)

	// Lock the DataFrame if the options specify not to skip locking.
	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	// Clear seriesIndex from keys after sorting is done
	defer func() {
		for i := range keys {
			key := &keys[i]
			key.seriesIndex = 0
		}
	}()

	// Convert the keys to their respective series index based on name or position
	for i := range keys {
		key := &keys[i]

		// If the key is a string (name of series), resolve it to an index
		name, ok := key.Key.(string)
		if ok {
			col, err := df.NameToColumn(name, dontLock)
			if err != nil {
				panic(err)
			}
			key.seriesIndex = col
		} else {
			// If the key is already an integer (position), use it directly
			key.seriesIndex = key.Key.(int)
		}
	}

	// Create a sorter and perform the sort operation based on the context
	s := &sorter{
		keys: keys,
		df:   df,
		ctx:  ctx,
	}

	// Sort the DataFrame using either stable or unstable sort
	if opts.Stable {
		sort.Stable(s)
	} else {
		sort.Sort(s)
	}

	// Return true if sorting is completed
	return true
}
