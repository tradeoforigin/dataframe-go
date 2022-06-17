package tests

import (
	"context"
	"dataframe-go"
	"math"
	"testing"
)

func TestSeriesCopyEqual(t *testing.T) {
	ctx := context.Background()

    s := dataframe.NewSeries[float64]("a", nil, 1, 2, 3, 4)
	sc1 := s.Copy()

	sc1.Rename("b")

	eq, err := s.IsEqual(ctx, sc1)

	if err != nil || !eq {
		t.Fatalf(`s.IsEqual(ctx, sc1) = %v, %v, want match for true, <nil>`, eq, err)
	}

	eq, err = s.IsEqual(ctx, sc1, dataframe.IsEqualOptions { CheckName: true })

	if err != nil || eq {
		t.Fatalf(`s.IsEqual(ctx, sc1) = %v, %v, want match for false, <nil>`, eq, err)
	}

	sc2 := s.Copy(dataframe.Range(0, 2))
	sc3 := s.Copy(dataframe.Range(0, 2))

	eq, err = sc2.IsEqual(ctx, sc3)

	if err != nil || !eq {
		t.Fatalf(`sc2.IsEqual(ctx, sc3) = %v, %v, want match for true, <nil>`, eq, err)
	}

	if sc3.NRows() != 3 {
		t.Fatalf(`sc3.NRows() = %v, want match for 3`, sc3.NRows())
	}

	sc2.Update(0, 100)

	if sc2.Value(0) == sc3.Value(0) {
		t.Fatalf(`sc2.Value(0) == sc3.Value(0) is true, want to match false`)
	}

	_ = s.Copy(dataframe.Range(0, s.NRows() - 1))
}

func TestSeriesValuesManipulation(t *testing.T) {
	s := dataframe.NewSeries[float64]("a", nil, 1, 2, 3, 4)
	var nRows = s.NRows()

	s.Append([]float64 { 0, 0 })
	s.Prepend([] float64 { 0, 0 })

	if s.NRows() != nRows + 4 {
		t.Fatalf(`s.NRows() = %v, want match for %v`, s.NRows(), nRows + 4)
	}

	nRows = s.NRows()

	if s.Value(0) != 0 || s.Value(1) != 0 {
		t.Fatalf(`s.Value(0) = %v and s.Value(1) = %v, want match for 0, 0`, s.Value(0), s.Value(1))
	}

	if s.Value(-1) != 0 || s.Value(-2) != 0 {
		t.Fatalf(`s.Value(1) = %v and s.Value(-2) = %v, want match for 0, 0`, s.Value(-1), s.Value(-2))
	}

	s.Insert(2, []float64 { -1 })

	if s.Value(2) != -1 {
		t.Fatalf(`s.Value(2) = %v, want match for -1`, s.Value(2))
	}

	if s.NRows() != nRows + 1 {
		t.Fatalf(`s.NRows() = %v, want match for %v`, s.NRows(), nRows + 1)
	}

	s.Update(-1, -1)
	if s.Value(-1) != -1 {
		t.Fatalf(`s.Value(-1) = %v, want match for -1`, s.Value(2))
	}
}

func TestSeriesIterator(t *testing.T) {
	s := dataframe.NewSeries("a", nil, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	var currentVal = 0
	var iterator = s.Iterator()
	for iterator.Next() {
		if iterator.Value != currentVal {
			t.Fatalf(`iterator.Value = %v, want match for %v`, iterator.Value, currentVal)
		}
		currentVal += 1
	}

	currentVal = 10
	iterator = s.Iterator(dataframe.IteratorOptions { InitialRow: -1, Step: -2 })
	for iterator.Next() {
		if iterator.Value != currentVal {
			t.Fatalf(`iterator.Value = %v, want match for %v`, iterator.Value, currentVal)
		}
		currentVal -= 2
	}

	currentVal = 5
	iterator = s.Iterator(dataframe.IteratorOptions { InitialRow: 5, Step: -1 })
	for iterator.Next() {
		if iterator.Value != currentVal {
			t.Fatalf(`iterator.Value = %v, want match for %v`, iterator.Value, currentVal)
		}
		currentVal -= 1
	}
}

func TestSeriesSort(t *testing.T) {
	ctx := context.Background()

	s := dataframe.NewSeries("a", nil, 0, 2, 1, 4, 3, 6, 5, 10, 9, 8, 7)

	s.SetIsLessThanFunc(dataframe.IsLessThanFunc[int])
	s.Sort(ctx)

	var currentVal = 0
	var iterator = s.Iterator()
	for iterator.Next() {
		if iterator.Value != currentVal {
			t.Fatalf(`iterator.Value = %v, want match for %v`, iterator.Value, currentVal)
		}
		currentVal += 1
	}

	s.Sort(ctx, dataframe.SortOptions { Desc: true })

	currentVal = 10
	iterator = s.Iterator()
	for iterator.Next() {
		if iterator.Value != currentVal {
			t.Fatalf(`iterator.Value = %v, want match for %v`, iterator.Value, currentVal)
		}
		currentVal -= 1
	}
}

func TestSeriesFillRand(t *testing.T) {
	s := dataframe.NewSeries("a", nil, math.NaN(), math.NaN(), math.NaN())
	s.FillRand(dataframe.RandFillerFloat64())

	var iterator = s.Iterator()
	for iterator.Next() {
		if math.IsNaN(iterator.Value) {
			t.Fatalf(`math.IsNaN(iterator.Value) = %v, want match for false`, math.IsNaN(iterator.Value))
		}
	}
}

func TestSeriesApply(t *testing.T) {
	ctx := context.Background()

	s := dataframe.NewSeries("a", nil, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	
	applySeriesFn := func (val int, row, nRows int) int {
		return val * val
	}

	ns, err := s.Apply(ctx, applySeriesFn)
	if err != nil {
		t.Fatal(err)
	}

	if ns.Value(-1) != 100 || s.Value(-1) == 100 {
		t.Fatalf(`ns.Value(-1) = %v and s.Value(-1) = %v, want match for 100, 100`, ns.Value(-1), s.Value(-1))
	}

	_, err = s.Apply(ctx, applySeriesFn, dataframe.ApplyOptions { InPlace: true })
	if err != nil {
		t.Fatal(err)
	}
	
	if eq, err := s.IsEqual(ctx, ns); err != nil || !eq {
		t.Fatalf(`s.IsEqual(ctx, ns) = (%v, %v), want match for true, <nil>`, eq, err)
	}
}

func TestSeriesFilter(t *testing.T) {
	ctx := context.Background()

	s := dataframe.NewSeries("a", nil, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	ns, err := s.Filter(ctx, func (val int, row, nRows int) (dataframe.FilterAction, error) {
		if row % 2 == 0 {
			return dataframe.KEEP, nil
		}
		return dataframe.DROP, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Filter(ctx, func (val int, row, nRows int) (dataframe.FilterAction, error) {
		if row % 2 == 1 {
			return dataframe.DROP, nil
		}
		return dataframe.KEEP, nil
	}, dataframe.FilterOptions { InPlace: true })

	if err != nil {
		t.Fatal(err)
	}

	if s.NRows() != 5 {
		t.Fatalf(`s.NRows() = %v, want match for 5`, s.NRows())
	}

	if eq, err := s.IsEqual(ctx, ns); err != nil || !eq {
		t.Fatalf(`s.IsEqual(ctx, ns) = (%v, %v), want match for true, <nil>`, eq, err)
	}
}
