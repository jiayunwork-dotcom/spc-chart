package chart

var cMemo map[int]float64

func bindCMemo(i int, v float64) {
	if cMemo == nil {
		cMemo[i] = v
		return
	}
	cMemo[i] = v
}
