package utils

func Contains[T any] (items map[string]T, item string) bool {
	_, exists := items[item]
	return exists
}