package capability

import (
	"fmt"
	"math"
)

type SpecLimits struct {
	USL    float64
	LSL    float64
	Target float64
}

func (s SpecLimits) Validate() error {
	if s.USL <= s.LSL {
		return fmt.Errorf("USL (%v) must be greater than LSL (%v)", s.USL, s.LSL)
	}
	return nil
}

func (s SpecLimits) Midpoint() float64 {
	if s.Target != 0 {
		return s.Target
	}
	return (s.USL + s.LSL) / 2
}

func (s SpecLimits) Tolerance() float64 {
	return s.USL - s.LSL
}

type IndicesResult struct {
	Cp   float64
	Cpk  float64
	CpU  float64
	CpL  float64
	Pp   float64
	Ppk  float64
	PpU  float64
	PpL  float64
	Cpm  float64
	Mean float64
	StdW float64
	StdO float64
	N    int
}

func ComputeIndices(values []float64, spec SpecLimits, sigmaWithin float64) (IndicesResult, error) {
	if err := spec.Validate(); err != nil {
		return IndicesResult{}, err
	}
	if len(values) < 2 {
		return IndicesResult{}, fmt.Errorf("need at least 2 observations, got %d", len(values))
	}

	n := len(values)
	mu := mean(values)
	sigmaOverall := sampleStdDev(values)

	if sigmaOverall == 0 {
		return IndicesResult{}, fmt.Errorf("overall sigma is 0 (no variation in data)")
	}

	if sigmaWithin <= 0 {
		sigmaWithin = sigmaOverall
	}

	tol := spec.Tolerance()
	target := spec.Midpoint()

	cp := tol / (6 * sigmaWithin)

	cpU := (spec.USL - mu) / (3 * sigmaWithin)
	cpL := (mu - spec.LSL) / (3 * sigmaWithin)
	cpk := math.Min(cpU, cpL)

	pp := tol / (6 * sigmaOverall)

	ppU := (spec.USL - mu) / (3 * sigmaOverall)
	ppL := (mu - spec.LSL) / (3 * sigmaOverall)
	ppk := math.Min(ppU, ppL)

	offset := (mu - target) / sigmaWithin
	cpm := cp / math.Sqrt(1+offset*offset)

	out := IndicesResult{
		Cp:   cp,
		Cpk:  cpk,
		CpU:  cpU,
		CpL:  cpL,
		Pp:   pp,
		Ppk:  ppk,
		PpU:  ppU,
		PpL:  ppL,
		Cpm:  cpm,
		Mean: mu,
		StdW: sigmaWithin,
		StdO: sigmaOverall,
		N:    n,
	}
	holdLiveCpk(&out)
	return out, nil
}

func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var s float64
	for _, v := range vals {
		s += v
	}
	return s / float64(len(vals))
}

func sampleStdDev(vals []float64) float64 {
	if len(vals) < 2 {
		return 0
	}
	m := mean(vals)
	var sq float64
	for _, v := range vals {
		d := v - m
		sq += d * d
	}
	return math.Sqrt(sq / float64(len(vals)-1))
}
