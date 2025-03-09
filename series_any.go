package dataframe

import "context"

// valueAny returns the value of a particular row.
func (s *Series[T]) valueAny(row int, options ...Options) any {
	return s.Value(row, options...)
}

// PrependAny is used to set a value to the beginning of the
// Series.
func (s *Series[T]) prependAny(val any, options ...Options) {
	switch v := val.(type) {
	case []T:
		s.Prepend(v, options...)
	case T:
		s.Prepend([]T{v}, options...)
	default:
		// raise panic
		_ = val.(T)
	}
}

// AppendAny is used to set a value to the end of the Series.
func (s *Series[T]) appendAny(val any, options ...Options) int {
	switch v := val.(type) {
	case []T:
		return s.Append(v, options...)
	case T:
		return s.Append([]T{v}, options...)
	}

	return s.Append([]T{val.(T)}, options...)
}

// InsertAny is used to set a value at an arbitrary row in
// the Series. All existing values from that row onwards
// are shifted by 1.
func (s *Series[T]) insertAny(row int, val any, options ...Options) {
	switch v := val.(type) {
	case []T:
		s.Insert(row, v, options...)
	case T:
		s.Insert(row, []T{v}, options...)
	default:
		// raise panic
		_ = val.(T)
	}
}

// UpdateAny is used to update the value of a particular row.
func (s *Series[T]) updateAny(row int, val any, options ...Options) {
	switch v := val.(type) {
	case T:
		s.Update(row, v, options...)
	default:
		// raise panic
		_ = val.(T)
	}
}

// IteratorAny will return a iterator that can be used to iterate through all the values.
func (s *Series[T]) iteratorAny(options ...IteratorOptions) Iterator[any] {
	return s.Iterator(options...).toAnyIterator()
}

// IsEqualAnyFunc	returns true if a is equal to b.
func (s *Series[T]) isEqualAnyFunc(a, b any) bool {
	return s.isEqualFunc(a.(T), b.(T))
}

// IsLessThanAnyFunc	returns true if a is less than b.
func (s *Series[T]) isLessThanAnyFunc(a, b any) bool {
	return s.isLessThanFunc(a.(T), b.(T))
}

// SetIsEqualAnyFunc	sets a function which can be used to determine
// if 2 values in the Series are equal.
func (s *Series[T]) setIsEqualAnyFunc(f CompareFn[any]) {
	s.SetIsEqualFunc(func(f1, f2 T) bool {
		return f(f1, f2)
	})
}

// SetIsLessThanAnyFunc	sets a function which can be used to determine
// if a value is less than another in the Series.
func (s *Series[T]) setIsLessThanAnyFunc(f CompareFn[any]) {
	s.SetIsLessThanFunc(func(f1, f2 T) bool {
		return f(f1, f2)
	})
}

// CopyAny will create a new copy of the Series.
// It is recommended that you lock the Series before attempting
// to Copy.
func (s *Series[T]) copyAny(options ...RangeOptions) SeriesAny {
	return s.Copy(options...)
}

func (s *Series[T]) cloneAsEmpty(size ...int) SeriesAny {
	var _size, _capacity = len(s.values), len(s.values)

	if len(size) > 1 {
		_size, _capacity = size[0], size[1]
	} else if len(size) == 1 {
		_size = size[0]
	}

	if _size > _capacity {
		_capacity = _size
	}

	return &Series[T]{
		name:         s.name,
		typeT:        s.typeT,
		values:       make([]T, _size, _capacity),
		valFormatter: DefaultValueFormatter,
		isEqualFunc:  IsEqualDefaultFunc[T],
	}
}

// FillRandAny will fill a Series with random data.
func (s *Series[T]) fillRandAny(rnd RandFn[any]) {
	s.FillRand(func() T {
		return rnd().(T)
	})
}

// IsEqualAny returns true if s2's values are equal to s.
func (s *Series[T]) isEqualAny(ctx context.Context, s2 SeriesAny, options ...IsEqualOptions) (bool, error) {
	return s.IsEqual(ctx, s2.(*Series[T]), options...)
}

var _ SeriesAny = (*Series[any])(nil)
