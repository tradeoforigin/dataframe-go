package dataframe

import (
	"context"
	"sort"
)

// SortKey is the key to sort a Dataframe
type SortKey struct {

	// Key can be an int (position of series) or string (name of series).
	Key any

	// Desc can be set to sort in descending order.
	Desc bool

	seriesIndex int
}

type sorter struct {
	keys []SortKey
	df   *DataFrame
	ctx  context.Context
}

func (s *sorter) Len() int {
	return s.df.n
}

func (s *sorter) Less(i, j int) bool {

	if err := s.ctx.Err(); err != nil {
		panic(err)
	}

	for _, key := range s.keys {
		series := s.df.Series[key.seriesIndex]

		left := series.ValueAny(i)
		right := series.ValueAny(j)

		// Check if left and right are not equal
		if !series.IsEqualAnyFunc(left, right) {
			if key.Desc {
				// Sort in descending order
				return !series.IsLessThanAnyFunc(left, right)
			}
			return series.IsLessThanAnyFunc(left, right)
		}
	}

	return false
}

func (s *sorter) Swap(i, j int) {
	s.df.Swap(i, j, DontLock)
}

// Sort is used to sort the Dataframe according to different keys.
// It will return true if sorting was completed or false when the context is canceled.
func (df *DataFrame) Sort(ctx context.Context, keys []SortKey, options ...SortOptions) (completed bool) {
	if len(keys) == 0 {
		return true
	}

	defer func() {
		if x := recover(); x != nil {
			if x == context.Canceled || x == context.DeadlineExceeded {
				completed = false
			} else {
				panic(x)
			}
		}
	}()

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		// Default
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	// Clear seriesIndex from keys
	defer func() {
		for i := range keys {
			key := &keys[i]
			key.seriesIndex = 0
		}
	}()

	// Convert keys to index
	for i := range keys {
		key := &keys[i]

		name, ok := key.Key.(string)
		if ok {
			col, err := df.NameToColumn(name, dontLock)
			if err != nil {
				panic(err)
			}
			key.seriesIndex = col
		} else {
			key.seriesIndex = key.Key.(int)
		}
	}

	s := &sorter{
		keys: keys,
		df:   df,
		ctx:  ctx,
	}

	if opts.Stable {
		sort.Stable(s)
	} else {
		// Default
		sort.Sort(s)
	}

	return true
}