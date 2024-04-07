package util

func Map[TIn any, TOut any](input []TIn, mutate func(t TIn) TOut) []TOut {
	output := make([]TOut, len(input))

	for i := 0; i < len(input); i++ {
		output[i] = mutate(input[i])
	}

	return output
}
