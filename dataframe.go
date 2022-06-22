package dataframe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/sync/errgroup"
)

// DataFrame allows you to handle numerous
//series of data conveniently.
type DataFrame struct {
	Series []SeriesAny
	n      int // Number of rows

	lock sync.RWMutex
}

func NewDataFrame(se ...SeriesAny) *DataFrame {
	df := &DataFrame{
		Series: []SeriesAny{},
	}
	
	if len(se) > 0 {
		var count int = -1
		names := map[string]bool{}

		for _, s := range se {
			if count == -1 {
				count = s.NRows()
				names[s.Name()] = true
			} else {
				if count != s.NRows() {
					panic("different number of rows in series")
				}
				if names[s.Name()] {
					panic("names of series must be unique: " + s.Name())
				}
				names[s.Name()] = true
			}
			df.Series = append(df.Series, s)
		}

		df.n = count
	}

	return df
}

// NRows returns the number of rows of data.
// Each series must contain the same number of rows.
func (df *DataFrame) NRows(options ...Options) int {
	opts := DefaultOptions(options...)
	
	if !opts.DontLock {
		df.lock.RLock(); defer df.lock.RUnlock()
	}

	return df.n
}

// Row returns the series' values for a particular row.
//
// Example:
//
//  df.Row(5, false)
//
func (df *DataFrame) Row(row int, options ...Options) map[string]any {
	opts := DefaultOptions(options...)
	
	if !opts.DontLock {
		df.lock.RLock(); defer df.lock.RUnlock()
	}

	out := map[string]any{}

	for _, aSeries := range df.Series {
		out[aSeries.Name()] = aSeries.ValueAny(row)
	}

	return out
}

// ValuesIterator will return a function that can be used to iterate through all the values.
func (df *DataFrame) valuesIterator(options ...IteratorOptions) IteratorFn[map[string]any] {
	opts := DefaultOptions(options...)

	var row, step = opts.InitialRow, 1

	if row < 0 {
		row = df.n + row
	}

	if opts.Step != 0 {
		step = opts.Step
	}

	initial := row

	return func() (int, map[string]any, int, bool) {
		if !opts.DontLock {
			df.lock.RLock()
			defer df.lock.RUnlock()
		}

		var t int
		if step > 0 {
			t = (df.n - initial - 1) / step + 1
		} else {
			t = -initial / step + 1
		}

		if row > df.n - 1 || row < 0 {
			// Don't iterate further
			return -1, nil, t, false
		}

		out := map[string]any{}

		for _, aSeries := range df.Series {
			out[aSeries.Name()] = aSeries.ValueAny(row)
		}

		row = row + step

		return row - step, out, t, true
	}
}

func (s *DataFrame) Iterator(options ...IteratorOptions) Iterator[map[string]any] {
	return NewIterator(s.valuesIterator(options...))
}

// Prepend inserts a row at the beginning.
func (df *DataFrame) Prepend(vals any, options ...Options) {
	df.Insert(0, vals, options...)
}

// Append inserts a row at the end.
func (df *DataFrame) Append(vals any, options ...Options) {
	df.Insert(df.n, vals, options...)
}

// Insert adds a row to a particular position.
func (df *DataFrame) Insert(row int, vals any, options ...Options) {
	opts := DefaultOptions(options...)
	
	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	df.insert(row, vals)
}

