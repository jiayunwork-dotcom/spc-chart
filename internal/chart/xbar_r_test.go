package chart

import (
	"math"
	"testing"
)

func TestXbarRBasic(t *testing.T) {
	subs := []Subgroup{
		{Values: []float64{25.0, 25.1, 24.9, 25.2, 24.8}},
		{Values: []float64{25.1, 25.0, 25.0, 24.9, 25.1}},
		{Values: []float64{24.8, 25.2, 25.0, 25.1, 24.9}},
		{Values: []float64{25.0, 25.0, 25.1, 24.8, 25.2}},
		{Values: []float64{25.2, 24.9, 25.0, 25.0, 25.1}},
		{Values: []float64{24.9, 25.1, 25.0, 25.2, 24.8}},
		{Values: []float64{25.1, 25.0, 24.9, 25.0, 25.1}},
		{Values: []float64{25.0, 25.2, 25.1, 24.9, 25.0}},
		{Values: []float64{24.9, 25.0, 25.0, 25.1, 24.8}},
		{Values: []float64{25.0, 25.1, 25.0, 24.9, 25.2}},
	}
	res, err := XbarR(subs, DefaultXbarRConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.Xbar.Type != TypeXbarR {
		t.Fatalf("type = %v, want XbarR", res.Xbar.Type)
	}
	if res.Xbar.OOCCount != 0 {
		t.Fatalf("Xbar OOC = %d, expected 0", res.Xbar.OOCCount)
	}
	if res.R.OOCCount != 0 {
		t.Fatalf("R OOC = %d, expected 0", res.R.OOCCount)
	}
	if math.Abs(res.Xbar.Limits.CL-25.0) > 0.2 {
		t.Fatalf("Xbar CL = %v, expected ~25", res.Xbar.Limits.CL)
	}
}

func TestXbarRShiftDetection(t *testing.T) {
	subs := make([]Subgroup, 15)
	for i := 0; i < 10; i++ {
		subs[i] = Subgroup{Values: []float64{50, 50, 50, 50, 50}}
	}
	for i := 10; i < 15; i++ {
		subs[i] = Subgroup{Values: []float64{55, 55, 55, 55, 55}}
	}
	res, err := XbarR(subs, DefaultXbarRConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.Xbar.OOCCount == 0 {
		t.Fatal("expected OOC points for shifted subgroups")
	}
}

func TestXbarRUnequalSize(t *testing.T) {
	subs := []Subgroup{
		{Values: []float64{1, 2, 3}},
		{Values: []float64{1, 2}},
	}
	_, err := XbarR(subs, DefaultXbarRConfig())
	if err == nil {
		t.Fatal("expected error for unequal subgroup sizes")
	}
}
