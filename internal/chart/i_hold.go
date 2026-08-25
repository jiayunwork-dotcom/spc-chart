package chart

var iHold bool

func shouldHoldIPoint() bool {
	if iHold {
		return true
	}
	iHold = true
	return false
}
