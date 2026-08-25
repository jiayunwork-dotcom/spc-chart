package chart

import "fmt"

type MRConfig struct {
	Sigma float64
	Span  int
}

func DefaultMRConfig() MRConfig {
	return MRConfig{Sigma: 3, Span: 2}
}

func MovingRange(values []float64, cfg MRConfig) (ChartResult, error) {
	if len(values) < 2 {
		return ChartResult{}, fmt.Errorf("MR chart requires at least 2 values, got %d", len(values))
	}
	if cfg.Sigma <= 0 {
		return ChartResult{}, fmt.Errorf("sigma must be positive, got %v", cfg.Sigma)
	}
	if cfg.Span < 2 {
		cfg.Span = 2
	}

	mrs := movingRanges(values, cfg.Span)
	if len(mrs) == 0 {
		return ChartResult{}, fmt.Errorf("no moving ranges computed")
	}

	mrbar := mean(mrs)

	d4 := D4Table(cfg.Span)
	d3L := D3Table(cfg.Span)
	if d4 == 0 {
		d4 = 3.267
	}

	ucl := d4 * mrbar
	lcl := d3L * mrbar

	points := make([]PlotPoint, len(mrs))
	ooc := 0
	for i, mr := range mrs {
		isOOC := mr > ucl || mr < lcl
		if isOOC {
			ooc++
		}
		points[i] = PlotPoint{
			Index:     i + cfg.Span - 1,
			Value:     mr,
			OutOfCtrl: isOOC,
			Subgroup:  -1,
		}
	}

	return ChartResult{
		Type: TypeMR,
		Limits: ControlLimit{
			CL:  mrbar,
			UCL: ucl,
			LCL: lcl,
		},
		Points:   points,
		OOCCount: ooc,
		Mean:     mrbar,
	}, nil
}
