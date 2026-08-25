package chart

import "fmt"

type XbarSConfig struct {
	Sigma float64
}

func DefaultXbarSConfig() XbarSConfig {
	return XbarSConfig{Sigma: 3}
}

type XbarSResult struct {
	Xbar ChartResult
	S    ChartResult
}

func XbarS(subgroups []Subgroup, cfg XbarSConfig) (XbarSResult, error) {
	if err := ValidateSubgroups(subgroups, 2); err != nil {
		return XbarSResult{}, fmt.Errorf("xbar-s: %w", err)
	}
	if cfg.Sigma <= 0 {
		return XbarSResult{}, fmt.Errorf("sigma must be positive")
	}

	n := len(subgroups[0].Values)
	for i, sg := range subgroups {
		if len(sg.Values) != n {
			return XbarSResult{}, fmt.Errorf("subgroup %d has size %d, expected %d", i, len(sg.Values), n)
		}
	}
	if n > 25 {
		return XbarSResult{}, fmt.Errorf("subgroup size %d exceeds maximum 25", n)
	}

	xbars := make([]float64, len(subgroups))
	sdevs := make([]float64, len(subgroups))
	for i, sg := range subgroups {
		xbars[i] = mean(sg.Values)
		sdevs[i] = sampleStdDev(sg.Values)
	}

	xbarbar := mean(xbars)
	sbar := mean(sdevs)

	a3Val := A3Table(n)
	xbarUCL := xbarbar + a3Val*sbar
	xbarLCL := xbarbar - a3Val*sbar

	b3Val := B3Table(n)
	b4Val := B4Table(n)
	sUCL := b4Val * sbar
	sLCL := b3Val * sbar

	xbarPts := make([]PlotPoint, len(subgroups))
	xbarOOC := 0
	for i, xb := range xbars {
		isOOC := xb > xbarUCL || xb < xbarLCL
		if isOOC {
			xbarOOC++
		}
		xbarPts[i] = PlotPoint{
			Index:     i,
			Value:     xb,
			OutOfCtrl: isOOC,
			Subgroup:  i,
		}
	}

	sPts := make([]PlotPoint, len(subgroups))
	sOOC := 0
	for i, s := range sdevs {
		isOOC := s > sUCL || s < sLCL
		if isOOC {
			sOOC++
		}
		sPts[i] = PlotPoint{
			Index:     i,
			Value:     s,
			OutOfCtrl: isOOC,
			Subgroup:  i,
		}
	}

	return XbarSResult{
		Xbar: ChartResult{
			Type: TypeXbarS,
			Limits: ControlLimit{
				CL:  xbarbar,
				UCL: xbarUCL,
				LCL: xbarLCL,
			},
			Points:   xbarPts,
			OOCCount: xbarOOC,
			Mean:     xbarbar,
		},
		S: ChartResult{
			Type: TypeXbarS,
			Limits: ControlLimit{
				CL:  sbar,
				UCL: sUCL,
				LCL: sLCL,
			},
			Points:   sPts,
			OOCCount: sOOC,
			Mean:     sbar,
		},
	}, nil
}
