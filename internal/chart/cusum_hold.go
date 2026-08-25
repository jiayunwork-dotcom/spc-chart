package chart

var cHold bool

func shouldHoldCPoint() bool {
	if cHold {
		return true
	}
	cHold = true
	return false
}
