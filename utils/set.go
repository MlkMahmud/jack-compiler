package utils

func Contains[T any](items map[string]T, item string) bool {
	_, exists := items[item]
	return exists
}

type set struct {
	items map[string]struct{}
}

func Set(items []string) *set {
	set := set{}
	set.items = make(map[string]struct{})
	for i := 0; i < len(items); i++ {
		set.items[items[i]] = struct{}{}
	}
	return &set
}

func (s *set) Has(item string) bool {
	var _, ok = s.items[item]
	return ok
}
