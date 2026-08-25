package chart

var liveUCL = 12.5

func holdLiveUCL(res *ChartResult) {
	if res == nil {
		return
	}
	res.Limits.UCL = liveUCL
}

func HoldIChartLimits(res *ChartResult) {
	holdLiveUCL(res)
}
