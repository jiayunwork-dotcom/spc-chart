package spc

import (
	"strings"
	"testing"
)

func TestAnalyzeBasic(t *testing.T) {
	vals := []float64{9, 10, 11}
	res, err := Analyze(vals, 3)
	if err != nil {
		t.Fatal(err)
	}
	if res.Mean != 10 {
		t.Fatalf("mean = %v, want 10", res.Mean)
	}
	if !approx(res.StdDev, 0.8165, 1e-3) {
		t.Fatalf("stddev = %v, want ~0.8165", res.StdDev)
	}
	if !approx(res.UCL, 12.4495, 1e-3) || !approx(res.LCL, 7.5505, 1e-3) {
		t.Fatalf("UCL/LCL = %v/%v, want ~12.4495/7.5505", res.UCL, res.LCL)
	}
	if res.Outlier != 0 {
		t.Fatalf("expected 0 outliers, got %d", res.Outlier)
	}
}

func TestAnalyzeOutlier(t *testing.T) {
	vals := make([]float64, 0, 21)
	for i := 0; i < 20; i++ {
		vals = append(vals, 10)
	}
	vals = append(vals, 100)
	res, err := Analyze(vals, 3)
	if err != nil {
		t.Fatal(err)
	}
	if res.Outlier != 1 {
		t.Fatalf("expected 1 outlier, got %d", res.Outlier)
	}
	if !res.Points[len(res.Points)-1].Outlier {
		t.Fatal("last point should be flagged outlier")
	}
}

func approx(a, b, eps float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= eps
}

func TestAnalyzeEmpty(t *testing.T) {
	if _, err := Analyze(nil, 3); err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseValues(t *testing.T) {
	in := "1.0\n2.0\n3.0\n"
	vals, err := ParseValues(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 3 {
		t.Fatalf("parsed %d values, want 3", len(vals))
	}
	vals, _ = ParseValues(strings.NewReader("1, 2\n3\t4"))
	if len(vals) != 4 {
		t.Fatalf("mixed parse got %d, want 4", len(vals))
	}
}
