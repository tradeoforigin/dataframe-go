package tests

import (
	"context"
	"github.com/tradeoforigin/dataframe-go"
	"math"
	"math/rand"
	"testing"
)

func TestDataFrameCopyEqual(t *testing.T) {
	ctx := context.Background()

	s1 := dataframe.NewSeries[float64]("a", nil, 1, 2, 3, 4)
	s2 := dataframe.NewSeries[float64]("b", nil, 1, 2, 3, 4)

	df   := dataframe.NewDataFrame(s1, s2)
	dfc1 := df.Copy()

	eq, err := df.IsEqual(ctx, dfc1)

	if err != nil || !eq {
		t.Fatalf(`df.IsEqual(ctx, dfc1) = %v, %v, want match for true, <nil>`, eq, err)
	}

	df.Series[0].Rename("c")

	eq, err = df.IsEqual(ctx, dfc1, dataframe.IsEqualOptions { CheckName: true })

	if err != nil || eq {
		t.Fatalf(`df.IsEqual(ctx, dfc1) = %v, %v, want match for false, <nil>`, eq, err)
	}

	dfc2 := df.Copy(dataframe.Range(0, 2))
	dfc3 := df.Copy(dataframe.Range(0, 2))

	eq, err = dfc2.IsEqual(ctx, dfc3)

	if err != nil || !eq {
		t.Fatalf(`dfc2.IsEqual(ctx, dfc3) = %v, %v, want match for true, <nil>`, eq, err)
	}

	if dfc3.NRows() != 3 {
		t.Fatalf(`dfc3.NRows() = %v, want match for 3`, dfc3.NRows())
	}

	dfc2.Update(0, "b", float64(100))

	if dfc2.Row(0)["b"] == dfc3.Row(0)["b"] {
		t.Fatalf(`dfc2.Row(0)["b"] == dfc3.Row(0)["b"] is true, want to match false`)
	}

	_ = df.Copy(dataframe.Range(0, df.NRows() - 1))
}

func TestDataFrameManipulation(t *testing.T) {
	s1 := dataframe.NewSeries[float64]("a", nil, 1, 2, 3, 4)
	s2 := dataframe.NewSeries[float64]("b", nil, 1, 2, 3, 4)

	df   := dataframe.NewDataFrame(s1, s2)

	var nRows = df.NRows()

	df.Append(map[string]any {
		"a": [] float64 { 0, 0 },
		"b": [] float64 { 0, 0 },
	})

	df.Prepend(map[string]any {
		"a": [] float64 { 0, 0 },
		"b": [] float64 { 0, 0 },
	})

	if df.NRows() != nRows + 4 {
		t.Fatalf(`df.NRows() = %v, want match for %v`, df.NRows(), nRows + 4)
	}

	nRows = df.NRows()
	
	if df.Row(0)["a"] != 0.0 || df.Row(-1)["b"] != 0.0 {
		t.Fatalf(`df.Row(0)["a"] = %v and df.Row(-1)["b"] = %v, want match for 0, 0`, df.Row(0)["a"], df.Row(-1)["b"])
	}

	df.Insert(2, map[string]any {
		"a": -1.0,
		"b": -1.0,
	})

	if df.Row(2)["a"] != -1.0 {
		t.Fatalf(`df.Row(2)["a"] = %v, want match for -1`, df.Row(2)["a"])
	}

	if df.NRows() != nRows + 1.0 {
		t.Fatalf(`df.NRows() = %v, want match for %v`, df.NRows(), nRows + 1)
	}

	df.Update(-1, "a", -1.0)

	if df.Row(-1)["a"] != -1.0 {
		t.Fatalf(`df.Row(-1)["a"] = %v, want match for -1`, df.Row(-1)["a"])
	}
}