func (df *DataFrame) insert(row int, vals any) {

	var nRows = df.n

	switch v := vals.(type) {
	case map[string]any:

		// Check if number of vals is equal to number of series
		if len(v) != len(df.Series) {
			panic("no. of args not equal to no. of series")
		}

		var idx int

		for name, val := range v {
			col := df.MustNameToColumn(name, dontLock)
			df.Series[col].InsertAny(row, val)

			sRows := df.Series[col].NRows(dontLock)
			if idx != 0 && nRows != sRows {
				panic("series length does not match")
			} 
			
			nRows = sRows
			idx += 1
		}

	case map[int]any:

		// Check if number of vals is equal to number of series
		if len(v) != len(df.Series) {
			panic("no. of args not equal to no. of series")
		}

		for idx, s := range df.Series {
			s.InsertAny(row, v[idx])

			sRows := s.NRows(dontLock)
			if idx != 0 && nRows != sRows {
				panic("series length does not match")
			} else {
				nRows = sRows
			}
		}

	case map[any]any:

		// Check if number of vals is equal to number of series
		names := map[string]bool{}

		for key := range v {
			switch kTyp := key.(type) {
			case int:
				names[df.Series[kTyp].Name(dontLock)] = true
			case string:
				names[kTyp] = true
			default:
				panic("unknown type in insert argument. Must be an int or string.")
			}
		}

		if len(names) != len(df.Series) {
			panic("no. of args not equal to no. of series")
		}

		var idx int

		for C, val := range v {
			var col = 0

			switch CTyp := C.(type) {
			case int: 
				col = CTyp
			case string:
				col = df.MustNameToColumn(CTyp, dontLock)
			default:
				panic("unknown type in insert argument. Must be an int or string.")
			}

			df.Series[col].InsertAny(row, val)

			sRows := df.Series[col].NRows(dontLock)
			if idx != 0 && nRows != sRows {
				panic("series length does not match")
			} 
			
			nRows = sRows
			idx += 1
		}

	case []any:
		// Check if number of vals is equal to number of series
		if len(v) != len(df.Series) {
			panic("no. of args not equal to no. of series")
		}

		for idx, val := range v {
			df.Series[idx].InsertAny(row, val)

			sRows := df.Series[idx].NRows(dontLock)
			if idx != 0 && nRows != sRows {
				panic("series length does not match")
			} else  {
				nRows = sRows
			}
		}

	default:
		panic("invalid type to insert")
	}

	df.n = nRows
}

// Remove deletes a row.
func (df *DataFrame) Remove(row int, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	for i := range df.Series {
		df.Series[i].Remove(row)
	}
	df.n--
}

// Update is used to update a specific entry.
// col can be the name of the series or the column number.
func (df *DataFrame) Update(row int, col any, val any, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	switch name := col.(type) {
	case string:
		col = df.MustNameToColumn(name, dontLock)
	}

	df.Series[col.(int)].UpdateAny(row, val)
}

// UpdateRow will update an entire row.
func (df *DataFrame) UpdateRow(row int, vals any, options ...Options) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	switch v := vals.(type) {
	case map[string]any:
		for name, val := range v {
			df.Series[df.MustNameToColumn(name, dontLock)].UpdateAny(row, val)
		}
	case map[int]any:
		for idx, val := range v {
			df.Series[idx].UpdateAny(row, val)
		}
	case map[any]any:
		for C, val := range v {
			switch CTyp := C.(type) {
			case int:
				df.Series[CTyp].UpdateAny(row, val)
			case string:
				df.Series[df.MustNameToColumn(CTyp, dontLock)].UpdateAny(row, val)
			default:
				panic("unknown type in UpdateRow argument. Must be an int or string.")
			}
		}
	case []any:
		// Check if number of vals is equal to number of series
		if len(v) != len(df.Series) {
			panic("no. of args not equal to no. of series")
		}

		for idx, val := range v {
			df.Series[idx].UpdateAny(row, val)
		}

	default:
		panic("invalid type to update")
	}
}

// Names will return a list of all the series names.
func (df *DataFrame) Names(options ...Options) []string {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.RLock()
		defer df.lock.RUnlock()
	}

	names := make([]string, 0, len(df.Series))
	for _, aSeries := range df.Series {
		names = append(names, aSeries.Name(options...))
	}

	return names
}

// MustNameToColumn returns the index of the series based on the name.
// The starting index is 0. If seriesName doesn't exist it panics.
func (df *DataFrame) MustNameToColumn(seriesName string, options ...Options) int {
	col, err := df.NameToColumn(seriesName, options...)
	if err != nil {
		panic(err)
	}

	return col
}

