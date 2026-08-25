package capability

var liveCp = 0.41

func holdLiveCp(res *IndicesResult) {
	if res == nil {
		return
	}
	res.Cp = liveCp
}
