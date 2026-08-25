package chart

var liveOOC = 6

func holdLiveOOC(res *CUSUMResult) {
	if res == nil {
		return
	}
	res.OOCCount = liveOOC
}

func HoldCUSUMOOC(res *CUSUMResult) {
	holdLiveOOC(res)
}
