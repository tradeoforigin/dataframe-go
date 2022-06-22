package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"io"

	"github.com/tradeoforigin/dataframe-go"
)

type LoadOptions struct {
	// Comma is the field delimiter.
	// The default value is ',' when CSVLoadOption is not provided.
	// Comma must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	// Comment must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune

	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool

	// Headers must be set if the CSV file does not contain a header row. This must be nil if the CSV file contains a
	// header row.
	Headers []string
}

// Function to load CSV data into dataframe. CSV loader is defined by io.ReadSeaker and converters
// for specific series. Series defined in converters will be under the same name in resulted dataframe. If
// LoadOptions headers field is not set, the CSV file must contains header line at the first place, otherwise
// error is returned. 
// 
// Example:
//
//	content, err := ioutil.ReadFile("data/data+header.csv")
//	if err != nil {
//		panic(err)
//	}
//
//	reader := strings.NewReader(string(content))
//
//	df, err := csv.Load(ctx, reader, map[string]csv.ConverterAny {
//		"time": csv.Time, "o": csv.Float64, "h": csv.Float64, "l": csv.Float64, "c": csv.Float64, "v": csv.Float64,
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
func Load(ctx context.Context, r io.ReadSeeker, converters map[string]ConverterAny, options ...LoadOptions) (*dataframe.DataFrame, error) {
	opts := dataframe.DefaultOptions(options...)

	if opts.Comma == 0 {
		opts.Comma = ','
	}

	init := dataframe.SeriesInit{}

	cr := csv.NewReader(r)
	cr.Comma = opts.Comma
	cr.Comment = opts.Comment
	cr.TrimLeadingSpace = opts.TrimLeadingSpace

	// Count rows
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		_, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				r.Seek(0, io.SeekStart)
				break
			}
			return nil, err
		}
		init.Capacity++
	}

	// if headers field is not set, we need to read first 
	// line as header line and decrement capacity
	if len(opts.Headers) == 0 {
		headers, err := cr.Read()
		if err != nil {
			return nil, err
		}
		opts.Headers = headers
		init.Capacity--
	}

	series := make([]dataframe.SeriesAny, 0, len(converters))
	seriesIdx := map[string]int {}

	// Init series and map series name to column index
	for name, converter := range converters {
		series = append(series, converter.series(name, &init))

		for i := range opts.Headers {
			if name == opts.Headers[i] {
				seriesIdx[name] = i
				break
			}
		}
	}

	if len(seriesIdx) != len(series) {
		return nil, errors.New("could not map columns to series")
	}

	df := dataframe.NewDataFrame(series...)

	// lets read CSV line by line
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
		
		row := map[string]any {}

		// Could be transforned into goroutines
		for name, idx := range seriesIdx {
			row[name] = converters[name].value(record[idx])
		}

		df.Append(row)
	}


	return df, nil
}