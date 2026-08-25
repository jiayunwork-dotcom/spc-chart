package chart

import "fmt"

type ChartType int

const (
	TypeIndividuals ChartType = iota
	TypeXbarR
	TypeXbarS
	TypeMR
	TypeCUSUM
	TypeEWMA
)

func (ct ChartType) String() string {
	switch ct {
	case TypeIndividuals:
		return "I-Chart"
	case TypeXbarR:
		return "Xbar-R"
	case TypeXbarS:
		return "Xbar-S"
	case TypeMR:
		return "MR"
	case TypeCUSUM:
		return "CUSUM"
	case TypeEWMA:
		return "EWMA"
	default:
		return fmt.Sprintf("Unknown(%d)", int(ct))
	}
}

type ControlLimit struct {
	CL  float64
	UCL float64
	LCL float64
}

type PlotPoint struct {
	Index     int
	Value     float64
	OutOfCtrl bool
	Subgroup  int
}

type ChartResult struct {
	Type     ChartType
	Limits   ControlLimit
	Points   []PlotPoint
	OOCCount int
	Mean     float64
}

type Subgroup struct {
	Values []float64
}

func ValidateSubgroups(subs []Subgroup, minSize int) error {
	if len(subs) == 0 {
		return fmt.Errorf("no subgroups provided")
	}
	for i, sg := range subs {
		if len(sg.Values) < minSize {
			return fmt.Errorf("subgroup %d has %d values, need at least %d", i, len(sg.Values), minSize)
		}
	}
	return nil
}
