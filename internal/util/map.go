package util

func Map[TIn any, TOut any](input []TIn, mutate func(t TIn) TOut) []TOut {
	output := make([]TOut, len(input))

	for i, inputItem := range input {
		output[i] = mutate(inputItem)
	}

	return output
}
