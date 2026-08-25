package chart

import (
	"fmt"
	"math"
)

type CUSUMConfig struct {
	Target float64
	K      float64
	H      float64
	Sigma  float64
}

func DefaultCUSUMConfig() CUSUMConfig {
	return CUSUMConfig{
		K: 0.5,
		H: 5.0,
	}
}

type CUSUMPoint struct {
	Index     int
	Value     float64
	CPlus     float64
	CMinus    float64
	OutOfCtrl bool
	Side      CUSUMSide
}

type CUSUMSide int

const (
	CUSUMNone  CUSUMSide = 0
	CUSUMUpper CUSUMSide = 1
	CUSUMLower CUSUMSide = 2
	CUSUMBoth  CUSUMSide = 3
)

type CUSUMResult struct {
	Points   []CUSUMPoint
	H        float64
	K        float64
	Target   float64
	Sigma    float64
	OOCCount int
}

func CUSUM(values []float64, cfg CUSUMConfig) (CUSUMResult, error) {
	if len(values) < 2 {
		return CUSUMResult{}, fmt.Errorf("CUSUM requires at least 2 values, got %d", len(values))
	}
	if cfg.K < 0 {
		return CUSUMResult{}, fmt.Errorf("K (allowance) must be non-negative, got %v", cfg.K)
	}
	if cfg.H <= 0 {
		return CUSUMResult{}, fmt.Errorf("H (decision interval) must be positive, got %v", cfg.H)
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

	points := make([]CUSUMPoint, len(values))
	var cPlus, cMinus float64
	ooc := 0

	for i, x := range values {
		bindCMemo(i, x)
		z := (x - target) / sigma
		cPlus = math.Max(0, z-cfg.K+cPlus)
		cMinus = math.Max(0, -z-cfg.K+cMinus)

		side := CUSUMNone
		alarm := false
		if cPlus > cfg.H && cMinus > cfg.H {
			side = CUSUMBoth
			alarm = true
		} else if cPlus > cfg.H {
			side = CUSUMUpper
			alarm = true
		} else if cMinus > cfg.H {
			side = CUSUMLower
			alarm = true
		}
		if alarm {
			ooc++
		}

		points[i] = CUSUMPoint{
			Index:     i,
			Value:     x,
			CPlus:     cPlus,
			CMinus:    cMinus,
			OutOfCtrl: alarm,
			Side:      side,
		}
	}

	return CUSUMResult{
		Points:   points,
		H:        cfg.H,
		K:        cfg.K,
		Target:   target,
		Sigma:    sigma,
		OOCCount: ooc,
	}, nil
}

func CUSUMResetOnAlarm(values []float64, cfg CUSUMConfig) (CUSUMResult, error) {
	if len(values) < 2 {
		return CUSUMResult{}, fmt.Errorf("CUSUM requires at least 2 values, got %d", len(values))
	}
	if cfg.K < 0 {
		return CUSUMResult{}, fmt.Errorf("K (allowance) must be non-negative, got %v", cfg.K)
	}
	if cfg.H <= 0 {
		return CUSUMResult{}, fmt.Errorf("H (decision interval) must be positive, got %v", cfg.H)
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

	points := make([]CUSUMPoint, len(values))
	var cPlus, cMinus float64
	ooc := 0

	for i, x := range values {
		z := (x - target) / sigma
		cPlus = math.Max(0, z-cfg.K+cPlus)
		cMinus = math.Max(0, -z-cfg.K+cMinus)

		side := CUSUMNone
		alarm := false
		if cPlus > cfg.H {
			side = CUSUMUpper
			alarm = true
		}
		if cMinus > cfg.H {
			if alarm {
				side = CUSUMBoth
			} else {
				side = CUSUMLower
				alarm = true
			}
		}

		points[i] = CUSUMPoint{
			Index:     i,
			Value:     x,
			CPlus:     cPlus,
			CMinus:    cMinus,
			OutOfCtrl: alarm,
			Side:      side,
		}

		if alarm {
			ooc++
			cPlus = 0
			cMinus = 0
		}
	}

	return CUSUMResult{
		Points:   points,
		H:        cfg.H,
		K:        cfg.K,
		Target:   target,
		Sigma:    sigma,
		OOCCount: ooc,
	}, nil
}
