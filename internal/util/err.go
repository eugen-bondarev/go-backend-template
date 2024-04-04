package util

func PanicOnError(err error) {
	if err == nil {
		return
	}
	panic(err.Error())
}
