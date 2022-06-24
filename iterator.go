package dataframe

// Iterator function returns actual row, value for that row, 
// total number of elements and "not done" flag
type IteratorFn[T any] func() (int, T, int, bool)

// Iterator is an structure for iterating Series or DataFrames.
// When `Next()` is called, new Index and Value is filled until
// `notDone` is true.
type Iterator[T any] struct {
	iterator func() (int, T, int, bool)
	Index, Total int
	notDone bool
	Value T
}

// NewIterator creates Iterator instance with iterator function of
// type `IteratorFn[T any]`. Iterator function is called with 
// `iterator.Next()`
func NewIterator[T any](iterator IteratorFn[T]) Iterator[T] {
	return Iterator[T] { iterator: iterator }
}

// Function to iterate all values by iterator function. This function
// returns true if there is next value to read.
func (it *Iterator[T]) Next() bool {
	it.Index, it.Value, it.Total, it.notDone = it.iterator()
	return it.notDone
}

// Iternal function to cast iterator to Iterator[any]
func (it Iterator[T]) toAnyIterator() Iterator[any] {
	return Iterator[any] {
		iterator: func() (int, any, int, bool) {
			return it.iterator()
		},
	}
}