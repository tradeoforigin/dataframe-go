package dataframe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type series[T any] struct {
	valFormatter ValueToStringFormatter

	name, typeT string

	isEqualFunc, isLessThanFunc CompareFn[T]

	values []T

	mutex sync.RWMutex
}

// NewSeries creates a series of type T with defined name. Size of the series
// can be prealocated by passing `init`. series can also by filled by data passed
// as vals.
//
// Example:
//
// x := NewSeries[float64]("x", nil, 1, 2, 3)
// y := NewSeries("y", nil, 1., 2., 3.)
func NewSeries[T any](name string, init *SeriesInit, vals ...T) *series[T] {
	s := &series[T]{
		name:         name,
		valFormatter: DefaultValueFormatter,
		typeT:        formatType[T](),
		isEqualFunc:  IsEqualDefaultFunc[T],
		values:       []T{},
	}

	var size, capacity int

	if init != nil {
		size = init.Size
		capacity = init.Capacity
	}

	if size < len(vals) {
		size = len(vals)
	}

	if capacity < size {
		capacity = size
	}

	s.values = make([]T, size, capacity)

	copy(s.values, vals)

	s.fillDefault(s.values, len(vals), size)
	return s
}

// Fill values as NaN if series is type of float32 or float 64
func (s *series[T]) fillDefault(vals any, lVals, size int) {
	// s.nilCount = s.nilCount + size - lVals
	switch v := vals.(type) {
	case []float64:
		for i := lVals; i < size; i++ {
			v[i] = math.NaN()
		}
	case []float32:
		for i := lVals; i < size; i++ {
			v[i] = float32(math.NaN())
		}
	}
}

// Name returns the series name.
func (s *series[T]) Name(options ...Options) string {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.RLock()
		defer s.RUnlock()
	}

	return s.name
}

// Rename renames the series.
func (s *series[T]) Rename(n string, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.RLock()
		defer s.RUnlock()
	}

	s.name = n
}

func (s *series[T]) Type() string {
	return s.typeT
}

// NRows returns how many rows the series contains.
func (s *series[T]) NRows(options ...Options) int {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.RLock()
		defer s.RUnlock()
	}

	return len(s.values)
}

func (s *series[T]) Values() []T {
	return s.values
}

// Value returns the value of a particular row.
func (s *series[T]) Value(row int, options ...Options) T {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.RLock()
		defer s.RUnlock()
	}

	if row < 0 {
		return s.values[len(s.values)+row]
	}

	return s.values[row]
}

// ValueString returns a string representation of a
// particular row. The string representation is defined
// by the function set in SetValueToStringFormatter.
func (s *series[T]) ValueString(row int, options ...Options) string {
	return s.valFormatter(s.Value(row, options...))
}

// Prepend is used to set a value to the beginning of the
// series.
func (s *series[T]) Prepend(val []T, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	// See: https://stackoverflow.com/questions/41914386/what-is-the-mechanism-of-using-append-to-prepend-in-go

	if cap(s.values) > len(s.values)+len(val) {
		// There is already extra capacity so copy current values by 1 spot
		s.values = s.values[:len(s.values)+len(val)]
		copy(s.values[len(val):], s.values)
		copy(s.values, val)
		return
	}

	// No room, new slice needs to be allocated:
	s.insert(0, val)
}

// Append is used to set a value to the end of the series.
func (s *series[T]) Append(val []T, options ...Options) int {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	row := len(s.values)
	s.insert(row, val)
	return row
}

// Insert is used to set a value at an arbitrary row in
// the series. All existing values from that row onwards
// are shifted by 1.
func (s *series[T]) Insert(row int, val []T, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	s.insert(row, val)
}

func (s *series[T]) insert(row int, val []T) {
	s.values = append(s.values[:row], append(val, s.values[row:]...)...)
}

// Remove is used to delete the value of a particular row.
func (s *series[T]) Remove(row int, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	s.values = append(s.values[:row], s.values[row+1:]...)
}

