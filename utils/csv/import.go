package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"io"

	"github.com/tradeoforigin/dataframe-go"
)

// LoadOptions holds the options for customizing CSV loading behavior.
type LoadOptions struct {
	// Comma is the field delimiter. The default value is ',' when CSVLoadOption is not provided.
	// Comma must be a valid rune and must not be \r, \n, or the Unicode replacement character (0xFFFD).
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the Comment character are ignored,
	// unless there is leading whitespace in which case the Comment character becomes part of the field.
	// Comment must be a valid rune and must not be \r, \n, or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune

	// If TrimLeadingSpace is true, leading whitespace in a field is ignored.
	// This is done even if the field delimiter (Comma) is whitespace.
	TrimLeadingSpace bool

	// Headers must be set if the CSV file does not contain a header row. 
	// This must be nil if the CSV file contains a header row.
	Headers []string
}

// Load reads a CSV file, parses it, and loads the data into a dataframe.
// The CSV data is read from an io.ReadSeeker and the series are converted using the specified converters.
// If LoadOptions' Headers field is not set, the first line in the CSV file is used as the header. If not, an error is returned.
//
// Example usage:
//
//	content, err := ioutil.ReadFile("data/data+header.csv")
//	if err != nil {
//		panic(err)
//	}
//
//	reader := strings.NewReader(string(content))
//
//	df, err := csv.Load(ctx, reader, map[string]csv.ConverterAny{
//		"time": csv.Time, "o": csv.Float64, "h": csv.Float64, "l": csv.Float64, "c": csv.Float64, "v": csv.Float64,
//	})
//
//	if err != nil {
//		panic(err)
//	}
func Load(ctx context.Context, r io.ReadSeeker, converters map[string]ConverterAny, options ...LoadOptions) (*dataframe.DataFrame, error) {
	// Apply the default options for loading the CSV if not provided.
	opts := dataframe.DefaultOptions(options...)

	// Default to comma (',') if no delimiter is set.
	if opts.Comma == 0 {
		opts.Comma = ','
	}

	// Initialize the dataframe series.
	init := dataframe.SeriesInit{}

	// Create a CSV reader with the specified delimiter and other options.
	cr := csv.NewReader(r)
	cr.Comma = opts.Comma
	cr.Comment = opts.Comment
	cr.TrimLeadingSpace = opts.TrimLeadingSpace

	// Count the rows in the CSV file.
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		_, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				// Rewind to the beginning of the file after counting rows.
				r.Seek(0, io.SeekStart)
				break
			}
			return nil, err
		}
		init.Capacity++
	}

	// If headers are not provided in options, use the first row of the CSV as headers.
	if len(opts.Headers) == 0 {
		headers, err := cr.Read()
		if err != nil {
			return nil, err
		}
		opts.Headers = headers
		init.Capacity-- // Exclude header row from capacity.
	}

	// Initialize series and map the series name to the column index.
	series := make([]dataframe.SeriesAny, 0, len(converters))
	seriesIdx := map[string]int{}

	for name, converter := range converters {
		series = append(series, converter.series(name, &init))

		for i := range opts.Headers {
			if name == opts.Headers[i] {
				seriesIdx[name] = i
				break
			}
		}
	}

	// Ensure that all series are mapped correctly to columns in the CSV.
	if len(seriesIdx) != len(series) {
		return nil, errors.New("could not map columns to series")
	}

	// Create a new dataframe with the initialized series.
	df := dataframe.NewDataFrame(series...)

	// Read and parse the CSV data line by line.
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		record, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Map the record values to the corresponding series in the dataframe.
		row := map[string]any{}
		for name, idx := range seriesIdx {
			row[name] = converters[name].value(record[idx])
		}

		// Append the row to the dataframe.
		df.Append(row)
	}

	return df, nil
}
