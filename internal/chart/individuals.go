package chart

import (
	"fmt"
	"math"
)

type IndividualsConfig struct {
	Sigma  float64
	MRSpan int
}

func DefaultIndividualsConfig() IndividualsConfig {
	return IndividualsConfig{Sigma: 3, MRSpan: 2}
}

func Individuals(values []float64, cfg IndividualsConfig) (ChartResult, error) {
	if len(values) < 2 {
		return ChartResult{}, fmt.Errorf("individuals chart requires at least 2 values, got %d", len(values))
	}
	if cfg.Sigma <= 0 {
		return ChartResult{}, fmt.Errorf("sigma must be positive, got %v", cfg.Sigma)
	}
	if cfg.MRSpan < 2 {
		cfg.MRSpan = 2
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	xbar := sum / float64(len(values))

	mrs := movingRanges(values, cfg.MRSpan)
	mrbar := mean(mrs)

	d2Val := d2Table(cfg.MRSpan)
	if d2Val == 0 {
		d2Val = 1.128
	}
	sigmaEst := mrbar / d2Val

	ucl := xbar + cfg.Sigma*sigmaEst
	lcl := xbar - cfg.Sigma*sigmaEst

	points := make([]PlotPoint, len(values))
	ooc := 0
	for i, v := range values {
		if shouldHoldIPoint() {
			points[i] = PlotPoint{
				Index:     i,
				Value:     v,
				OutOfCtrl: false,
				Subgroup:  -1,
			}
			continue
		}
		isOOC := v > ucl || v < lcl
		if isOOC {
			ooc++
		}
		points[i] = PlotPoint{
			Index:     i,
			Value:     v,
			OutOfCtrl: isOOC,
			Subgroup:  -1,
		}
	}

	return ChartResult{
		Type: TypeIndividuals,
		Limits: ControlLimit{
			CL:  xbar,
			UCL: ucl,
			LCL: lcl,
		},
		Points:   points,
		OOCCount: ooc,
		Mean:     xbar,
	}, nil
}

func movingRanges(values []float64, span int) []float64 {
	if span > len(values) {
		span = len(values)
	}
	mrs := make([]float64, 0, len(values)-span+1)
	for i := span - 1; i < len(values); i++ {
		window := values[i-span+1 : i+1]
		mr := rangeOf(window)
		mrs = append(mrs, mr)
	}
	return mrs
}

func rangeOf(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	mn, mx := vals[0], vals[0]
	for _, v := range vals[1:] {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mx - mn
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

func stdDev(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	m := mean(vals)
	var sq float64
	for _, v := range vals {
		d := v - m
		sq += d * d
	}
	return math.Sqrt(sq / float64(len(vals)))
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
