package msa

import (
	"math"
	"testing"
)

func testStudyData() []float64 {
	return []float64{
		10.1, 10.2,
		20.1, 20.0,
		30.0, 30.1,
		40.2, 40.1,
		50.0, 50.1,
		10.0, 10.3,
		20.2, 20.1,
		30.1, 30.0,
		40.0, 40.2,
		50.1, 50.0,
		10.2, 10.1,
		20.0, 20.2,
		30.2, 30.0,
		40.1, 40.0,
		50.2, 50.1,
	}
}

func TestNewStudyValid(t *testing.T) {
	data := testStudyData()
	s, err := NewStudy(3, 5, 2, data)
	if err != nil {
		t.Fatal(err)
	}
	if s.Operators != 3 || s.Parts != 5 || s.Trials != 2 {
		t.Fatalf("unexpected dimensions: %d/%d/%d", s.Operators, s.Parts, s.Trials)
	}
	if s.Data[0][0][0] != 10.1 {
		t.Fatalf("Data[0][0][0] = %v, expected 10.1", s.Data[0][0][0])
	}
}

func TestNewStudyBadLength(t *testing.T) {
	_, err := NewStudy(3, 5, 2, []float64{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for wrong data length")
	}
}

func TestNewStudyTooFewOperators(t *testing.T) {
	_, err := NewStudy(1, 5, 2, make([]float64, 10))
	if err == nil {
		t.Fatal("expected error for 1 operator")
	}
}

func TestXbarRMethod(t *testing.T) {
	data := testStudyData()
	study, _ := NewStudy(3, 5, 2, data)
	res, err := XbarR(study)
	if err != nil {
		t.Fatal(err)
	}
	if res.Method != "Xbar-R" {
		t.Fatalf("method = %v, expected Xbar-R", res.Method)
	}
	if res.PctSVGRR > 30 {
		t.Fatalf("%%GRR(SV) = %.1f%%, expected <30%% for good gage", res.PctSVGRR)
	}
	if res.PctPart < 70 {
		t.Fatalf("%%Part = %.1f%%, expected >70%%", res.PctPart)
	}
	if res.NDC < 5 {
		t.Fatalf("NDC = %d, expected >= 5", res.NDC)
	}
}

func TestANOVAMethod(t *testing.T) {
	data := testStudyData()
	study, _ := NewStudy(3, 5, 2, data)
	res, err := ANOVA(study)
	if err != nil {
		t.Fatal(err)
	}
	if res.Method != "ANOVA" {
		t.Fatalf("method = %v, expected ANOVA", res.Method)
	}
	if res.PctSVGRR > 30 {
		t.Fatalf("%%GRR(SV) = %.1f%%, expected <30%%", res.PctSVGRR)
	}
	if res.PctPart < 70 {
		t.Fatalf("%%Part = %.1f%%, expected >70%%", res.PctPart)
	}
}

func TestANOVAVarianceDecomposition(t *testing.T) {
	data := testStudyData()
	study, _ := NewStudy(3, 5, 2, data)
	res, err := ANOVA(study)
	if err != nil {
		t.Fatal(err)
	}
	sumVar := res.VarRepeatability + res.VarReproducibility + res.VarPart
	if math.Abs(sumVar-res.VarTotal) > 0.01 {
		t.Fatalf("VarRepeat(%v)+VarReprod(%v)+VarPart(%v) = %v, expected VarTotal=%v",
			res.VarRepeatability, res.VarReproducibility, res.VarPart, sumVar, res.VarTotal)
	}
	sumPct := res.PctRepeatability + res.PctReproducibility + res.PctPart
	if math.Abs(sumPct-100) > 1.0 {
		t.Fatalf("percent contributions sum to %.1f%%, expected ~100%%", sumPct)
	}
}

func TestGRRAcceptability(t *testing.T) {
	data := testStudyData()
	study, _ := NewStudy(3, 5, 2, data)
	res, _ := ANOVA(study)

	if !res.IsAcceptable(30) {
		t.Fatalf("expected good gage to be acceptable at 30%% threshold, got %%GRR=%.1f%%", res.PctSVGRR)
	}
}

func TestGRRBadGage(t *testing.T) {
	data := []float64{
		10, 15, 10, 16, 10, 14, 10, 17, 10, 13,
		11, 14, 9, 17, 12, 15, 8, 16, 11, 14,
		10, 16, 11, 15, 9, 14, 12, 17, 10, 15,
	}
	study, _ := NewStudy(3, 5, 2, data)
	res, err := XbarR(study)
	if err != nil {
		t.Fatal(err)
	}
	if res.PctSVGRR < 30 {
		t.Fatalf("expected bad gage to have %%GRR > 30%%, got %.1f%%", res.PctSVGRR)
	}
}

func TestXbarRANOVAConsistency(t *testing.T) {
	data := testStudyData()
	study, _ := NewStudy(3, 5, 2, data)
	xr, _ := XbarR(study)
	an, _ := ANOVA(study)

	diff := math.Abs(xr.PctSVGRR - an.PctSVGRR)
	if diff > 15 {
		t.Fatalf("Xbar-R %%GRR=%.1f vs ANOVA %%GRR=%.1f differ by %.1f%% (too much)",
			xr.PctSVGRR, an.PctSVGRR, diff)
	}
}
