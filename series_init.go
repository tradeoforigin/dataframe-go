package dataframe

// SeriesInit is used to configure the series when it is initialized.
// It allows specifying the size and memory allocation for the underlying data of the series.
type SeriesInit struct {
	
	// Size defines the number of rows in the series.
	// If the series needs to be prefilled with data, this value determines
	// how many rows will be initialized, usually with `nil` or a default value (e.g., "NaN").
	// 
	// This field is required for the initialization of a series.
	Size int
	
	// Capacity specifies how much memory to preallocate for the series.
	// If the size of the series is known in advance, preallocating memory
	// helps to optimize memory management and improve performance.
	// This is especially useful when dealing with large data sets to avoid reallocations.
	//
	// It is not required, but it is recommended to preallocate the capacity if the series size
	// is known in advance.
	Capacity int
}
