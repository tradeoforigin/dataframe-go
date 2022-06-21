package tests

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

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