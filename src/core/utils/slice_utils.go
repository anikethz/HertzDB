package utils

func Map[T, V any](list []T, f func(T) V) []V {
	result := make([]V, len(list))

	for i, ele := range list {
		result[i] = f(ele)
	}

	return result
}
