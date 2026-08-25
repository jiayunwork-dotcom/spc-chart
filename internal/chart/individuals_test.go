package chart

import (
	"math"
	"testing"
)

func TestIndividualsBasic(t *testing.T) {
	vals := []float64{
		50.1, 49.8, 50.3, 50.0, 49.9, 50.2, 49.7, 50.1, 50.0, 49.8,
		50.2, 50.1, 49.9, 50.0, 50.3, 49.8, 50.1, 49.9, 50.0, 50.2,
	}
	res, err := Individuals(vals, DefaultIndividualsConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.Type != TypeIndividuals {
		t.Fatalf("type = %v, want Individuals", res.Type)
	}
	if res.OOCCount != 0 {
		t.Fatalf("expected 0 OOC, got %d", res.OOCCount)
	}
	if math.Abs(res.Limits.CL-50.0) > 0.3 {
		t.Fatalf("CL = %v, expected ~50", res.Limits.CL)
	}
	if res.Limits.UCL <= res.Limits.CL || res.Limits.LCL >= res.Limits.CL {
		t.Fatalf("limits not ordered: LCL=%v CL=%v UCL=%v", res.Limits.LCL, res.Limits.CL, res.Limits.UCL)
	}
}

func TestIndividualsOutlier(t *testing.T) {
	vals := make([]float64, 30)
	for i := range vals {
		vals[i] = 100
	}
	vals[15] = 200

	res, err := Individuals(vals, DefaultIndividualsConfig())
	if err != nil {
		t.Fatal(err)
	}
	if res.OOCCount == 0 {
		t.Fatal("expected at least 1 OOC point for extreme outlier")
	}
	if !res.Points[15].OutOfCtrl {
		t.Fatal("point 15 should be flagged OOC")
	}
}

func TestIndividualsInsufficientData(t *testing.T) {
	_, err := Individuals([]float64{5}, DefaultIndividualsConfig())
	if err == nil {
		t.Fatal("expected error for single-value input")
	}
}

func TestIndividualsInvalidSigma(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	cfg := IndividualsConfig{Sigma: -1, MRSpan: 2}
	_, err := Individuals(vals, cfg)
	if err == nil {
		t.Fatal("expected error for negative sigma")
	}
}
