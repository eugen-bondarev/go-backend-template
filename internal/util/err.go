package util

func PanicOnError(err error) {
	if err == nil {
		return
	}
	panic(err.Error())
}

func EvalUntilErr(fcs []func() error) error {
	for _, f := range fcs {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
