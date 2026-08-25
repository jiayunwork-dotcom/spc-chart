package chart

import (
	"math"
	"testing"
)

func TestXbarSBasic(t *testing.T) {
	subs := []Subgroup{
		{Values: []float64{10.0, 10.1, 10.2, 9.9, 10.0}},
		{Values: []float64{10.1, 10.0, 9.9, 10.0, 10.1}},
		{Values: []float64{9.9, 10.0, 10.1, 10.0, 9.8}},
		{Values: []float64{10.0, 10.0, 10.2, 9.9, 10.1}},
		{Values: []float64{10.1, 9.9, 10.0, 10.0, 10.1}},
		{Values: []float64{10.0, 10.1, 10.0, 9.8, 10.2}},
		{Values: []float64{9.9, 10.0, 10.0, 10.1, 10.0}},
		{Values: []float64{10.2, 10.0, 9.9, 10.1, 10.0}},
	}
	res, err := XbarS(subs, DefaultXbarSConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.Xbar.OOCCount != 0 {
		t.Fatalf("Xbar OOC = %d, expected 0", res.Xbar.OOCCount)
	}
	if res.S.OOCCount != 0 {
		t.Fatalf("S OOC = %d, expected 0", res.S.OOCCount)
	}
	if res.S.Limits.CL <= 0 {
		t.Fatalf("Sbar = %v, expected positive", res.S.Limits.CL)
	}
	if math.Abs(res.Xbar.Limits.CL-10) > 0.2 {
		t.Fatalf("Xbar CL = %v, expected ~10", res.Xbar.Limits.CL)
	}
}

func TestXbarSVariabilityOOC(t *testing.T) {
	subs := make([]Subgroup, 10)
	for i := 0; i < 9; i++ {
		subs[i] = Subgroup{Values: []float64{50, 50, 50, 50, 50}}
	}
	subs[9] = Subgroup{Values: []float64{40, 60, 30, 70, 50}}
	res, err := XbarS(subs, DefaultXbarSConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.S.OOCCount == 0 {
		t.Fatal("expected S chart to flag high-variability subgroup")
	}
}

func TestXbarSEmptySubgroups(t *testing.T) {
	_, err := XbarS(nil, DefaultXbarSConfig())
	if err == nil {
		t.Fatal("expected error for nil subgroups")
	}
}
