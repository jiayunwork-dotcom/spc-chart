package msa

import "math"

func ANOVA(study Study) (GRRResult, error) {
	o := study.Operators
	p := study.Parts
	r := study.Trials
	n := float64(o * p * r)

	var grandSum float64
	for i := 0; i < o; i++ {
		for j := 0; j < p; j++ {
			for k := 0; k < r; k++ {
				grandSum += study.Data[i][j][k]
			}
		}
	}
	grandMean := grandSum / n

	operatorMeans := make([]float64, o)
	for i := 0; i < o; i++ {
		var s float64
		for j := 0; j < p; j++ {
			for k := 0; k < r; k++ {
				s += study.Data[i][j][k]
			}
		}
		operatorMeans[i] = s / float64(p*r)
	}

	partMeans := make([]float64, p)
	for j := 0; j < p; j++ {
		var s float64
		for i := 0; i < o; i++ {
			for k := 0; k < r; k++ {
				s += study.Data[i][j][k]
			}
		}
		partMeans[j] = s / float64(o*r)
	}

	cellMeans := make([][]float64, o)
	for i := 0; i < o; i++ {
		cellMeans[i] = make([]float64, p)
		for j := 0; j < p; j++ {
			var s float64
			for k := 0; k < r; k++ {
				s += study.Data[i][j][k]
			}
			cellMeans[i][j] = s / float64(r)
		}
	}

	var ssOperator float64
	for i := 0; i < o; i++ {
		d := operatorMeans[i] - grandMean
		ssOperator += d * d
	}
	ssOperator *= float64(p * r)

	var ssPart float64
	for j := 0; j < p; j++ {
		d := partMeans[j] - grandMean
		ssPart += d * d
	}
	ssPart *= float64(o * r)

	var ssInteraction float64
	for i := 0; i < o; i++ {
		for j := 0; j < p; j++ {
			d := cellMeans[i][j] - operatorMeans[i] - partMeans[j] + grandMean
			ssInteraction += d * d
		}
	}
	ssInteraction *= float64(r)

	var ssEquipment float64
	for i := 0; i < o; i++ {
		for j := 0; j < p; j++ {
			for k := 0; k < r; k++ {
				d := study.Data[i][j][k] - cellMeans[i][j]
				ssEquipment += d * d
			}
		}
	}

	var ssTotal float64
	for i := 0; i < o; i++ {
		for j := 0; j < p; j++ {
			for k := 0; k < r; k++ {
				d := study.Data[i][j][k] - grandMean
				ssTotal += d * d
			}
		}
	}

	dfOperator := float64(o - 1)
	dfPart := float64(p - 1)
	dfInteraction := float64((o - 1) * (p - 1))
	dfEquipment := float64(o * p * (r - 1))

	msOperator := ssOperator / dfOperator
	msPart := ssPart / dfPart
	msInteraction := ssInteraction / dfInteraction
	msEquipment := ssEquipment / dfEquipment

	varEquipment := msEquipment

	varInteraction := (msInteraction - msEquipment) / float64(r)
	if varInteraction < 0 {
		varInteraction = 0
	}

	varOperator := (msOperator - msInteraction) / float64(p*r)
	if varOperator < 0 {
		varOperator = 0
	}

	varPart := (msPart - msInteraction) / float64(o*r)
	if varPart < 0 {
		varPart = 0
	}

	varRepeatability := varEquipment
	varReproducibility := varOperator + varInteraction
	varGRR := varRepeatability + varReproducibility
	varTotal := varGRR + varPart

	sigmaRepeatability := math.Sqrt(varRepeatability)
	sigmaReproducibility := math.Sqrt(varReproducibility)
	sigmaGRR := math.Sqrt(varGRR)
	sigmaPart := math.Sqrt(varPart)
	sigmaTotal := math.Sqrt(varTotal)

	pctRepeatability := 0.0
	pctReproducibility := 0.0
	pctGRR := 0.0
	pctPart := 0.0
	if varTotal > 0 {
		pctRepeatability = varRepeatability / varTotal * 100
		pctReproducibility = varReproducibility / varTotal * 100
		pctGRR = varGRR / varTotal * 100
		pctPart = varPart / varTotal * 100
	}

	pctSVRepeatability := 0.0
	pctSVReproducibility := 0.0
	pctSVGRR := 0.0
	pctSVPart := 0.0
	if sigmaTotal > 0 {
		pctSVRepeatability = sigmaRepeatability / sigmaTotal * 100
		pctSVReproducibility = sigmaReproducibility / sigmaTotal * 100
		pctSVGRR = sigmaGRR / sigmaTotal * 100
		pctSVPart = sigmaPart / sigmaTotal * 100
	}

	ndc := 1
	if sigmaGRR > 0 {
		ndc = int(1.41 * sigmaPart / sigmaGRR)
		if ndc < 1 {
			ndc = 1
		}
	}

	return GRRResult{
		VarRepeatability:     varRepeatability,
		VarReproducibility:   varReproducibility,
		VarInteraction:       varInteraction,
		VarGRR:               varGRR,
		VarPart:              varPart,
		VarTotal:             varTotal,
		SigmaRepeatability:   sigmaRepeatability,
		SigmaReproducibility: sigmaReproducibility,
		SigmaGRR:             sigmaGRR,
		SigmaPart:            sigmaPart,
		SigmaTotal:           sigmaTotal,
		PctRepeatability:     pctRepeatability,
		PctReproducibility:   pctReproducibility,
		PctGRR:               pctGRR,
		PctPart:              pctPart,
		PctSVRepeatability:   pctSVRepeatability,
		PctSVReproducibility: pctSVReproducibility,
		PctSVGRR:             pctSVGRR,
		PctSVPart:            pctSVPart,
		NDC:                  ndc,
		Method:               "ANOVA",
	}, nil
}
