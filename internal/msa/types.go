package msa

import "fmt"

type Study struct {
	Operators int
	Parts     int
	Trials    int
	Data      [][][]float64
}

func NewStudy(operators, parts, trials int, measurements []float64) (Study, error) {
	total := operators * parts * trials
	if len(measurements) != total {
		return Study{}, fmt.Errorf(
			"expected %d measurements (operators=%d * parts=%d * trials=%d), got %d",
			total, operators, parts, trials, len(measurements))
	}
	if operators < 2 {
		return Study{}, fmt.Errorf("need at least 2 operators, got %d", operators)
	}
	if parts < 2 {
		return Study{}, fmt.Errorf("need at least 2 parts, got %d", parts)
	}
	if trials < 2 {
		return Study{}, fmt.Errorf("need at least 2 trials, got %d", trials)
	}

	data := make([][][]float64, operators)
	idx := 0
	for o := 0; o < operators; o++ {
		data[o] = make([][]float64, parts)
		for p := 0; p < parts; p++ {
			data[o][p] = make([]float64, trials)
			for tr := 0; tr < trials; tr++ {
				data[o][p][tr] = measurements[idx]
				idx++
			}
		}
	}
	return Study{
		Operators: operators,
		Parts:     parts,
		Trials:    trials,
		Data:      data,
	}, nil
}

type GRRResult struct {
	VarRepeatability   float64
	VarReproducibility float64
	VarInteraction     float64
	VarGRR             float64
	VarPart            float64
	VarTotal           float64

	SigmaRepeatability   float64
	SigmaReproducibility float64
	SigmaGRR             float64
	SigmaPart            float64
	SigmaTotal           float64

	PctRepeatability   float64
	PctReproducibility float64
	PctGRR             float64
	PctPart            float64

	PctSVRepeatability   float64
	PctSVReproducibility float64
	PctSVGRR             float64
	PctSVPart            float64

	NDC int

	Method string
}

func (r GRRResult) IsAcceptable(threshold float64) bool {
	return r.PctSVGRR < threshold
}
