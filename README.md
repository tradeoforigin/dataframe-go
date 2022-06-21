# dataframe-go

Dataframes are used for statistics, machine-learning, and data manipulation/exploration. This package is based on [rocketlaunchr/dataframe-go](https://github.com/rocketlaunchr/dataframe-go) and rewritten with Go 1.18 generics. This package is still in progress and all of the [rocketlaunchr/dataframe-go](https://github.com/rocketlaunchr/dataframe-go) features will be added in the future. If you are interested in contribution, your help is welcome. 

## Series

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

### Series manipulation

Series provides a few functions for data manipulation. Let series `s` is series of type `dataframe.Series[float64]`:

1. `s.Value(row int, options ...Options) float64` to get value at row. This function also provides negative indexing e.g. `s.Value(-1)` to get value from the end of the series `s`.
2. `s.Prepend(val []float64, options ...Options)` to preppend one or many values into series.
3. `s.Append(val []float64, options ...Options) int` to append one or many values into series. It means that values are added at end of the series `s`.
4. `s.Insert(row int, val []T, options ...Options)` inserts one or many values into series at row.
5. `s.Remove(row int, options ...Options)` to remove data at row.
6. `s.Reset(options ...Options)` clears all of the data from series.
7. `s.Update(row int, val T, options ...Options)` is used to change single value at given row.

Example:
```go
s := dataframe.NewSeries[float64]("numbers", nil, 1, 2, 3)
s.Append([]float64 { 0, 0 })
s.Prepend([] float64 { 0, 0 })
s.Insert(2, []float64 { -1 })
s.Update(-1, -1)
s.Remove(0)
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

### Fill rand

There is possibility to fill series by random values:

```go
s := dataframe.NewSeries("a", nil, math.NaN(), math.NaN(), math.NaN())
s.FillRand(dataframe.RandFillerFloat64())
```

You can also define your own `RandFiller` as function of type `dataframe.RandFn[T any]`.

### Sorting

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

### Values iterator

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

### Apply and Filter

You can apply function to modify values of series. As well as you can filter data of series and `DROP` or `KEEP` values. 

Apply:

```go
s := dataframe.NewSeries("apply", nil, 1., 2., 3.)
	
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

_, err := s.Filter(ctx, filterFn, dataframe.ApplyOptions { InPlace: true })
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

### Copy and Equality

You can create copy of series as well as you can equal two different series.

```go
s1 := dataframe.NewSeries[float64]("s1", nil, 1, 2, 3, 4)
s2 := s1.Copy() // To copy series s1

eq, err := s.IsEqual(ctx, sc1) // returns true, nil 

// // lines below returns false, nil
// s2.Rename("s2")
// eq, err := s.IsEqual(ctx, sc1, dataframe.IsEqualOptions { CheckName: true }) 
```

