package dataframe

// Options is used to configure operations that can bypass locking the RWMutex. 
// By default, all operations on Series or DataFrames are performed with an 
// RWMutex lock for thread-safety. However, setting `DontLock` to true allows 
// operations to be performed without locking, which can be useful in certain
// scenarios where thread safety is not required, or when working with data 
// that is already protected by other means.
//
// Defaults:
//   Options { DontLock: false }
//
// Properties:
//   - `DontLock`: If set to true, the operation is performed without locking 
//     the RWMutex.
type Options struct {
	DontLock bool
}

// IsEqualOptions defines optional parameters for the IsEqual(...) function 
// that compares two Series or DataFrames. This struct is used to specify 
// whether the equality comparison should include the names of the series/columns 
// and whether the operation should bypass the RWMutex lock.
//
// Defaults:
//   IsEqualOptions { CheckName: false, DontLock: false }
//
// Properties:
//   - `CheckName`: If true, the names of the Series or DataFrame columns will
//     be checked for equality during the comparison.
//   - `DontLock`: If true, the operation will be performed without locking 
//     the RWMutex.
type IsEqualOptions struct {
	CheckName, DontLock bool
}

// FilterOptions defines optional parameters for the Filter(...) function 
// applied to Series or DataFrame. This struct allows you to specify whether 
// the filter should modify the data in-place or return a new filtered version.
//
// Defaults:
//   FilterOptions { InPlace: false, DontLock: false }
//
// Properties:
//   - `InPlace`: If true, the filter will modify the existing Series/DataFrame 
//     and no new one will be returned.
//   - `DontLock`: If true, the operation will be performed without locking 
//     the RWMutex.
type FilterOptions struct {
	InPlace, DontLock bool
}

// ApplyOptions is defined as an optional parameter for the Apply(...) function
// on Series or DataFrame. It works similarly to FilterOptions, specifying whether 
// the operation should be applied in-place or require a new object to be returned.
//
// Defaults:
//   ApplyOptions { InPlace: false, DontLock: false }
//
// Properties:
//   - `InPlace`: If true, the Apply operation will modify the current Series/DataFrame 
//     in place and no new one will be returned.
//   - `DontLock`: If true, the operation will be performed without locking the RWMutex.
type ApplyOptions = FilterOptions

// RangeOptions defines optional parameters for functions like Copy(...), Apply(...),
// Filter(...), etc., that require a range of rows or indices. This struct allows 
// you to specify a start and end index for iteration or copying data.
//
// Defaults:
//   RangeOptions { Start: 0, End: nil }
//
// Properties:
//   - `Start`: Defines the starting row/index for iteration or copy.
//   - `End`: Defines the end row/index for iteration or copy. If nil, it indicates 
//     the end of the available data, and negative values can be used to count from 
//     the end (e.g., -1 means the last element).
type RangeOptions struct {
	Start int
	End   *int
}

// TableOptions defines optional parameters for generating a table representation 
// of a Series or DataFrame using the Table(...) function. This struct allows 
// you to select which columns to include in the table and specify a range of rows.
//
// Defaults:
//   TableOptions { 
//     Series: nil, 
//     Range: RangeOptions { Start: 0, End: nil },
//     DontLock: false 
//   }
//
// Properties:
//   - `Series`: A list of column indices or names that specifies which columns 
//     should be included in the table. This only affects DataFrame.
//   - `Range`: Specifies the range of rows to include in the table.
//   - `DontLock`: If true, the operation will be performed without locking 
//     the RWMutex.
type TableOptions struct {
	Series []any
	Range  RangeOptions
	DontLock bool
}

// SortOptions defines optional parameters for sorting a Series or DataFrame 
// using the Sort(...) function. This struct allows you to specify if the sorting 
// should be stable (preserving the order of equal values) and if the sorting 
// should be in descending order.
//
// Defaults:
//   SortOptions { 
//     Stable: false,
//     Desc: false,
//     DontLock: false 
//   }
//
// Properties:
//   - `Stable`: If true, the sorting will be stable, meaning that elements 
//     that compare equal will remain in their original order.
//   - `Desc`: If true, the values will be sorted in descending order.
//   - `DontLock`: If true, the operation will be performed without locking 
//     the RWMutex.
type SortOptions struct {
	Stable, Desc, DontLock bool
}

// IteratorOptions defines optional parameters for the Iterator(...) function 
// on Series or DataFrame. This struct allows you to specify the starting row 
// and the step size for iteration (positive or negative values to iterate forward 
// or backward).
//
// Defaults:
//   IteratorOptions { 
//     InitialRow: 0,
//     Step: 1,
//     DontLock: false 
//   }
//
// Properties:
//   - `InitialRow`: If set, the iterator will start at this row. Defaults to 0.
//   - `Step`: Specifies the step size for iteration. A positive value will 
//     iterate forward, while a negative value will iterate backward.
//   - `DontLock`: If true, the operation will be performed without locking 
//     the RWMutex.
type IteratorOptions struct {
	InitialRow, Step int
	DontLock bool
}

// dontLock is a shortcut for Options { DontLock: true }, allowing you to 
// quickly disable locking operations.
var dontLock = Options{DontLock: true}

// DontLock is a pre-configured shortcut to disable locking during operations.
var DontLock = dontLock
