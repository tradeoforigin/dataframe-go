package dataframe

import (
	"context"
)

type SeriesAny interface {
	
	// Name returns the series name.
	Name(options ...Options) string

	// Rename renames the series.
	Rename(n string, options ...Options) 

	// Type returns type of the series as string value.
	Type() string

	// NRows returns how many rows the series contains.
	NRows(options ...Options) int

	// ValueAny returns the value of a particular row.
	ValueAny(row int, options ...Options) any

	// ValueString returns a string representation of a
	// particular row. The string representation is defined
	// by the function set in SetValueToStringFormatter.
	// By default, a nil value is returned as "NaN".
	ValueString(row int, options ...Options) string

	// Prepend is used to set a value to the beginning of the
	// series.
	PrependAny(val any, options ...Options)

	// AppendAny is used to set a value to the end of the series.
	AppendAny(val any, options ...Options) int

	// InsertAny is used to set a value at an arbitrary row in
	// the series. All existing values from that row onwards
	// are shifted by 1.
	InsertAny(row int, val any, options ...Options)

	// Remove is used to delete the value of a particular row.
	Remove(row int, options ...Options)

	// Reset is used clear all data contained in the Series.
	Reset(options ...Options)

	// Update is used to update the value of a particular row.
	UpdateAny(row int, val any, options ...Options)

	// IteratorAny will return a iterator that can be used to iterate through all the values.
	IteratorAny(options ...IteratorOptions) Iterator[any]

	// SetValueToStringFormatter is used to set a function
	// to convert the value of a particular row to a string
	// representation.
	SetValueToStringFormatter(f ValueToStringFormatter)

	// Swap is used to swap 2 values based on their row position.
	Swap(row1, row2 int, options ...Options)

	// IsEqualAnyFunc	returns true if a is equal to b.
	IsEqualAnyFunc(a, b any) bool

	// IsLessThanAnyFunc	returns true if a is less than b.
	IsLessThanAnyFunc(a, b any) bool

	// SetIsEqualAnyFunc	sets a function which can be used to determine
	// if 2 values in the series are equal.
	SetIsEqualAnyFunc(f CompareFn[any])

	// SetIsLessThanAnyFunc	sets a function which can be used to determine
	// if a value is less than another in the series.
	SetIsLessThanAnyFunc(f CompareFn[any])

	// Sort will sort the series.
	// It will return true if sorting was completed or false when the context is canceled.
	Sort(ctx context.Context, options ...SortOptions) (completed bool)

	// CopyAny will create a new copy of the series.
	// It is recommended that you lock the Series before attempting
	// to Copy.
	CopyAny(options ...RangeOptions) SeriesAny

	// Table will produce the Series in a table.
	Table(options ...TableOptions) string

	// String implements the fmt.Stringer interface. It does not lock the Series.
	String() string

	// FillRandAny will fill a Series with random data.
	FillRandAny(rnd RandFn[any])

	// IsEqualAny returns true if s2's values are equal to s.
	IsEqualAny(ctx context.Context, s2 SeriesAny, options ...IsEqualOptions) (bool, error)

	// RWMutex Lock
	Lock()

	// RWMutex Unlock 
	Unlock()

	// RWMutex RLock
	RLock()

	// RWMutex RUnlock
	RUnlock()

	// Creates clone with empty Values
	cloneAsEmpty(size ...int) SeriesAny
}