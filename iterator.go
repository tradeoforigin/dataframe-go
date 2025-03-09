package dataframe

// IteratorFn is a function type that returns the current row index, value for that row,
// total number of elements in the collection, and a "not done" flag indicating if there are
// more elements to iterate over.
type IteratorFn[T any] func() (int, T, int, bool)

// Iterator is a structure for iterating over Series or DataFrames. It uses the provided
// iterator function to fetch the current index, value, total number of elements, and a flag 
// indicating whether there are more elements to read.
//
// When the `Next()` method is called, the iterator will update the index, value, total, 
// and the "not done" flag, continuing the iteration until "not done" is false.
type Iterator[T any] struct {
	// iterator is the function that is called to get the current state of the iterator
	iterator func() (int, T, int, bool)

	// Index is the current index in the iteration.
	Index int

	// Total is the total number of elements in the Series or DataFrame.
	Total int

	// notDone is a flag that indicates whether the iteration has more values to return.
	notDone bool

	// Value is the current value at the row being iterated.
	Value T
}

// NewIterator creates a new Iterator instance with the provided iterator function 
// of type `IteratorFn[T]`. The iterator function will be called with each call to `Next()`.
func NewIterator[T any](iterator IteratorFn[T]) Iterator[T] {
	return Iterator[T] { iterator: iterator }
}

// Next advances the iterator and fills in the current row's index, value, total number of 
// elements, and the "not done" flag. It returns `true` if there are more values to read,
// or `false` if the iteration is complete.
func (it *Iterator[T]) Next() bool {
	// Call the iterator function to update the iterator state
	it.Index, it.Value, it.Total, it.notDone = it.iterator()
	return it.notDone
}

// toAnyIterator is an internal function that casts the current Iterator[T] to an 
// Iterator[any], where the value type is generalized to `any` for more flexible iteration.
func (it Iterator[T]) toAnyIterator() Iterator[any] {
	// Return a new Iterator[any] by wrapping the existing iterator function
	return Iterator[any] {
		iterator: func() (int, any, int, bool) {
			// Call the original iterator function and return the result with the any type
			return it.iterator()
		},
	}
}
