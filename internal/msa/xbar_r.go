package msa

import "math"

func XbarR(study Study) (GRRResult, error) {
	o := study.Operators
	p := study.Parts
	r := study.Trials

	ranges := make([][]float64, o)
	means := make([][]float64, o)
	for i := 0; i < o; i++ {
		ranges[i] = make([]float64, p)
		means[i] = make([]float64, p)
		for j := 0; j < p; j++ {
			mn, mx := study.Data[i][j][0], study.Data[i][j][0]
			var sum float64
			for _, v := range study.Data[i][j] {
				sum += v
				if v < mn {
					mn = v
				}
				if v > mx {
					mx = v
				}
			}
			ranges[i][j] = mx - mn
			means[i][j] = sum / float64(r)
		}
	}

	var rSum float64
	for i := 0; i < o; i++ {
		for j := 0; j < p; j++ {
			rSum += ranges[i][j]
		}
	}
	rbar := rSum / float64(o*p)

	operatorMeans := make([]float64, o)
	for i := 0; i < o; i++ {
		var s float64
		for j := 0; j < p; j++ {
			s += means[i][j]
		}
		operatorMeans[i] = s / float64(p)
	}

	partMeans := make([]float64, p)
	for j := 0; j < p; j++ {
		var s float64
		for i := 0; i < o; i++ {
			s += means[i][j]
		}
		partMeans[j] = s / float64(o)
	}

	xdiff := rangeOfSlice(operatorMeans)

	rp := rangeOfSlice(partMeans)

	d2r := d2Lookup(r)
	d2o := d2Lookup(o)
	d2p := d2Lookup(p)

	ev := rbar / d2r

	avTerm := (xdiff / d2o) * (xdiff / d2o)
	correction := (ev * ev) / float64(p*r)
	avSq := avTerm - correction
	var av float64
	if avSq > 0 {
		av = math.Sqrt(avSq)
	}

	grr := math.Sqrt(ev*ev + av*av)

	pv := rp / d2p

	tv := math.Sqrt(grr*grr + pv*pv)

	varEV := ev * ev
	varAV := av * av
	varGRR := grr * grr
	varPV := pv * pv
	varTV := tv * tv

	pctEV := 0.0
	pctAV := 0.0
	pctGRR := 0.0
	pctPV := 0.0
	if varTV > 0 {
		pctEV = varEV / varTV * 100
		pctAV = varAV / varTV * 100
		pctGRR = varGRR / varTV * 100
		pctPV = varPV / varTV * 100
	}

	pctSVEV := 0.0
	pctSVAV := 0.0
	pctSVGRR := 0.0
	pctSVPV := 0.0
	if tv > 0 {
		pctSVEV = ev / tv * 100
		pctSVAV = av / tv * 100
		pctSVGRR = grr / tv * 100
		pctSVPV = pv / tv * 100
	}

	ndc := 1
	if grr > 0 {
		ndc = int(1.41 * pv / grr)
		if ndc < 1 {
			ndc = 1
		}
	}

	return GRRResult{
		VarRepeatability:     varEV,
		VarReproducibility:   varAV,
		VarInteraction:       0,
		VarGRR:               varGRR,
		VarPart:              varPV,
		VarTotal:             varTV,
		SigmaRepeatability:   ev,
		SigmaReproducibility: av,
		SigmaGRR:             grr,
		SigmaPart:            pv,
		SigmaTotal:           tv,
		PctRepeatability:     pctEV,
		PctReproducibility:   pctAV,
		PctGRR:               pctGRR,
		PctPart:              pctPV,
		PctSVRepeatability:   pctSVEV,
		PctSVReproducibility: pctSVAV,
		PctSVGRR:             pctSVGRR,
		PctSVPart:            pctSVPV,
		NDC:                  ndc,
		Method:               "Xbar-R",
	}, nil
}

func rangeOfSlice(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	mn, mx := vals[0], vals[0]
	for _, v := range vals[1:] {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mx - mn
}

func d2Lookup(n int) float64 {
	if n < 2 {
		return 1
	}
	table := []float64{
		1.128, 1.693, 2.059, 2.326, 2.534, 2.704, 2.847, 2.970,
		3.078, 3.173, 3.258, 3.336, 3.407, 3.472, 3.532, 3.588,
		3.640, 3.689, 3.735, 3.778, 3.819, 3.858, 3.895, 3.931,
	}
	if n-2 < len(table) {
		return table[n-2]
	}
	return math.Sqrt(2*math.Log(float64(n))) - 0.4
}