func TestDataFrameIterator(t *testing.T) {
	s1 := dataframe.NewSeries("a", nil, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	s2 := dataframe.NewSeries("b", nil, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	df := dataframe.NewDataFrame(s1, s2)

	var currentVal = 0
	var iterator = df.Iterator()
	for iterator.Next() {
		if iterator.Value["a"] != currentVal {
			t.Fatalf(`iterator.Value["a"] = %v, want match for %v`, iterator.Value["a"], currentVal)
		}
		if iterator.Value["b"] != currentVal {
			t.Fatalf(`iterator.Value["b"] = %v, want match for %v`, iterator.Value["b"], currentVal)
		}
		currentVal += 1
	}

	currentVal = 10
	iterator = df.Iterator(dataframe.IteratorOptions { InitialRow: -1, Step: -2 })
	for iterator.Next() {
		if iterator.Value["a"] != currentVal {
			t.Fatalf(`iterator.Value["a"] = %v, want match for %v`, iterator.Value["a"], currentVal)
		}
		if iterator.Value["b"] != currentVal {
			t.Fatalf(`iterator.Value["b"] = %v, want match for %v`, iterator.Value["b"], currentVal)
		}
		currentVal -= 2
	}

	currentVal = 5
	iterator = df.Iterator(dataframe.IteratorOptions { InitialRow: 5, Step: -1 })
	for iterator.Next() {
		if iterator.Value["a"] != currentVal {
			t.Fatalf(`iterator.Value["a"] = %v, want match for %v`, iterator.Value["a"], currentVal)
		}
		if iterator.Value["b"] != currentVal {
			t.Fatalf(`iterator.Value["b"] = %v, want match for %v`, iterator.Value["b"], currentVal)
		}
		currentVal -= 1
	}
}

func TestDataFrameSort(t *testing.T) {
	ctx := context.Background()

	s1 := dataframe.NewSeries("a", nil, 0, 2, 1, 4, 3, 6, 5, 10, 9, 8, 7)
	s2 := dataframe.NewSeries("b", nil, 0, 2, 1, 4, 3, 6, 5, 10, 9, 8, 7)

	s1.SetIsLessThanFunc(dataframe.IsLessThanFunc[int])
	s2.SetIsLessThanFunc(dataframe.IsLessThanFunc[int])

	df := dataframe.NewDataFrame(s1, s2)
	
	df.Sort(ctx, []dataframe.SortKey {
		{ Key: "a" },
		{ Key: "b" },
	})

	var currentVal = 0
	var iterator = df.Iterator()
	for iterator.Next() {
		if iterator.Value["a"] != currentVal {
			t.Fatalf(`iterator.Value["a"] = %v, want match for %v`, iterator.Value["a"], currentVal)
		}
		if iterator.Value["b"] != currentVal {
			t.Fatalf(`iterator.Value["b"] = %v, want match for %v`, iterator.Value["b"], currentVal)
		}
		currentVal += 1
	}

	df.Sort(ctx, []dataframe.SortKey {
		{ Key: "a", Desc: true },
		{ Key: "b", Desc: true },
	})
	
	currentVal = 10
	iterator = df.Iterator()
	for iterator.Next() {
		if iterator.Value["a"] != currentVal {
			t.Fatalf(`iterator.Value["a"] = %v, want match for %v`, iterator.Value["a"], currentVal)
		}
		if iterator.Value["b"] != currentVal {
			t.Fatalf(`iterator.Value["b"] = %v, want match for %v`, iterator.Value["b"], currentVal)
		}
		currentVal -= 1
	}
}

func TestDataFrameFillRand(t *testing.T) {
	s1 := dataframe.NewSeries("a", nil, math.NaN(), math.NaN(), math.NaN())
	s2 := dataframe.NewSeries("b", nil, math.NaN(), math.NaN(), math.NaN())

	df := dataframe.NewDataFrame(s1, s2)

	df.FillRand(func() any {
		return rand.Float64()
	})

	var iterator = df.Iterator()
	for iterator.Next() {
		if math.IsNaN(iterator.Value["a"].(float64)) {
			t.Fatalf(
				`math.IsNaN(iterator.Value["a"].(float64)) = %v, want match for false`, 
				math.IsNaN(iterator.Value["a"].(float64)),
			)
		}
		if math.IsNaN(iterator.Value["b"].(float64)) {
			t.Fatalf(
				`math.IsNaN(iterator.Value["b"].(float64)) = %v, want match for false`, 
				math.IsNaN(iterator.Value["b"].(float64)),
			)
		}
	}
}
