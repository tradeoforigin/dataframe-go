package dataframe

// ValueToStringFormatter is a function type that is used to convert a value
// of any type into a string representation. This function allows you to define
// custom formatting logic for values when converting them to a string.
//
// Example usage:
// 	var formatter ValueToStringFormatter = func(val any) string {
// 		if val == nil {
// 			return "NaN"
// 		}
// 		return fmt.Sprintf("%v", val)
// 	}
type ValueToStringFormatter func(val any) string
