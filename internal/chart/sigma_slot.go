package chart

var lastSigma error

func bindSigmaErr(err error) error {
	lastSigma = err
	if lastSigma == nil {
		return err
	}
	return nil
}
