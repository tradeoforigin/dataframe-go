package dataframe

// Options is used to perform operation with DontLock.
// Notice that all operations on the series or
// dataframes are performed with locked RWMutex.
//
// Defaults:
// 		Options { DontLock: false }
//
// Properties:
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type Options struct {
	DontLock bool
}

// IsEqualOptions is defined as an optional parameters
// for IsEqual(...) on top of Series or DataFrame.
//
// Defaults:
// 		IsEqualOptions { CheckName: false, DontLock: false }
//
// Properties:
//	• `CheckName` - indicates that name should be checked in form of equality
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type IsEqualOptions struct {
	CheckName, DontLock bool
}

// FilterOptions is defined as an optional parameters
// for Filter(...) on top of Series or DataFrame.
//
// Defaults:
// 		FilterOptions { InPlace: false, DontLock: false }
//
// Properties:
//	• `InPlace` - Filter affects current Series/DataFrame and no new one is returned
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type FilterOptions struct {
	InPlace, DontLock bool
}

// ApplyOptions is defined as an optional parameters
// for Apply(...) on top of Series or DataFrame.
//
// Defaults:
// 		ApplyOptions { InPlace: false, DontLock: false }
//
// Properties:
//	• `InPlace` - Apply affects current Series/DataFrame and no new one is returned
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type ApplyOptions = FilterOptions

// RangeOptions is defined as an optional parameters
// for functions which needs range like Copy(...), Apply(...),
// Filter(...), etc.
//
// Notice that DataFrame and Series calls Limits(length) on top of 
// RangeOptions passed. In case of `r.End == nil`, end is set to -1.
// Negative values provides indexing from the end. For example the
// Range(0, -1) is the same as Range(0, len(arr) - 1)
//
// Defaults:
// 		RangeOptions { Start: 0, End: nil }
//
// Properties:
//	• `Start` - Defines start row/index for iteration/copy
//	• `End` - Defines where iteration/copy should end
type RangeOptions struct {
	Start int 
	End   *int
}


// TableOptions is defined as an optional parameters
// for Table(...) on top of Series or DataFrame.
//
// Defaults:
//		TableOptions { 
//			Series: nil, 
//			Range: RangeOptions { Start: 0, End: nil }
//			DontLock: false 
//		}
//
// Properties:
//	• `Series` - is int or string and indicates which series should table contains. Affets only DataFrame
//	• `Range` - specifies range for displayed table
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type TableOptions struct {
	Series []any
	Range RangeOptions
	DontLock bool
}

// SortOptions is defined as an optional parameters
// for Sort(...) on top of Series or DataFrame.
//
// Defaults:
//		SortOptions { 
//			Stable: false,
//			Desc: false,
//			DontLock: false 
//		}
//
// Properties:
//	• `Desc` - if true, then values will be sorted in descending order
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type SortOptions struct {
	Stable, Desc, DontLock bool
}

// SortOptions is defined as an optional parameters
// for Iterator(...) on top of Series or DataFrame.
//
// Defaults:
//		IteratorOptions { 
//			InitialRow: 0,
//			Step: 1,
//			DontLock: false 
//		}
//
// Properties:
//	• `InitialRow` - if set then iterator will start at this row.
//	• `Step` - iteration step. Negative values causes backward iterations
//	• `DontLock` - if set to true, then operation is performed without locking RWMutex 
type IteratorOptions struct {
	InitialRow, Step int
	DontLock bool
}

var dontLock = Options { DontLock: true }
// shortcut for Options { DontLock: true }
var DontLock = dontLock