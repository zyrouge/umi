package utils

func SliceToSet[T comparable](list []T) map[T]struct{} {
	set := make(map[T]struct{})
	for _, item := range list {
		set[item] = struct{}{}
	}
	return set
}

func SliceToAny[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
