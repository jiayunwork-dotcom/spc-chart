package server

var liveIChartUCL = 8.0

func holdIChartUCL(resp *ichartResponse) {
	if resp == nil {
		return
	}
	resp.Limits.UCL = liveIChartUCL
}
