# dataframe-go

Dataframes are used for statistics, machine-learning, and data manipulation/exploration. This package is based on [rocketlaunchr/dataframe-go](https://github.com/rocketlaunchr/dataframe-go) and rewritten with Go 1.18 generics. This package is still in progress and all of the [rocketlaunchr/dataframe-go](https://github.com/rocketlaunchr/dataframe-go) features will be added in the future. If you are interested in contribution, your help is welcome. 

## 1. Installation and usage

```
go get -u github.com/tradeoforigin/dataframe-go
```

```go
import "github.com/tradeoforigin/dataframe-go"
```

## 2. Series

Series is generic struct to store any kind of data you wish. Series is also type of `interface SeriesAny` to handle different types in `DataFrame`. 

```go
s := dataframe.NewSeries("weight", nil, 115.5, 93.1)
fmt.Println(s.Table())
```

Output:
```
+-----+---------+
|     | WEIGHT  |
+-----+---------+
| 0:  |  115.5  |
| 1:  |  93.1   |
+-----+---------+
| 2X1 | FLOAT64 |
+-----+---------+
```

Series with type definition:
```go
s := dataframe.NewSeries[float64]("weight", nil, 115, 93.1)
fmt.Println(s.Table())
```

Output:
```
+-----+---------+
|     | WEIGHT  |
+-----+---------+
| 0:  |   115   |
| 1:  |  93.1   |
+-----+---------+
| 2X1 | FLOAT64 |
+-----+---------+
```

You can also define series of your own type:

```go
type Dog struct {
    name string
}

s := dataframe.NewSeries("dogs", nil, 
    Dog { "Abby" }, 
    Dog { "Agas" },
)
fmt.Println(s.Table())
```

Output:
```
+-----+----------+
|     |   DOGS   |
+-----+----------+
| 0:  |  {Abby}  |
| 1:  |  {Agas}  |
+-----+----------+
| 2X1 | MAIN DOG |
+-----+----------+
```

Or series of any type:
```go
s := dataframe.NewSeries[any]("numbers", nil, 10, "ten", 10.0)
fmt.Println(s.Table())
```

Output:
```
+-----+---------+
|     | NUMBERS |
+-----+---------+
| 0:  |   10    |
| 1:  |   ten   |
| 2:  |   10    |
+-----+---------+
| 3X1 |   ANY   |
+-----+---------+
```

### 2.1. Series manipulation

Series provides a few functions for data manipulation:

1. `s.Value(row int, options ...Options) T` to get value at row. This function also provides negative indexing e.g. `s.Value(-1)` to get value from the end of the series `s`.
2. `s.Prepend(val []T, options ...Options)` to preppend one or more values into series.
3. `s.Append(val []T, options ...Options) int` to append one or more values into series. It means that values are added at end of the series `s`.
4. `s.Insert(row int, val []T, options ...Options)` inserts one or more values into series at row.
5. `s.Remove(row int, options ...Options)` to remove data at row.
6. `s.Reset(options ...Options)` clears all of the data from series.
7. `s.Update(row int, val T, options ...Options)` is used to change single value at given row.

Example:
```go
s := dataframe.NewSeries[float64]("numbers", nil, 1, 2, 3) // [1, 2, 3]
s.Append([]float64 { 0, 0 }) // [1, 2, 3, 0, 0]
s.Prepend([] float64 { 0, 0 }) // [0, 0, 1, 2, 3, 0, 0]
s.Insert(2, []float64 { -1 }) // [0, 0, -1, 1, 2, 3, 0, 0]
s.Update(-1, -1) // [0, 0, -1, 1, 2, 3, 0, -1]
s.Remove(0) // [0, -1, 1, 2, 3, 0, -1]
fmt.Println(s.Table())
```

Output:
```
+-----+---------+
| 0:  |    0    |
| 1:  |   -1    |
| 2:  |    1    |
| 3:  |    2    |
| 4:  |    3    |
| 5:  |    0    |
| 6:  |   -1    |
+-----+---------+
| 7X1 | FLOAT64 |
+-----+---------+
```

### 2.2. Fill values randomly

There is possibility to fill series with random values:

```go
s := dataframe.NewSeries("rand", nil, math.NaN(), math.NaN(), math.NaN())
s.FillRand(dataframe.RandFillerFloat64())
```

You can also define your own `RandFiller` as function of type `dataframe.RandFn[T any]`.

### 2.3. Sorting

To sort series values you need to provide `CompareFn[T any]` as series less than function: 

```go
s := dataframe.NewSeries("sorted", nil, 0, 2, 1, 4, 3, 6, 5, 10, 9, 8, 7)

s.SetIsLessThanFunc(dataframe.IsLessThanFunc[int])
s.Sort(ctx) // DESC -> s.Sort(ctx, dataframe.SortOptions { Desc: true })

fmt.Println(s.Table())
```

Output:
```
+------+--------+
|      | SORTED |
+------+--------+
|  0:  |   0    |
|  1:  |   1    |
|  2:  |   2    |
|  3:  |   3    |
|  4:  |   4    |
|  5:  |   5    |
|  6:  |   6    |
|  7:  |   7    |
|  8:  |   8    |
|  9:  |   9    |
| 10:  |   10   |
+------+--------+
| 11X1 |  INT   |
+------+--------+
```

### 2.4. Values iterator

Values iterator is used to iterate series data. Iterator provides options to set:
    
1. `InitialRow` - iterator starts at this row. Can be negative value for indexing from the end of the series.
2. `Step` - iteration steps. Can be negative value to iterate backwards.
3. `DontLock` - if true is passed then series is not locked by iterator.

```go
s := dataframe.NewSeries("iterate", nil, 1, 2, 3)

iterator := s.Iterator()
for iterator.Next() {
    fmt.Println(iterator.Index, "->", iterator.Value)
}
```

Output:
```
0 -> 1
1 -> 2
2 -> 3
```

### 2.5. Apply and Filter

You can apply function to modify values of series. As well as you can filter data of series and `DROP` or `KEEP` values. 

Apply:

```go
s := dataframe.NewSeries("apply", nil, 1., 2., 3.) // *dataframe.Series[float64]
	
applyFn := func (val float64, row, nRows int) float64 {
	return val / 2
}

_, err := s.Apply(ctx, applyFn, dataframe.ApplyOptions { InPlace: true })
if err != nil {
	panic(err)
}

fmt.Println(s.Table())
```

Output: 

```
+-----+---------+
|     |  APPLY  |
+-----+---------+
| 0:  |   0.5   |
| 1:  |    1    |
| 2:  |   1.5   |
+-----+---------+
| 3X1 | FLOAT64 |
+-----+---------+
```

Filter:

```go
s := dataframe.NewSeries("filter", nil, 1., math.NaN(), 3.)
	
filterFn := func (val float64, row, nRows int) (dataframe.FilterAction, error) {
	if math.IsNaN(val) {
		return dataframe.DROP, nil
	}
	return dataframe.KEEP, nil
}

_, err := s.Filter(ctx, filterFn, dataframe.FilterOptions { InPlace: true })
if err != nil {
	panic(err)
}

fmt.Println(s.Table())
```

Output: 

```
+-----+---------+
|     | FILTER  |
+-----+---------+
| 0:  |    1    |
| 1:  |    3    |
+-----+---------+
| 2X1 | FLOAT64 |
+-----+---------+
```

### 2.6. Copy and Equality

You can create copy of series as well as you can compare two different series.

```go
s1 := dataframe.NewSeries[float64]("s1", nil, 1, 2, 3, 4)
s2 := s1.Copy() // copy series s1

eq, err := s.IsEqual(ctx, sc1) // returns true, nil 

// // lines below returns false, nil
// s2.Rename("s2")
// eq, err := s.IsEqual(ctx, sc1, dataframe.IsEqualOptions { CheckName: true }) 
```

## 3. DataFrame

DataFrame is container for Series of any kind. You can think of a Dataframe as an excel spreadsheet. 

```go
x := dataframe.NewSeries("x", nil, 1., 2., 3.)
y := dataframe.NewSeries("y", nil, 1., 2., 3.)

df := dataframe.NewDataFrame(x, y)

fmt.Println(df.Table())
```

Output: 

```
+-----+---------+---------+
|     |    X    |    Y    |
+-----+---------+---------+
| 0:  |    1    |    1    |
| 1:  |    2    |    2    |
| 2:  |    3    |    3    |
+-----+---------+---------+
| 3X2 | FLOAT64 | FLOAT64 |
+-----+---------+---------+
```

### 3.1. DataFrame manipulation

DataFrame provides functions for manipulation with data. Similarly like for the series:

1. `df.Row(row int, options ...Options) map[string]any` returns values for specific row.
2. `df.Prepend(vals any, options ...Options)` adds value to the start of all series
3. `df.Append(vals any, options ...Options)` adds value to the end of all series
4. `df.Insert(row int, vals any, options ...Options)` inserts value to all of the series at specific row.
5. `df.Remove(row int, options ...Options)` removes row at index `row`
6. `df.UpdateRow(row int, vals any, options ...Options)` updatese whole row. Means that all of the series are updated.
7. `df.Update(row int, col any, val any, options ...Options)` - updates value for specific row and column (series)
8. `df.ReorderColumns(newOrder []string, options ...Options) error` changes orders of columns / series indexes. 
9. `df.RemoveSeries(seriesName string, options ...Options) error` removes whole series by name.
10. `df.AddSeries(s SeriesAny, colN *int, options ...Options) error` adds new series into DataFrame
11. `df.Swap(row1, row2 int, options ...Options)` swaps two rows.

In many cases the values should be provided as `map[string]any`, `map[int]any` or `[]any`.

```go
s1 := dataframe.NewSeries[float64]("a", nil, 1, 2, 3, 4)
s2 := dataframe.NewSeries[float64]("b", nil, 1, 2, 3, 4)

df := dataframe.NewDataFrame(s1, s2)

df.Append(map[string]any {
	"a": [] float64 { 0, 0 },
	"b": [] float64 { 0, 0 },
})

df.Prepend(map[string]any {
	"a": [] float64 { 0, 0 },
	"b": [] float64 { 0, 0 },
})

df.Insert(2, map[string]any {
	"a": -1.0,
	"b": -1.0,
})

df.Update(-1, "a", -1.0)

fmt.Println(df.Table())
```

Output:

```
+-----+---------+---------+
|     |    A    |    B    |
+-----+---------+---------+
| 0:  |    0    |    0    |
| 1:  |    0    |    0    |
| 2:  |   -1    |   -1    |
| 3:  |    1    |    1    |
| 4:  |    2    |    2    |
| 5:  |    3    |    3    |
| 6:  |    4    |    4    |
| 7:  |    0    |    0    |
| 8:  |   -1    |    0    |
+-----+---------+---------+
| 9X2 | FLOAT64 | FLOAT64 |
+-----+---------+---------+
```

### 3.2. Fill values randomly

You can fill values with RandFiller at once:

```go
s1 := dataframe.NewSeries("a", nil, math.NaN(), math.NaN(), math.NaN())
s2 := dataframe.NewSeries("b", nil, math.NaN(), math.NaN(), math.NaN())

df := dataframe.NewDataFrame(s1, s2)

df.FillRand(func() any {
	return rand.Float64()
})
```

### 3.3. Sorting

To sort DataFrame you need to provide `CompareFn[T any]` for all of the series as is less than function: 

```go
s1 := dataframe.NewSeries("a", nil, 0, 2, 1, 4, 3, 6, 5, 10, 9, 8, 7)
s2 := dataframe.NewSeries("b", nil, 0, 2, 1, 4, 3, 6, 5, 10, 9, 8, 7)

s1.SetIsLessThanFunc(dataframe.IsLessThanFunc[int])
s2.SetIsLessThanFunc(dataframe.IsLessThanFunc[int])

df := dataframe.NewDataFrame(s1, s2)
	
df.Sort(ctx, []dataframe.SortKey {
	{ Key: "a" }, // Desc: true
	{ Key: "b" }, // Desc: true
})

fmt.Println(df.Table())
```

Output: 

```
+------+-----+-----+
|      |  A  |  B  |
+------+-----+-----+
|  0:  |  0  |  0  |
|  1:  |  1  |  1  |
|  2:  |  2  |  2  |
|  3:  |  3  |  3  |
|  4:  |  4  |  4  |
|  5:  |  5  |  5  |
|  6:  |  6  |  6  |
|  7:  |  7  |  7  |
|  8:  |  8  |  8  |
|  9:  |  9  |  9  |
| 10:  | 10  | 10  |
+------+-----+-----+
| 11X2 | INT | INT |
+------+-----+-----+
```

### 3.4. Values iterator

Values iterator is used to iterate dataframe rows. Iterator provides options to set:
    
1. `InitialRow` - iterator starts at this row. Can be negative value for indexing from the end of the series.
2. `Step` - iteration steps. Can be negative value to iterate backwards.
3. `DontLock` - if true is passed then dataframe is not locked by iterator.

```go
s1 := dataframe.NewSeries("a", nil, 1, 2, 3)
s2 := dataframe.NewSeries("b", nil, 1, 2, 3)

df := dataframe.NewDataFrame(s1, s2)

var iterator = df.Iterator()
for iterator.Next() {
    fmt.Println(iterator.Index, iterator.Value)
}
```

Output:

```
0 map[a:1 b:1]
1 map[a:2 b:2]
2 map[a:3 b:3]
```

### 3.5. Apply and Filter

You can apply function to modify rows of dataframe. As well as you can filter data of dataframe and DROP or KEEP values.

Apply:

```go
y1  := dataframe.NewSeries[float64]("y1", &dataframe.SeriesInit{Size: 24})
y2 := dataframe.NewSeries[float64]("y2", &dataframe.SeriesInit{Size: 24})
	
df := dataframe.NewDataFrame(y1, y2)

//  // Simple example
//  fn := func (vals map[string]any, row, nRows int) map[string]any {
// 	    x := float64(row + 1)
//      return map[string]any{
// 		    "y": math.Sin(2 * math.Pi * x / float64(nRows)),
// 	    }
//  }

fn := func (vals map[string]any, row, nRows int) map[string]any {
	x := float64(row + 1)
	y := math.Sin(2 * math.Pi * x / 24)

	if y == 1 || y == -1 {
		return map[string]any{
			"y1": y,
			"y2": y,
		}
	}

    // We can also update just one column
	return map[string]any{
		"y1": y,
	}
}

_, err := df.Apply(ctx, fn, dataframe.ApplyOptions { InPlace: true })
if err != nil {
	panic(err)
}

fmt.Println(df.Table())
```

Output:

```
+------+------------------------+---------+
|      |           Y1           |   Y2    |
+------+------------------------+---------+
|  0:  |  0.25881904510252074   |   NaN   |
|  1:  |  0.49999999999999994   |   NaN   |
|  2:  |   0.7071067811865475   |   NaN   |
|  3:  |   0.8660254037844386   |   NaN   |
|  4:  |   0.9659258262890683   |   NaN   |
|  5:  |           1            |    1    |
|  6:  |   0.9659258262890683   |   NaN   |
|  7:  |   0.8660254037844388   |   NaN   |
|  8:  |   0.7071067811865476   |   NaN   |
|  9:  |  0.49999999999999994   |   NaN   |
| 10:  |   0.258819045102521    |   NaN   |
| 11:  | 1.2246467991473515e-16 |   NaN   |
| 12:  |  -0.2588190451025208   |   NaN   |
| 13:  |  -0.4999999999999998   |   NaN   |
| 14:  |  -0.7071067811865471   |   NaN   |
| 15:  |  -0.8660254037844384   |   NaN   |
| 16:  |  -0.9659258262890683   |   NaN   |
| 17:  |           -1           |   -1    |
| 18:  |  -0.9659258262890684   |   NaN   |
| 19:  |  -0.8660254037844386   |   NaN   |
| 20:  |  -0.7071067811865477   |   NaN   |
| 21:  |  -0.5000000000000004   |   NaN   |
| 22:  |  -0.2588190451025215   |   NaN   |
| 23:  | -2.449293598294703e-16 |   NaN   |
+------+------------------------+---------+
| 24X2 |        FLOAT64         | FLOAT64 |
+------+------------------------+---------+
```

Filter:

```go
s := dataframe.NewSeries("s", nil, 1, 2, 3, 4, 5)
df := dataframe.NewDataFrame(s)
	
fn := func (vals map[string]any, row, nRows int) (dataframe.FilterAction, error) {
	if row % 2 != 0 {
		return dataframe.DROP, nil
	}
	return dataframe.KEEP, nil
}

_, err := df.Filter(ctx, fn, dataframe.FilterOptions { InPlace: true })
if err != nil {
	panic(err)
}

fmt.Println(df.Table())
```

Output:

```
+-----+-----+
|     |  S  |
+-----+-----+
| 0:  |  1  |
| 1:  |  3  |
| 2:  |  5  |
+-----+-----+
| 3X1 | INT |
+-----+-----+
```

### 3.6. Copy and Equality

You can create copy of dataframe as well as you can compare two different dataframes.

```go
s := dataframe.NewSeries[float64]("s", nil, 1, 2, 3, 4)

df1 := dataframe.NewDataFrame(s)
df2 := df1.Copy() // To copy series s1

eq, err := df1.IsEqual(ctx, df2) // returns true, nil 
```

### 3.7. Import dataframe from CSV

There is possibility to import dataframe directly from CSV:

```go
csvString := `
A,B,C,D
0.0,0.0,0.02,0
0.0,1.6739,0.04,0
0.0,1.6739,0.06,0
0.0,1.673738,0.06,0
0.0,1.6736,0.06,0
0.0,1.673456,0.08,0
0.0,1.67302752,0.08,0
0.0,1.6726333184,0.08,0
1.6681,0.0,0.02,1`

reader := strings.NewReader(csvString)

df, err := csv.Load(ctx, reader, map[string]csv.ConverterAny {
	"A": csv.Float64,
	"B": csv.Float64,
	"C": csv.Float64,
	"D": csv.Float64,
})

if err != nil {
	t.Fatal(err)
}

fmt.Println(df.Table())
```

Output:

```
+-----+---------+---------+--------------+---------+
|     |    D    |    A    |      B       |    C    |
+-----+---------+---------+--------------+---------+
| 0:  |    0    |    0    |      0       |  0.02   |
| 1:  |    0    |    0    |    1.6739    |  0.04   |
| 2:  |    0    |    0    |    1.6739    |  0.06   |
| 3:  |    0    |    0    |   1.673738   |  0.06   |
| 4:  |    0    |    0    |    1.6736    |  0.06   |
| 5:  |    0    |    0    |   1.673456   |  0.08   |
| 6:  |    0    |    0    |  1.67302752  |  0.08   |
| 7:  |    0    |    0    | 1.6726333184 |  0.08   |
| 8:  |    1    | 1.6681  |      0       |  0.02   |
+-----+---------+---------+--------------+---------+
| 9X4 | FLOAT64 | FLOAT64 |   FLOAT64    | FLOAT64 |
+-----+---------+---------+--------------+---------+
```

You can also define custom converter to fit your needs.

### 3.8. Math functions and fakers

There is no need for creating series by string expressions. Math functions for series can be covered by `df.Apply` or `s.Apply` function. The faker can be covered by custom `RandFillers`. Math functions and fakers may be added in future.