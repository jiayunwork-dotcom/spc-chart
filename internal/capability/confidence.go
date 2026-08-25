package capability

import (
	"fmt"
	"math"
)

type ConfidenceInterval struct {
	Lower    float64
	Upper    float64
	Level    float64
	PointEst float64
}

func CpkConfidence(cpk float64, n int, confidence float64) (ConfidenceInterval, error) {
	if n < 2 {
		return ConfidenceInterval{}, fmt.Errorf("need n >= 2, got %d", n)
	}
	if confidence <= 0 || confidence >= 1 {
		return ConfidenceInterval{}, fmt.Errorf("confidence must be in (0,1), got %v", confidence)
	}
	if cpk <= 0 {
		return ConfidenceInterval{}, fmt.Errorf("Cpk must be positive for CI computation, got %v", cpk)
	}

	alpha := 1 - confidence
	z := normalQuantile(1 - alpha/2)

	nf := float64(n)
	se := math.Sqrt(1/(9*nf*cpk*cpk) + 1/(2*(nf-1)))

	lower := cpk - z*se
	upper := cpk + z*se

	return ConfidenceInterval{
		Lower:    lower,
		Upper:    upper,
		Level:    confidence,
		PointEst: cpk,
	}, nil
}

func CpConfidence(cp float64, n int, confidence float64) (ConfidenceInterval, error) {
	if n < 2 {
		return ConfidenceInterval{}, fmt.Errorf("need n >= 2, got %d", n)
	}
	if confidence <= 0 || confidence >= 1 {
		return ConfidenceInterval{}, fmt.Errorf("confidence must be in (0,1), got %v", confidence)
	}
	if cp <= 0 {
		return ConfidenceInterval{}, fmt.Errorf("Cp must be positive for CI computation, got %v", cp)
	}

	alpha := 1 - confidence
	df := float64(n - 1)

	chiLower := chiSquareQuantile(alpha/2, df)
	chiUpper := chiSquareQuantile(1-alpha/2, df)

	lower := cp * math.Sqrt(df/chiUpper)
	upper := cp * math.Sqrt(df/chiLower)

	return ConfidenceInterval{
		Lower:    lower,
		Upper:    upper,
		Level:    confidence,
		PointEst: cp,
	}, nil
}

func chiSquareQuantile(p float64, df float64) float64 {
	z := normalQuantile(p)
	term := 1 - 2/(9*df) + z*math.Sqrt(2/(9*df))
	result := df * term * term * term
	if result < 0 {
		result = 0
	}
	return result
}

func PPMEstimate(cpk float64) float64 {
	if cpk <= 0 {
		return 1e6
	}
	tailProb := 1 - normalCDF(3*cpk)
	return 2 * tailProb * 1e6
}

func SigmaLevel(cpk float64) float64 {
	return 3 * cpk
}

func SigmaLevelWithShift(cpk float64) float64 {
	return 3*cpk + 1.5
}
