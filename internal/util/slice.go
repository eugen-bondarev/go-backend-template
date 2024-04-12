package util

func Remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func Diff[T comparable](src []T, arrays ...[]T) []T {
	elemMap := make(map[T]struct{})
	for _, array := range arrays {
		for i := range array {
			elemMap[array[i]] = struct{}{}
		}
	}
	result := make([]T, 0, len(src))
	for i := range src {
		if _, ok := elemMap[src[i]]; !ok {
			result = append(result, src[i])
		}
	}
	return result
}
