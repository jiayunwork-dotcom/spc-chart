package capability

import (
	"math"
	"testing"
)

func TestCpkConfidenceBasic(t *testing.T) {
	ci, err := CpkConfidence(1.33, 50, 0.95)
	if err != nil {
		t.Fatal(err)
	}
	if ci.Lower >= ci.PointEst || ci.Upper <= ci.PointEst {
		t.Fatalf("CI [%v, %v] should bracket point estimate %v", ci.Lower, ci.Upper, ci.PointEst)
	}
	if ci.Lower <= 0 {
		t.Fatalf("lower bound %v should be positive for Cpk=1.33", ci.Lower)
	}
}

func TestCpkConfidenceWiderWithLessData(t *testing.T) {
	ci30, _ := CpkConfidence(1.5, 30, 0.95)
	ci100, _ := CpkConfidence(1.5, 100, 0.95)

	width30 := ci30.Upper - ci30.Lower
	width100 := ci100.Upper - ci100.Lower
	if width30 <= width100 {
		t.Fatalf("CI with n=30 (%v) should be wider than n=100 (%v)", width30, width100)
	}
}

func TestCpConfidence(t *testing.T) {
	ci, err := CpConfidence(2.0, 100, 0.95)
	if err != nil {
		t.Fatal(err)
	}
	if ci.Lower >= ci.Upper {
		t.Fatalf("CI lower=%v should be < upper=%v", ci.Lower, ci.Upper)
	}
	if ci.Lower <= 0 {
		t.Fatalf("lower=%v should be positive for Cp=2.0", ci.Lower)
	}
	if ci.PointEst < ci.Lower || ci.PointEst > ci.Upper {
		t.Fatalf("point est %v not in [%v, %v]", ci.PointEst, ci.Lower, ci.Upper)
	}
}

func TestPPMEstimate(t *testing.T) {
	ppm := PPMEstimate(1.0)
	if ppm < 2000 || ppm > 3500 {
		t.Fatalf("PPM for Cpk=1 = %v, expected ~2700", ppm)
	}
	ppm2 := PPMEstimate(2.0)
	if ppm2 > 1 {
		t.Fatalf("PPM for Cpk=2 = %v, expected <1", ppm2)
	}
}

func TestSigmaLevel(t *testing.T) {
	if math.Abs(SigmaLevel(1.0)-3.0) > 1e-9 {
		t.Fatal("SigmaLevel(1) should be 3")
	}
	if math.Abs(SigmaLevel(2.0)-6.0) > 1e-9 {
		t.Fatal("SigmaLevel(2) should be 6")
	}
}

func TestSigmaLevelWithShift(t *testing.T) {
	if math.Abs(SigmaLevelWithShift(1.0)-4.5) > 1e-9 {
		t.Fatal("SigmaLevelWithShift(1) should be 4.5")
	}
}

func TestCpkConfidenceInvalid(t *testing.T) {
	_, err := CpkConfidence(-1, 50, 0.95)
	if err == nil {
		t.Fatal("expected error for negative Cpk")
	}
	_, err = CpkConfidence(1, 50, 1.5)
	if err == nil {
		t.Fatal("expected error for confidence > 1")
	}
}
