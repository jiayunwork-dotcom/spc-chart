package chart

import "testing"

func TestCUSUMNoShift(t *testing.T) {
	vals := make([]float64, 50)
	for i := range vals {
		vals[i] = 100
	}
	cfg := CUSUMConfig{
		Target: 100,
		K:      0.5,
		H:      5,
		Sigma:  1,
	}
	res, err := CUSUM(vals, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if res.OOCCount != 0 {
		t.Fatalf("expected 0 alarms for stable process, got %d", res.OOCCount)
	}
	for i, p := range res.Points {
		if p.CPlus != 0 || p.CMinus != 0 {
			t.Fatalf("point %d: CPlus=%v CMinus=%v, expected 0", i, p.CPlus, p.CMinus)
		}
	}
}

func TestCUSUMUpwardShift(t *testing.T) {
	vals := make([]float64, 40)
	for i := 0; i < 20; i++ {
		vals[i] = 50
	}
	for i := 20; i < 40; i++ {
		vals[i] = 52
	}
	cfg := CUSUMConfig{
		Target: 50,
		K:      0.5,
		H:      5,
		Sigma:  1,
	}
	res, err := CUSUM(vals, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if res.OOCCount == 0 {
		t.Fatal("expected alarms after upward shift")
	}
	for _, p := range res.Points {
		if p.OutOfCtrl {
			if p.Side != CUSUMUpper && p.Side != CUSUMBoth {
				t.Fatalf("first alarm side = %v, expected Upper", p.Side)
			}
			break
		}
	}
}

func TestCUSUMResetBehavior(t *testing.T) {
	vals := make([]float64, 30)
	for i := 0; i < 10; i++ {
		vals[i] = 50
	}
	for i := 10; i < 20; i++ {
		vals[i] = 53
	}
	for i := 20; i < 30; i++ {
		vals[i] = 50
	}
	cfg := CUSUMConfig{
		Target: 50,
		K:      0.5,
		H:      4,
		Sigma:  1,
	}
	res, err := CUSUMResetOnAlarm(vals, cfg)
	if err != nil {
		t.Fatal(err)
	}
	postAlarmOOC := 0
	for _, p := range res.Points[20:] {
		if p.OutOfCtrl {
			postAlarmOOC++
		}
	}
	if postAlarmOOC > 0 {
		t.Fatalf("expected 0 alarms after reset and return to target, got %d", postAlarmOOC)
	}
}

func TestCUSUMInvalidParams(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 5}
	_, err := CUSUM(vals, CUSUMConfig{K: -1, H: 5, Sigma: 1, Target: 3})
	if err == nil {
		t.Fatal("expected error for negative K")
	}
	_, err = CUSUM(vals, CUSUMConfig{K: 0.5, H: 0, Sigma: 1, Target: 3})
	if err == nil {
		t.Fatal("expected error for zero H")
	}
}
