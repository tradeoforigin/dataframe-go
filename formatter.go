package dataframe

// ValueToStringFormatter is used to convert a value
// into a string. Val can be nil or the concrete
// type stored by the series.
type ValueToStringFormatter func(val any) string