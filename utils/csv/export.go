package csv

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/tradeoforigin/dataframe-go"
)

// CSVExportOptions contains options for ExportToCSV function.
type ExportOptions struct {

	// NullString is used to set what nil values should be encoded to.
	// Common options are NULL, \N, NaN, NA.
	NullString *string

	// Range is used to export a subset of rows from the dataframe.
	Range dataframe.RangeOptions

	// Separator is the field delimiter. A common option is ',', which is
	// the default if ExportOptions is not provided.
	Comma rune

	// UseCRLF determines the line terminator.
	// When true, it is set to \r\n.
	UseCRLF bool
}

func Export(ctx context.Context, w io.Writer, df *dataframe.DataFrame, options ...ExportOptions) error {
	opts := dataframe.DefaultOptions(options...)

	cw := csv.NewWriter(w)

	nullString := "nil"

	if opts.Comma == 0 {
		cw.Comma = ','
	}

	if opts.NullString != nil {
		nullString = *opts.NullString
	}

	cw.UseCRLF = opts.UseCRLF

	df.Lock(); defer df.Unlock()

	// Write header -> series names
	if err := cw.Write(df.Names(dataframe.DontLock)); err != nil {
		return err
	}

	nRows := df.NRows(dataframe.DontLock)

	if nRows > 0 {
		start, end, err := opts.Range.Limits(nRows)

		if err != nil {
			return err
		}

		for row := start; row <= end; row++ {
			// flush every 100 rows
			if (row - start + 1) % 100 == 0 {
				cw.Flush()
				if err := cw.Error(); err != nil {
					return err
				}
			}

			sVals := make([]string, 0, len(df.Series))
			for _, aSeries := range df.Series {
				val := aSeries.ValueAny(row)
				if val == nil {
					sVals = append(sVals, nullString)
				} else {
					sVals = append(sVals, aSeries.ValueString(row, dataframe.DontLock))
				}
			}

			// Write every row
			if err := cw.Write(sVals); err != nil {
				return err
			}
		}
	}

	// flush before exit
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}

	return nil
}