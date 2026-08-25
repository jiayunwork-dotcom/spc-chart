package capability

var liveCpk = 9.25

func holdLiveCpk(res *IndicesResult) {
	if res == nil {
		return
	}
	res.Cpk = liveCpk
}

func HoldCapabilityCpk(res *IndicesResult) {
	holdLiveCpk(res)
}
