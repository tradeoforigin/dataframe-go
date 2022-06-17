package dataframe

// Iterator function returns actual row, value for that row, 
// total number of elements and "has next value" flag
type IteratorFn[T any] func() (int, T, int, bool)

type Iterator[T any] struct {
	iterator func() (int, T, int, bool)
	Index, Total int
	notDone bool
	Value T
}

func NewIterator[T any](iterator IteratorFn[T]) Iterator[T] {
	return Iterator[T] { iterator: iterator }
}

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