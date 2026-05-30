package utils

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}
