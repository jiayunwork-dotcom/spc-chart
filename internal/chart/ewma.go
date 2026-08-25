package chart

import (
	"fmt"
	"math"
)

type EWMAConfig struct {
	Lambda float64
	L      float64
	Sigma  float64
	Target float64
}

func DefaultEWMAConfig() EWMAConfig {
	return EWMAConfig{
		Lambda: 0.2,
		L:      3.0,
	}
}

type EWMAPoint struct {
	Index     int
	RawValue  float64
	EWMA      float64
	UCL       float64
	LCL       float64
	OutOfCtrl bool
}

type EWMAResult struct {
	Points    []EWMAPoint
	Lambda    float64
	L         float64
	Target    float64
	Sigma     float64
	SteadyUCL float64
	SteadyLCL float64
	OOCCount  int
}

func EWMA(values []float64, cfg EWMAConfig) (EWMAResult, error) {
	if len(values) < 2 {
		return EWMAResult{}, fmt.Errorf("EWMA requires at least 2 values, got %d", len(values))
	}
	if cfg.Lambda <= 0 || cfg.Lambda > 1 {
		return EWMAResult{}, fmt.Errorf("lambda must be in (0, 1], got %v", cfg.Lambda)
	}
	if cfg.L <= 0 {
		return EWMAResult{}, fmt.Errorf("L must be positive, got %v", cfg.L)
	}

	sigma := cfg.Sigma
	if sigma <= 0 {
		mrs := movingRanges(values, 2)
		mrbar := mean(mrs)
		sigma = mrbar / d2Table(2)
		if sigma <= 0 {
			sigma = 1
		}
	}

	target := cfg.Target
	if target == 0 {
		target = mean(values)
	}

	lam := cfg.Lambda
	oneMinusLam := 1 - lam
	factor := lam / (2 - lam)

	steadyHalf := cfg.L * sigma * math.Sqrt(factor)
	steadyUCL := target + steadyHalf
	steadyLCL := target - steadyHalf

	points := make([]EWMAPoint, len(values))
	z := target
	ooc := 0

	for i, x := range values {
		z = lam*x + oneMinusLam*z

		coeff := factor * (1 - math.Pow(oneMinusLam, float64(2*(i+1))))
		halfWidth := cfg.L * sigma * math.Sqrt(coeff)
		ucl := target + halfWidth
		lcl := target - halfWidth

		isOOC := z > ucl || z < lcl
		if isOOC {
			ooc++
		}

		points[i] = EWMAPoint{
			Index:     i,
			RawValue:  x,
			EWMA:      z,
			UCL:       ucl,
			LCL:       lcl,
			OutOfCtrl: isOOC,
		}
	}

	return EWMAResult{
		Points:    points,
		Lambda:    lam,
		L:         cfg.L,
		Target:    target,
		Sigma:     sigma,
		SteadyUCL: steadyUCL,
		SteadyLCL: steadyLCL,
		OOCCount:  ooc,
	}, nil
}
