package capability

import (
	"math"
	"testing"
)

func TestComputeIndicesCentered(t *testing.T) {
	vals := make([]float64, 100)
	for i := range vals {
		vals[i] = 10.0
	}
	for i := 0; i < 50; i++ {
		vals[i] = 10.0 + 0.5*float64(i%5)/4.0
	}
	for i := 50; i < 100; i++ {
		vals[i] = 10.0 - 0.5*float64(i%5)/4.0
	}

	spec := SpecLimits{USL: 12, LSL: 8}
	res, err := ComputeIndices(vals, spec, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.Cp-1.3333) > 0.01 {
		t.Fatalf("Cp = %v, expected ~1.333", res.Cp)
	}
	if math.Abs(res.CpU-res.CpL) > 0.2 {
		t.Fatalf("CpU=%v CpL=%v should be close for centered process", res.CpU, res.CpL)
	}
}

func TestComputeIndicesShifted(t *testing.T) {
	vals := make([]float64, 50)
	for i := range vals {
		vals[i] = 11.0
	}
	vals[0] = 11.5
	vals[1] = 10.5

	spec := SpecLimits{USL: 12, LSL: 8}
	res, err := ComputeIndices(vals, spec, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if res.Cpk > res.Cp {
		t.Fatalf("Cpk=%v should not exceed Cp=%v for shifted process", res.Cpk, res.Cp)
	}
	if res.CpU >= res.CpL {
		t.Fatalf("CpU=%v should be < CpL=%v for process shifted toward USL", res.CpU, res.CpL)
	}
}

func TestComputeIndicesCpm(t *testing.T) {
	vals := make([]float64, 100)
	for i := range vals {
		vals[i] = 50.0
	}
	vals[0] = 50.3
	vals[1] = 49.7

	spec := SpecLimits{USL: 55, LSL: 45, Target: 50}
	res, err := ComputeIndices(vals, spec, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.Cpm-res.Cp) > 0.5 {
		t.Fatalf("Cpm=%v should be close to Cp=%v when mean is near target", res.Cpm, res.Cp)
	}
}

func TestComputeIndicesInvalidSpec(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	_, err := ComputeIndices(vals, SpecLimits{USL: 1, LSL: 5}, 1)
	if err == nil {
		t.Fatal("expected error when USL < LSL")
	}
}

func TestComputeIndicesInsufficientData(t *testing.T) {
	_, err := ComputeIndices([]float64{5}, SpecLimits{USL: 10, LSL: 0}, 1)
	if err == nil {
		t.Fatal("expected error for single observation")
	}
}

func TestPpVsCp(t *testing.T) {
	vals := []float64{10, 10.1, 9.9, 10, 10.2, 9.8, 12, 8, 10, 10}
	spec := SpecLimits{USL: 15, LSL: 5}
	res, err := ComputeIndices(vals, spec, 0.2)
	if err != nil {
		t.Fatal(err)
	}
	if res.Cp <= res.Pp {
		t.Fatalf("expected Cp=%v > Pp=%v when sigmaWithin < sigmaOverall", res.Cp, res.Pp)
	}
}