// Reset is used clear all data contained in the series.
func (s *series[T]) Reset(options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	s.values = []T{}
}

// Update is used to update the value of a particular row.
func (s *series[T]) Update(row int, val T, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	if row < 0 {
		row = len(s.values) + row
	}

	s.values[row] = val
}

// valuesIterator will return a function that can be used to iterate through all the values.
func (s *series[T]) valuesIterator(options ...IteratorOptions) IteratorFn[T] {
	opts := DefaultOptions(options...)

	var row, step = opts.InitialRow, 1

	if row < 0 {
		row = len(s.values) + row
	}

	if opts.Step != 0 {
		step = opts.Step
	}

	initial := row

	return func() (int, T, int, bool) {
		if !opts.DontLock {
			s.RLock()
			defer s.RUnlock()
		}

		var t int
		if step > 0 {
			t = (len(s.values)-initial-1)/step + 1
		} else {
			t = -initial/step + 1
		}

		if row > len(s.values)-1 || row < 0 {
			// Don't iterate further
			return -1, *new(T), t, false
		}

		out := s.values[row]

		row = row + step
		return row - step, out, t, true
	}
}

// Iterator will return a iterator that can be used to iterate through all the values.
func (s *series[T]) Iterator(options ...IteratorOptions) Iterator[T] {
	return NewIterator(s.valuesIterator(options...))
}

// SetValueToStringFormatter is used to set a function
// to convert the value of a particular row to a string
// representation.
func (s *series[T]) SetValueToStringFormatter(f ValueToStringFormatter) {
	if f == nil {
		s.valFormatter = DefaultValueFormatter
		return
	}
	s.valFormatter = f
}

// Swap is used to swap 2 values based on their row position.
func (s *series[T]) Swap(row1, row2 int, options ...Options) {
	if row1 == row2 {
		return
	}

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	s.values[row1], s.values[row2] = s.values[row2], s.values[row1]
}

// IsEqualFunc returns true if a is equal to b.
func (s *series[T]) IsEqualFunc(a, b T) bool {
	if s.isEqualFunc == nil {
		panic(errors.New("IsEqualFunc not set"))
	}

	return s.isEqualFunc(a, b)
}

// IsLessThanFunc returns true if a is less than b.
func (s *series[T]) IsLessThanFunc(a, b T) bool {

	if s.isLessThanFunc == nil {
		panic(errors.New("IsLessThanFunc not set"))
	}

	return s.isLessThanFunc(a, b)
}

// SetIsEqualFunc sets a function which can be used to determine
// if 2 values in the series are equal.
func (s *series[T]) SetIsEqualFunc(f CompareFn[T]) {
	if f == nil {
		// Return to default
		s.isEqualFunc = IsEqualDefaultFunc[T]
	} else {
		s.isEqualFunc = f
	}
}

// SetIsLessThanFunc sets a function which can be used to determine
// if a value is less than another in the series.
func (s *series[T]) SetIsLessThanFunc(f CompareFn[T]) {
	s.isLessThanFunc = f
}

// Sort will sort the series.
// It will return true if sorting was completed or false when the context is canceled.
func (s *series[T]) Sort(ctx context.Context, options ...SortOptions) (completed bool) {

	if s.isLessThanFunc == nil {
		panic(errors.New("cannot sort without setting IsLessThanFunc"))
	}

	defer func() {
		if x := recover(); x != nil {
			completed = false
		}
	}()

	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.Lock()
		defer s.Unlock()
	}

	sortFunc := func(i, j int) (ret bool) {
		if err := ctx.Err(); err != nil {
			panic(err)
		}

		if opts.Desc {
			return !s.isLessThanFunc(s.values[i], s.values[j])
		}

		return s.isLessThanFunc(s.values[i], s.values[j])
	}

	if opts.Stable {
		sort.SliceStable(s.values, sortFunc)
	} else {
		sort.Slice(s.values, sortFunc)
	}

	return true
}

