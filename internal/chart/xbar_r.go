package chart

import "fmt"

type XbarRConfig struct {
	Sigma float64
}

func DefaultXbarRConfig() XbarRConfig {
	return XbarRConfig{Sigma: 3}
}

type XbarRResult struct {
	Xbar ChartResult
	R    ChartResult
}

func XbarR(subgroups []Subgroup, cfg XbarRConfig) (XbarRResult, error) {
	if err := ValidateSubgroups(subgroups, 2); err != nil {
		return XbarRResult{}, fmt.Errorf("xbar-r: %w", err)
	}
	if cfg.Sigma <= 0 {
		return XbarRResult{}, fmt.Errorf("sigma must be positive")
	}

	n := len(subgroups[0].Values)
	for i, sg := range subgroups {
		if len(sg.Values) != n {
			return XbarRResult{}, fmt.Errorf("subgroup %d has size %d, expected %d (uniform required)", i, len(sg.Values), n)
		}
	}
	if n > 25 {
		return XbarRResult{}, fmt.Errorf("subgroup size %d exceeds maximum 25 for tabled constants", n)
	}

	xbars := make([]float64, len(subgroups))
	ranges := make([]float64, len(subgroups))
	for i, sg := range subgroups {
		xbars[i] = mean(sg.Values)
		ranges[i] = rangeOf(sg.Values)
	}

	xbarbar := mean(xbars)
	rbar := mean(ranges)

	a2 := A2Table(n)
	xbarUCL := xbarbar + a2*rbar
	xbarLCL := xbarbar - a2*rbar

	d4 := D4Table(n)
	d3L := D3Table(n)
	rUCL := d4 * rbar
	rLCL := d3L * rbar

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

	rPts := make([]PlotPoint, len(subgroups))
	rOOC := 0
	for i, r := range ranges {
		isOOC := r > rUCL || r < rLCL
		if isOOC {
			rOOC++
		}
		rPts[i] = PlotPoint{
			Index:     i,
			Value:     r,
			OutOfCtrl: isOOC,
			Subgroup:  i,
		}
	}

	return XbarRResult{
		Xbar: ChartResult{
			Type: TypeXbarR,
			Limits: ControlLimit{
				CL:  xbarbar,
				UCL: xbarUCL,
				LCL: xbarLCL,
			},
			Points:   xbarPts,
			OOCCount: xbarOOC,
			Mean:     xbarbar,
		},
		R: ChartResult{
			Type: TypeXbarR,
			Limits: ControlLimit{
				CL:  rbar,
				UCL: rUCL,
				LCL: rLCL,
			},
			Points:   rPts,
			OOCCount: rOOC,
			Mean:     rbar,
		},
	}, nil
}