// NameToColumn returns the index of the series based on the name.
// The starting index is 0.
func (df *DataFrame) NameToColumn(seriesName string, options ...Options) (int, error) {
	opts := DefaultOptions(options...)
	
	if !opts.DontLock {
		df.lock.RLock()
		defer df.lock.RUnlock()
	}

	for idx, aSeries := range df.Series {
		if aSeries.Name() == seriesName {
			return idx, nil
		}
	}

	return 0, errors.New("no series contains name")
}

// ReorderColumns reorders the columns based on an ordered list of
// column names. The length of newOrder must match the number of columns
// in the Dataframe. The column names in newOrder must be unique.
func (df *DataFrame) ReorderColumns(newOrder []string, options ...Options) error {
	opts := DefaultOptions(options...)
	
	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	if len(newOrder) != len(df.Series) {
		return errors.New("length of newOrder must match number of columns")
	}

	// Check if newOrder contains duplicates
	fields := map[string]bool{}
	for _, v := range newOrder {
		fields[v] = true
	}

	if len(fields) != len(df.Series) {
		return errors.New("newOrder must not contain duplicate values")
	}

	series := []SeriesAny{}

	for _, v := range newOrder {
		idx, err := df.NameToColumn(v, dontLock)
		if err != nil {
			return errors.New(err.Error() + ": " + v)
		}

		series = append(series, df.Series[idx])
	}

	df.Series = series

	return nil
}

// RemoveSeries will remove a Series from the Dataframe.
func (df *DataFrame) RemoveSeries(seriesName string, options ...Options) error {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	idx, err := df.NameToColumn(seriesName, dontLock)
	if err != nil {
		return errors.New(err.Error() + ": " + seriesName)
	}

	df.Series = append(df.Series[:idx], df.Series[idx+1:]...)
	return nil
}

// AddSeries will add a Series to the end of the DataFrame, unless set by ColN.
func (df *DataFrame) AddSeries(s SeriesAny, colN *int, options ...Options) error {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	if s.NRows(dontLock) != df.n {
		panic("different number of rows in series")
	}

	if colN == nil {
		df.Series = append(df.Series, s)
	} else {
		df.Series = append(df.Series, nil)
		copy(df.Series[*colN+1:], df.Series[*colN:])
		df.Series[*colN] = s
	}

	return nil
}

