package dataframe

import "errors"

func (r RangeOptions) Limits(length int) (int, int, error) {
	var start, end = r.Start, -1

	if r.End != nil {
		end = *r.End
	}

	if start < 0 {
		// negative
		start = length + start
	}

	if end < 0 {
		// negative
		end = length + end
	}

	if start < 0 || end < 0 {
		return 0, 0, errors.New("invalid range")
	}

	if start > end {
		return 0, 0, errors.New("invalid range")
	}

	if start >= length || end >= length {
		return 0, 0, errors.New("invalid range")
	}

	return start, end, nil
}
