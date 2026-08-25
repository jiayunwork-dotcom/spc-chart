package capability

import (
	"math"
	"math/rand"
	"testing"
)

func TestAndersonDarlingNormal(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	vals := make([]float64, 100)
	for i := 0; i < 100; i += 2 {
		u1 := rng.Float64()
		u2 := rng.Float64()
		z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		z1 := math.Sqrt(-2*math.Log(u1)) * math.Sin(2*math.Pi*u2)
		vals[i] = 50 + 2*z0
		if i+1 < 100 {
			vals[i+1] = 50 + 2*z1
		}
	}

	res, err := AndersonDarling(vals, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsNormal {
		t.Fatalf("expected normal data to pass AD test: A2*=%v, crit=%v", res.Statistic, res.Critical)
	}
}

func TestAndersonDarlingUniform(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	vals := make([]float64, 200)
	for i := range vals {
		vals[i] = rng.Float64() * 100
	}

	res, err := AndersonDarling(vals, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if res.IsNormal {
		t.Fatalf("expected uniform data to fail AD test: A2*=%v, crit=%v", res.Statistic, res.Critical)
	}
}

func TestAndersonDarlingInsufficientData(t *testing.T) {
	_, err := AndersonDarling([]float64{1, 2, 3}, 0.05)
	if err == nil {
		t.Fatal("expected error for n < 8")
	}
}

func TestAndersonDarlingBadAlpha(t *testing.T) {
	vals := make([]float64, 20)
	for i := range vals {
		vals[i] = float64(i)
	}
	_, err := AndersonDarling(vals, 0.03)
	if err == nil {
		t.Fatal("expected error for unsupported alpha")
	}
}

func TestShapiroWilkApproxNormal(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	vals := make([]float64, 50)
	for i := 0; i < 50; i += 2 {
		u1 := rng.Float64()
		u2 := rng.Float64()
		z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		z1 := math.Sqrt(-2*math.Log(u1)) * math.Sin(2*math.Pi*u2)
		vals[i] = z0
		if i+1 < 50 {
			vals[i+1] = z1
		}
	}

	res, err := ShapiroWilkApprox(vals)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsNormal {
		t.Fatalf("expected normal data to pass SW approx: W=%v, crit=%v", res.Statistic, res.Critical)
	}
}

func TestShapiroWilkApproxExponential(t *testing.T) {
	rng := rand.New(rand.NewSource(55))
	vals := make([]float64, 100)
	for i := range vals {
		vals[i] = -math.Log(1 - rng.Float64())
	}

	res, err := ShapiroWilkApprox(vals)
	if err != nil {
		t.Fatal(err)
	}
	if res.IsNormal {
		t.Fatalf("expected exponential data to fail SW approx: W=%v, crit=%v", res.Statistic, res.Critical)
	}
}

func TestNormalCDFBoundary(t *testing.T) {
	if math.Abs(normalCDF(0)-0.5) > 1e-10 {
		t.Fatalf("normalCDF(0) = %v, expected 0.5", normalCDF(0))
	}
	if normalCDF(-10) > 1e-10 {
		t.Fatalf("normalCDF(-10) = %v, expected ~0", normalCDF(-10))
	}
	if normalCDF(10) < 1-1e-10 {
		t.Fatalf("normalCDF(10) = %v, expected ~1", normalCDF(10))
	}
}

func TestNormalQuantileRoundTrip(t *testing.T) {
	for _, z := range []float64{-3, -1, 0, 0.5, 1.96, 2.576} {
		p := normalCDF(z)
		got := normalQuantile(p)
		if math.Abs(got-z) > 0.001 {
			t.Fatalf("normalQuantile(normalCDF(%v)) = %v, expected %v", z, got, z)
		}
	}
}
