package rule

import "testing"

func TestRule1Beyond3Sigma(t *testing.T) {
	values := []float64{100, 100, 100, 100, 100, 103.5, 100, 100}
	cz := NewChartZones(values, 100, 1.0)

	vs := CheckRule1(cz)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Index != 5 {
		t.Fatalf("violation at index %d, expected 5", vs[0].Index)
	}
}

func TestRule1NegativeSide(t *testing.T) {
	values := []float64{100, 100, 96, 100, 100}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule1(cz)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation (below -3sigma), got %d", len(vs))
	}
	if vs[0].Index != 2 {
		t.Fatalf("violation at index %d, expected 2", vs[0].Index)
	}
}

func TestRule2NineConsecutiveSameSide(t *testing.T) {
	values := make([]float64, 12)
	for i := range values {
		values[i] = 100
	}
	for i := 2; i <= 10; i++ {
		values[i] = 101.5
	}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule2(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 2 violation for 9 consecutive above CL")
	}
}

func TestRule2NoViolation(t *testing.T) {
	values := make([]float64, 20)
	for i := range values {
		if i%2 == 0 {
			values[i] = 101
		} else {
			values[i] = 99
		}
	}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule2(cz)
	if len(vs) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(vs))
	}
}

func TestRule3SixIncreasing(t *testing.T) {
	values := []float64{100, 100, 101, 102, 103, 104, 105, 106, 100}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule3(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 3 violation for 6 increasing points")
	}
}

func TestRule3SixDecreasing(t *testing.T) {
	values := []float64{106, 105, 104, 103, 102, 101, 100, 100}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule3(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 3 violation for 6 decreasing points")
	}
}

func TestRule4FourteenAlternating(t *testing.T) {
	values := make([]float64, 16)
	for i := range values {
		if i%2 == 0 {
			values[i] = 101
		} else {
			values[i] = 99
		}
	}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule4(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 4 violation for 14 alternating points")
	}
}

func TestRule5TwoOfThreeBeyond2Sigma(t *testing.T) {
	values := []float64{102.5, 100, 102.5, 100, 100}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule5(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 5 violation (2/3 beyond +2-sigma)")
	}
}

func TestRule6FourOfFiveBeyond1Sigma(t *testing.T) {
	values := []float64{101.5, 101.5, 100.5, 101.5, 101.5, 100}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule6(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 6 violation (4/5 beyond +1-sigma)")
	}
}

func TestRule7FifteenWithin1Sigma(t *testing.T) {
	values := make([]float64, 18)
	for i := range values {
		values[i] = 100.3
	}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule7(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 7 violation (15 within 1-sigma)")
	}
}

func TestRule8EightBeyond1SigmaBothSides(t *testing.T) {
	values := []float64{101.5, 98.5, 102, 97, 101.5, 98, 102.5, 97.5, 100}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule8(cz)
	if len(vs) == 0 {
		t.Fatal("expected Rule 8 violation (8 beyond 1-sigma both sides)")
	}
}

func TestRule8NotTriggeredSingleSide(t *testing.T) {
	values := []float64{101.5, 102, 101.8, 103, 101.5, 102, 101.8, 103}
	cz := NewChartZones(values, 100, 1.0)
	vs := CheckRule8(cz)
	if len(vs) != 0 {
		t.Fatal("Rule 8 should not fire when all points are on one side")
	}
}
