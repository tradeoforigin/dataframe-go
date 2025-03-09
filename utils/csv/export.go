package csv

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/tradeoforigin/dataframe-go"
)

// ExportOptions contains configuration options for the Export function, 
// which exports a dataframe to CSV format. These options control various 
// aspects of the export process, such as delimiters, line terminators, and 
// handling of missing values.

type ExportOptions struct {
	// NullString specifies the string representation for nil values in the CSV.
	// For example, "NULL", "\N", "NaN", "NA" are common representations for nulls.
	// If not provided, the default "nil" will be used.
	NullString *string

	// Range defines the subset of rows to be exported from the dataframe.
	// This allows exporting a specific range of rows instead of the entire dataframe.
	Range dataframe.RangeOptions

	// Comma is the field delimiter for the CSV. The default delimiter is ','.
	// If no value is provided, the default comma delimiter will be used.
	Comma rune

	// UseCRLF determines the line terminator format. If true, the line terminator
	// will be set to \r\n (CRLF). By default, the function uses \n (LF).
	UseCRLF bool
}

// Export generates a CSV representation of a dataframe and writes it to
// the provided writer (e.g., file or network connection). The function
// handles options like delimiters, null value encoding, and row ranges 
// during the export. It writes the header row and then the data rows to 
// the CSV format, one row at a time.
//
// Example usage:
//
//	ctx := context.Background()
//
//	// Create some sample series
//	s1 := dataframe.NewSeries("str", nil, "one", "one,two", "one,two,three")
//	s2 := dataframe.NewSeries("num", nil, 1, 12, 123)
//
//	// Create a dataframe with the series
//	df1 := dataframe.NewDataFrame(s1, s2)
//
//	// Open a file for writing
//	f, err := os.OpenFile("data/export.csv", os.O_WRONLY|os.O_CREATE, 0600)
//	if err != nil {
//		panic(err)
//	}
//
//	// Export the dataframe to the CSV file
//	err = csv.Export(ctx, f, df1)
//	if err != nil {
//		panic(err)
//	}
//
//	// Close the file after writing
//	f.Close()
func Export(ctx context.Context, w io.Writer, df dataframe.DataFrame, options ...ExportOptions) error {
	opts := dataframe.DefaultOptions(options...)

	// Create a new CSV writer
	cw := csv.NewWriter(w)

	// Default null string representation for nil values
	nullString := "nil"

	// Set the field delimiter (comma by default)
	if opts.Comma == 0 {
		cw.Comma = ','
	}

	// Set the null value string representation if provided
	if opts.NullString != nil {
		nullString = *opts.NullString
	}

	// Set the line terminator (CRLF if true, LF if false)
	cw.UseCRLF = opts.UseCRLF

	// Lock the dataframe for safe concurrent access
	df.Lock()
	defer df.Unlock()

	// Write the header row (series names)
	if err := cw.Write(df.Names(dataframe.DontLock)); err != nil {
		return err
	}

	// Get the number of rows in the dataframe
	nRows := df.NRows(dataframe.DontLock)

	if nRows > 0 {
		// Get the range of rows to export
		start, end, err := opts.Range.Limits(nRows)
		if err != nil {
			return err
		}

		// Get the column names (keys)
		var keys []string = df.Names(dataframe.DontLock)

		// Create an iterator for the dataframe
		iterator := df.Iterator(dataframe.IteratorOptions{
			InitialRow: start,
			Step:       1,
			DontLock:   true,
		})

		// Iterate through the rows and write them to the CSV
		for iterator.Next() && iterator.Index <= end {
			// Flush every 100 rows to avoid memory overload
			if (iterator.Index-start+1)%100 == 0 {
				cw.Flush()
				if err := cw.Error(); err != nil {
					return err
				}
			}

			// Prepare a slice to store the row values as strings
			sVals := make([]string, 0, len(df.Series()))
			for i, key := range keys {
				val := iterator.Value[key]
				if val == nil {
					// If value is nil, append the null string representation
					sVals = append(sVals, nullString)
				} else {
					// Otherwise, convert the value to a string
					sVals = append(sVals, df.Series()[i].ValueString(iterator.Index, dataframe.DontLock))
				}
			}

			// Write the row to the CSV file
			if err := cw.Write(sVals); err != nil {
				return err
			}
		}
	}

	// Flush remaining data before exiting
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}

	return nil
}
