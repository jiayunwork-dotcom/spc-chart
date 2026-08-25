package chart

import (
	"math"
	"testing"
)

func TestEWMAStable(t *testing.T) {
	vals := make([]float64, 30)
	for i := range vals {
		vals[i] = 20
	}
	cfg := EWMAConfig{
		Lambda: 0.2,
		L:      3,
		Sigma:  1,
		Target: 20,
	}
	res, err := EWMA(vals, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if res.OOCCount != 0 {
		t.Fatalf("expected 0 OOC for stable process, got %d", res.OOCCount)
	}
	for i, p := range res.Points {
		if math.Abs(p.EWMA-20) > 1e-9 {
			t.Fatalf("point %d: EWMA=%v, expected 20", i, p.EWMA)
		}
	}
}

func TestEWMAShiftDetection(t *testing.T) {
	vals := make([]float64, 50)
	for i := 0; i < 20; i++ {
		vals[i] = 50
	}
	for i := 20; i < 50; i++ {
		vals[i] = 55
	}
	cfg := EWMAConfig{
		Lambda: 0.2,
		L:      3,
		Sigma:  1,
		Target: 50,
	}
	res, err := EWMA(vals, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if res.OOCCount == 0 {
		t.Fatal("expected OOC after 5-sigma sustained shift")
	}
}

func TestEWMATimeVaryingLimits(t *testing.T) {
	vals := make([]float64, 100)
	for i := range vals {
		vals[i] = 0
	}
	cfg := EWMAConfig{
		Lambda: 0.1,
		L:      2.7,
		Sigma:  1,
		Target: 0,
	}
	res, err := EWMA(vals, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if res.Points[0].UCL >= res.Points[99].UCL {
		t.Fatalf("expected UCL to widen: UCL[0]=%v >= UCL[99]=%v",
			res.Points[0].UCL, res.Points[99].UCL)
	}
	diff := math.Abs(res.Points[99].UCL - res.SteadyUCL)
	if diff > 0.01 {
		t.Fatalf("UCL[99]=%v not close to steady %v (diff=%v)", res.Points[99].UCL, res.SteadyUCL, diff)
	}
}

func TestEWMAInvalidLambda(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	_, err := EWMA(vals, EWMAConfig{Lambda: 0, L: 3, Sigma: 1})
	if err == nil {
		t.Fatal("expected error for lambda=0")
	}
	_, err = EWMA(vals, EWMAConfig{Lambda: 1.5, L: 3, Sigma: 1})
	if err == nil {
		t.Fatal("expected error for lambda>1")
	}
}
