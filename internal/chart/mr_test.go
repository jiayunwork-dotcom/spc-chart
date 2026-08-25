package chart

import (
	"math"
	"testing"
)

func TestMRBasic(t *testing.T) {
	vals := []float64{10.0, 10.1, 9.9, 10.0, 10.2, 9.8, 10.1, 10.0, 9.9, 10.1}
	res, err := MovingRange(vals, DefaultMRConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.Type != TypeMR {
		t.Fatalf("type = %v, want MR", res.Type)
	}
	if res.OOCCount != 0 {
		t.Fatalf("expected 0 OOC, got %d", res.OOCCount)
	}
	if res.Limits.CL <= 0 {
		t.Fatalf("CL = %v, expected positive", res.Limits.CL)
	}
	if res.Limits.UCL <= res.Limits.CL {
		t.Fatalf("UCL=%v should exceed CL=%v", res.Limits.UCL, res.Limits.CL)
	}
}

func TestMRSpikeDetection(t *testing.T) {
	vals := make([]float64, 20)
	for i := range vals {
		vals[i] = 30
	}
	vals[10] = 50

	res, err := MovingRange(vals, DefaultMRConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.OOCCount == 0 {
		t.Fatal("expected OOC for spike")
	}
}

func TestMRPointCount(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res, err := MovingRange(vals, MRConfig{Sigma: 3, Span: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Points) != 8 {
		t.Fatalf("expected 8 MR points, got %d", len(res.Points))
	}
}

func TestMRLimitsConsistency(t *testing.T) {
	vals := []float64{5, 6, 7, 8, 9, 10, 11, 12}
	res, err := MovingRange(vals, DefaultMRConfig())
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.Limits.LCL) > 1e-9 {
		t.Fatalf("LCL = %v, expected 0 for span=2", res.Limits.LCL)
	}
}