// Copy will create a new copy of the series.
// It is recommended that you lock the series before attempting
// to Copy.
func (s *series[T]) Copy(options ...RangeOptions) Series[T] {

	if len(s.values) == 0 {
		return &series[T]{
			valFormatter:   s.valFormatter,
			isEqualFunc:    s.isEqualFunc,
			isLessThanFunc: s.isLessThanFunc,
			name:           s.name,
			typeT:          s.typeT,
			values:         []T{},
		}
	}

	opts := DefaultOptions(options...)

	start, end, err := opts.Limits(len(s.values))
	if err != nil {
		panic(err)
	}

	// Copy slice
	x := s.values[start : end+1]
	newSlice := append(x[:0:0], x...)

	return &series[T]{
		valFormatter:   s.valFormatter,
		isEqualFunc:    s.isEqualFunc,
		isLessThanFunc: s.isLessThanFunc,
		name:           s.name,
		typeT:          s.typeT,
		values:         newSlice,
	}
}

// Table will produce the series in a table.
func (s *series[T]) Table(options ...TableOptions) string {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.RLock()
		defer s.RUnlock()
	}

	data := [][]string{}

	headers := []string{"", s.name} // row header is blank
	footers := []string{fmt.Sprintf("%dx%d", len(s.values), 1), s.Type()}

	if len(s.values) > 0 {
		start, end, err := opts.Range.Limits(len(s.values))
		if err != nil {
			panic(err)
		}

		for row := start; row <= end; row++ {
			sVals := []string{fmt.Sprintf("%d:", row), s.ValueString(row)}
			data = append(data, sVals)
		}
	}

	var buf bytes.Buffer

	table := tablewriter.NewWriter(&buf)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.SetFooter(footers)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Render()

	return buf.String()
}

// String implements the fmt.Stringer interface. It does not lock the series.
func (s *series[T]) String() string {

	count := len(s.values)

	out := s.name + ": [ "

	if count > 6 {
		idx := []int{0, 1, 2, count - 3, count - 2, count - 1}
		for j, row := range idx {
			if j == 3 {
				out = out + "... "
			}
			out = out + s.valFormatter(s.values[row]) + " "
		}
		return out + "]"
	}

	for row := range s.values {
		out = out + s.valFormatter(s.values[row]) + " "
	}
	return out + "]"

}

// FillRand will fill a series with random data.
func (s *series[T]) FillRand(rnd RandFn[T]) {

	for i := 0; i < len(s.values); i++ {
		s.values[i] = rnd()
	}

	capacity := cap(s.values)
	length := len(s.values)

	for i := 0; i < length; i++ {
		s.values[i] = rnd()
	}

	if capacity > length {
		excess := capacity - length
		for i := 0; i < excess; i++ {
			s.values = append(s.values, rnd())
		}
	}
}

// IsEqual returns true if s2's values are equal to s.
func (s *series[T]) IsEqual(ctx context.Context, s2 Series[T], options ...IsEqualOptions) (bool, error) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		s.RLock()
		defer s.RUnlock()
	}

	// Check number of values
	if len(s.values) != len(s2.Values()) {
		return false, nil
	}

	// Check name
	if opts.CheckName && s.name != s2.Name(Options{DontLock: true}) {
		return false, nil
	}

	// Check values
	for i, v := range s.values {
		if err := ctx.Err(); err != nil {
			return false, err
		}

		if !s.isEqualFunc(v, s2.Values()[i]) {
			return false, nil
		}
	}

	return true, nil
}

func (s *series[T]) Lock() {
	s.mutex.Lock()
}

func (s *series[T]) Unlock() {
	s.mutex.Unlock()
}

func (s *series[T]) RLock() {
	s.mutex.RLock()
}

func (s *series[T]) RUnlock() {
	s.mutex.RUnlock()
}

var _ Series[any] = (*series[any])(nil)
