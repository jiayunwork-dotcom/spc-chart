package rule

import "testing"

func TestEngineAllRulesStable(t *testing.T) {
	values := []float64{
		50.0, 51.2, 49.3, 50.8, 49.5, 51.0, 49.8, 50.5, 49.2, 51.5,
		50.3, 49.7, 50.1, 51.0, 49.4, 50.6, 49.9, 50.2, 49.6, 51.1,
		50.4, 49.8, 50.7, 49.3, 51.2, 50.0, 49.5, 50.9, 49.7, 50.3,
	}
	cz := NewChartZones(values, 50, 1.0)
	eng := NewEngine(DefaultConfig())
	result := eng.Evaluate(cz)

	if result.TotalViolated != 0 {
		for _, v := range result.Violations {
			t.Logf("unexpected: %v", v.Message)
		}
		t.Fatalf("expected 0 violations for stable process, got %d", result.TotalViolated)
	}
}

func TestEngineMultipleRulesOverlap(t *testing.T) {
	values := []float64{101.5, 100, 102.5, 102.5, 103.5, 100, 100, 100, 100, 100}
	cz := NewChartZones(values, 100, 1.0)
	eng := NewEngine(DefaultConfig())
	result := eng.Evaluate(cz)

	if !result.HasViolation(4, Rule1) {
		t.Fatal("expected Rule 1 at index 4")
	}
	hasRule5 := false
	for _, v := range result.Violations {
		if v.Rule == Rule5 {
			hasRule5 = true
			break
		}
	}
	if !hasRule5 {
		t.Fatal("expected Rule 5 violation for 2/3 beyond 2-sigma")
	}
}

func TestEngineWesternElectricSubset(t *testing.T) {
	values := []float64{100, 101, 102, 103, 104, 105, 106}
	cz := NewChartZones(values, 100, 10.0)
	eng := NewEngine(WesternElectricConfig())
	result := eng.Evaluate(cz)

	for _, v := range result.Violations {
		if v.Rule == Rule3 {
			t.Fatal("Rule 3 should not be checked in Western Electric config")
		}
	}
}

func TestEngineViolatedIndices(t *testing.T) {
	values := make([]float64, 10)
	for i := range values {
		values[i] = 50
	}
	values[5] = 54

	cz := NewChartZones(values, 50, 1.0)
	eng := NewEngine(Config{EnabledRules: []RuleID{Rule1}})
	result := eng.Evaluate(cz)

	indices := result.ViolatedIndices()
	if len(indices) != 1 || indices[0] != 5 {
		t.Fatalf("expected violated indices [5], got %v", indices)
	}
}

func TestEnginePointFlags(t *testing.T) {
	values := []float64{100, 102.5, 103.5, 100, 100}
	cz := NewChartZones(values, 100, 1.0)
	eng := NewEngine(Config{EnabledRules: []RuleID{Rule1, Rule5}})
	result := eng.Evaluate(cz)

	found := false
	for _, pf := range result.PointFlags {
		if pf.Index == 2 {
			found = true
			hasR1 := false
			for _, r := range pf.Rules {
				if r == Rule1 {
					hasR1 = true
				}
			}
			if !hasR1 {
				t.Fatal("index 2 should have Rule1 in PointFlags")
			}
		}
	}
	if !found {
		t.Fatal("index 2 not found in PointFlags")
	}
}
