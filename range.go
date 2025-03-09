package dataframe

import "errors"

// Limits calculates and returns the start and end indices for a given range
// based on the `RangeOptions` struct. The function takes into account 
// negative values for both the `Start` and `End` parameters to allow for 
// indexing relative to the end of the range (e.g., -1 for the last element).
// It also ensures that the calculated indices are within the bounds of the
// provided length, and returns an error if the range is invalid.
//
// Parameters:
//   length (int): The total length of the data, used to calculate valid start and end indices.
//
// Returns:
//   (int, int, error): The function returns the start and end indices, 
//   along with any potential error if the range is invalid. 
//   If the range is valid, the error will be nil.
//
// Errors:
//   - "invalid range": This error is returned if any of the following are true:
//     - The start or end index is negative but out of bounds (after adjustment).
//     - The start index is greater than the end index.
//     - The start or end index is outside the valid range [0, length).
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