// Swap is used to swap 2 values based on their row position.
func (df *DataFrame) Swap(row1, row2 int, options ...Options) {
	opts := DefaultOptions(options...)
	
	if !opts.DontLock {
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	for idx := range df.Series {
		df.Series[idx].Swap(row1, row2)
	}
}

// Lock will lock the Dataframe allowing you to directly manipulate
// the underlying Series with confidence.
func (df *DataFrame) Lock(deep ...bool) {
	df.lock.Lock()

	if len(deep) > 0 && deep[0] {
		for i := range df.Series {
			df.Series[i].Lock()
		}
	}
}

// Unlock will unlock the Dataframe that was previously locked.
func (df *DataFrame) Unlock(deep ...bool) {
	if len(deep) > 0 && deep[0] {
		for i := range df.Series {
			df.Series[i].Unlock()
		}
	}

	df.lock.Unlock()
}


// Lock will lock the Dataframe allowing you to directly manipulate
// the underlying Series with confidence.
func (df *DataFrame) RLock(deep ...bool) {
	df.lock.RLock()

	if len(deep) > 0 && deep[0] {
		for i := range df.Series {
			df.Series[i].RLock()
		}
	}
}

// Unlock will unlock the Dataframe that was previously locked.
func (df *DataFrame) RUnlock(deep ...bool) {
	if len(deep) > 0 && deep[0] {
		for i := range df.Series {
			df.Series[i].RUnlock()
		}
	}

	df.lock.RUnlock()
}

// Copy will create a new copy of the Dataframe.
// It is recommended that you lock the Dataframe
// before attempting to Copy.
func (df *DataFrame) Copy(options ...RangeOptions) *DataFrame {

	series := []SeriesAny{}
	for i := range df.Series {
		series = append(series, df.Series[i].CopyAny(options...))
	}

	newDF := &DataFrame{
		Series: series,
	}

	if len(series) > 0 {
		newDF.n = series[0].NRows(dontLock)
	}

	return newDF
}

// Table will produce the DataFrame in a table.
func (df *DataFrame) Table(options ...TableOptions) string {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.RLock(); defer df.lock.RUnlock()
	}

	columns := map[any]bool{}
	for _, v := range opts.Series {
		columns[v] = true
	}

	data := [][]string{}

	headers := []string{""} // row header is blank
	footers := []string{fmt.Sprintf("%dx%d", df.n, len(df.Series))}
	for idx, aSeries := range df.Series {
		if len(columns) == 0 {
			headers = append(headers, aSeries.Name())
			footers = append(footers, aSeries.Type())
		} else {
			// Check idx
			if columns[idx] {
				headers = append(headers, aSeries.Name())
				footers = append(footers, aSeries.Type())
				continue
			}

			// Check series name
			if columns[aSeries.Name()] {
				headers = append(headers, aSeries.Name())
				footers = append(footers, aSeries.Type())
				continue
			}
		}
	}

	if df.n > 0 {
		start, end, err := opts.Range.Limits(df.n)
		if err != nil {
			panic(err)
		}

		for row := start; row <= end; row++ {

			sVals := []string{ fmt.Sprintf("%d:", row )}

			for idx, aSeries := range df.Series {
				if len(columns) == 0 {
					sVals = append(sVals, aSeries.ValueString(row))
				} else {
					// Check idx
					if columns[idx] {
						sVals = append(sVals, aSeries.ValueString(row))
						continue
					}

					// Check series name
					if columns[aSeries.Name()] {
						sVals = append(sVals, aSeries.ValueString(row))
						continue
					}
				}
			}

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

// String implements the fmt.Stringer interface. It does not lock the DataFrame.
func (df *DataFrame) String() string {

	if df.NRows() <= 6 {
		return df.Table(TableOptions{ DontLock: true })
	}

	idx := []int{0, 1, 2, df.n - 3, df.n - 2, df.n - 1}

	data := [][]string{}

	headers := []string{""} // row header is blank
	footers := []string{fmt.Sprintf("%dx%d", df.n, len(df.Series))}
	for _, aSeries := range df.Series {
		headers = append(headers, aSeries.Name())
		footers = append(footers, aSeries.Type())
	}

	for j, row := range idx {

		if j == 3 {
			sVals := []string{"⋮"}

			for range df.Series {
				sVals = append(sVals, "⋮")
			}

			data = append(data, sVals)
		}

		sVals := []string{fmt.Sprintf("%d:", row)}

		for _, aSeries := range df.Series {
			sVals = append(sVals, aSeries.ValueString(row))
		}

		data = append(data, sVals)
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


// FillRand will randomly fill all the Series in the Dataframe.
func (df *DataFrame) FillRand(rnd RandFn[any]) {
	for _, s := range df.Series {
		// should be changed
		s.FillRandAny(rnd)
	}
}

var errNotEqual = errors.New("not equal")

// IsEqual returns true if df2's values are equal to df.
func (df *DataFrame) IsEqual(ctx context.Context, df2 *DataFrame, options ...IsEqualOptions) (bool, error) {
	opts := DefaultOptions(options...)

	if !opts.DontLock {
		df.lock.RLock()
		defer df.lock.RUnlock()
	}

	// Check if number of columns are the same
	if len(df.Series) != len(df2.Series) {
		return false, nil
	}

	// Check values
	g, newCtx := errgroup.WithContext(ctx)

	for i := range df.Series {
		i := i
		g.Go(func() error {

			eq, err := df.Series[i].IsEqualAny(newCtx, df2.Series[i], options...)
			if err != nil {
				return err
			}

			if !eq {
				return errNotEqual
			}

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		if err == errNotEqual {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
