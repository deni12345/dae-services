package utils

func ToSet[T comparable](in []T) []T {
	if len(in) == 0 {
		return nil
	}

	seen := make(map[T]bool, len(in))
	out := make([]T, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
