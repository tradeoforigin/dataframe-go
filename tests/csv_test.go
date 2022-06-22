package tests

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/tradeoforigin/dataframe-go"
	"github.com/tradeoforigin/dataframe-go/utils/csv"
)

func TestCSVLoad(t *testing.T) {
	ctx := context.Background()

	var content, err = ioutil.ReadFile("data/data+header.csv")
	if err != nil {
		t.Fatal(err)
	}

	var reader = strings.NewReader(string(content))

	df1, err := csv.Load(ctx, reader, map[string]csv.ConverterAny {
		"A": csv.Float64,
		"B": csv.Float64,
		"C": csv.Float64,
		"D": csv.Float64,
	})

	if err != nil {
		t.Fatal(err)
	}

	content, err = ioutil.ReadFile("data/data-header.csv")
	if err != nil {
		t.Fatal(err)
	}

	reader = strings.NewReader(string(content))

	df2, err := csv.Load(ctx, reader, map[string]csv.ConverterAny {
		"A": csv.Float64,
		"B": csv.Float64,
		"C": csv.Float64,
		"D": csv.Float64,
	}, csv.LoadOptions{Headers: []string {"A", "B", "C", "D"}})

	if err != nil {
		t.Fatal(err)
	}

	df1.ReorderColumns([]string { "A", "B", "C", "D" })
	df2.ReorderColumns([]string { "A", "B", "C", "D" })
	
	if eq, err := df1.IsEqual(ctx, df2); !eq || err != nil {
		t.Fatalf(`eq, err := df1.IsEqual(ctx, df2) = %v, %v, want match for true, <nil>`, eq, err)
	}
}

func TestCSVExport(t *testing.T) {
	ctx := context.Background()

	s1 := dataframe.NewSeries("str", nil, "one", "one,two", "one,two,three")
	s2 := dataframe.NewSeries("num", nil, 1, 12, 123)

	df1 := dataframe.NewDataFrame(s1, s2)

	f, err := os.OpenFile("data/export.csv", os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        t.Fatal(err)
    }

	err = csv.Export(ctx, f, df1)
	if err != nil {
		t.Fatal(err)
	}

	f.Close()

	content, err := ioutil.ReadFile("data/export.csv")
	if err != nil {
		t.Fatal(err)
	}

	r := strings.NewReader(string(content))

	df2, err := csv.Load(ctx, r, map[string]csv.ConverterAny {
		"str": csv.String,
		"num": csv.Int,
	})

	if err != nil {
		t.Fatal(err)
	}

	df1.ReorderColumns([]string { "str", "num" })
	df2.ReorderColumns([]string { "str", "num" })
	
	if eq, err := df1.IsEqual(ctx, df2); !eq || err != nil {
		t.Fatalf(`eq, err := df1.IsEqual(ctx, df2) = %v, %v, want match for true, <nil>`, eq, err)
	}
}