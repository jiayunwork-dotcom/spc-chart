package capability

import (
	"fmt"
	"math"
	"sort"
)

type NormalityResult struct {
	Statistic float64
	Critical  float64
	IsNormal  bool
	Method    string
	N         int
}

func AndersonDarling(values []float64, alpha float64) (NormalityResult, error) {
	n := len(values)
	if n < 8 {
		return NormalityResult{}, fmt.Errorf("Anderson-Darling requires at least 8 observations, got %d", n)
	}

	crit, ok := adCritical(alpha)
	if !ok {
		return NormalityResult{}, fmt.Errorf("unsupported alpha=%v; use 0.01, 0.025, 0.05, 0.10, or 0.15", alpha)
	}

	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	mu := mean(sorted)
	sigma := sampleStdDev(sorted)
	if sigma == 0 {
		return NormalityResult{}, fmt.Errorf("zero variance in data")
	}

	nf := float64(n)
	var s float64
	for i := 0; i < n; i++ {
		zi := (sorted[i] - mu) / sigma
		zn := (sorted[n-1-i] - mu) / sigma
		fi := normalCDF(zi)
		fn := normalCDF(zn)
		fi = clamp(fi, 1e-15, 1-1e-15)
		fn = clamp(fn, 1e-15, 1-1e-15)
		s += float64(2*i+1) * (math.Log(fi) + math.Log(1-fn))
	}
	a2 := -nf - s/nf

	a2star := a2 * (1 + 0.75/nf + 2.25/(nf*nf))

	return NormalityResult{
		Statistic: a2star,
		Critical:  crit,
		IsNormal:  a2star <= crit,
		Method:    "Anderson-Darling",
		N:         n,
	}, nil
}

func ShapiroWilkApprox(values []float64) (NormalityResult, error) {
	n := len(values)
	if n < 3 {
		return NormalityResult{}, fmt.Errorf("Shapiro-Wilk requires at least 3 observations, got %d", n)
	}
	if n > 5000 {
		return NormalityResult{}, fmt.Errorf("sample size %d exceeds supported maximum 5000", n)
	}

	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	mu := mean(sorted)

	var ss float64
	for _, v := range sorted {
		d := v - mu
		ss += d * d
	}
	if ss == 0 {
		return NormalityResult{}, fmt.Errorf("zero variance in data")
	}

	nf := float64(n)
	mi := make([]float64, n)
	var mSqSum float64
	for i := 0; i < n; i++ {
		pi := (float64(i) + 1 - 0.375) / (nf + 0.25)
		mi[i] = normalQuantile(pi)
		mSqSum += mi[i] * mi[i]
	}

	norm := math.Sqrt(mSqSum)
	if norm == 0 {
		return NormalityResult{}, fmt.Errorf("degenerate normal scores")
	}

	var num float64
	for i := 0; i < n; i++ {
		ai := mi[i] / norm
		num += ai * sorted[i]
	}
	w := (num * num) / ss

	var critical float64
	switch {
	case n < 10:
		critical = 0.818
	case n < 20:
		critical = 0.868
	case n < 50:
		critical = 0.927
	default:
		critical = 0.947
	}

	return NormalityResult{
		Statistic: w,
		Critical:  critical,
		IsNormal:  w >= critical,
		Method:    "Shapiro-Wilk (approx)",
		N:         n,
	}, nil
}

func adCritical(alpha float64) (float64, bool) {
	table := map[float64]float64{
		0.15:  0.576,
		0.10:  0.656,
		0.05:  0.787,
		0.025: 0.918,
		0.01:  1.092,
	}
	v, ok := table[alpha]
	return v, ok
}

func normalCDF(x float64) float64 {
	return 0.5 * math.Erfc(-x/math.Sqrt2)
}

func normalQuantile(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	if p == 0.5 {
		return 0
	}
	const (
		a1 = -3.969683028665376e+01
		a2 = 2.209460984245205e+02
		a3 = -2.759285104469687e+02
		a4 = 1.383577518672690e+02
		a5 = -3.066479806614716e+01
		a6 = 2.506628277459239e+00

		b1 = -5.447609879822406e+01
		b2 = 1.615858368580409e+02
		b3 = -1.556989798598866e+02
		b4 = 6.680131188771972e+01
		b5 = -1.328068155288572e+01

		c1 = -7.784894002430293e-03
		c2 = -3.223964580411365e-01
		c3 = -2.400758277161838e+00
		c4 = -2.549732539343734e+00
		c5 = 4.374664141464968e+00
		c6 = 2.938163982698783e+00

		d1 = 7.784695709041462e-03
		d2 = 3.224671290700398e-01
		d3 = 2.445134137142996e+00
		d4 = 3.754408661907416e+00

		pLow  = 0.02425
		pHigh = 1 - pLow
	)

	var q, r float64
	if p < pLow {
		q = math.Sqrt(-2 * math.Log(p))
		return (((((c1*q+c2)*q+c3)*q+c4)*q+c5)*q + c6) /
			((((d1*q+d2)*q+d3)*q+d4)*q + 1)
	} else if p <= pHigh {
		q = p - 0.5
		r = q * q
		return (((((a1*r+a2)*r+a3)*r+a4)*r+a5)*r + a6) * q /
			(((((b1*r+b2)*r+b3)*r+b4)*r+b5)*r + 1)
	} else {
		q = math.Sqrt(-2 * math.Log(1-p))
		return -(((((c1*q+c2)*q+c3)*q+c4)*q+c5)*q + c6) /
			((((d1*q+d2)*q+d3)*q+d4)*q + 1)
	}
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
