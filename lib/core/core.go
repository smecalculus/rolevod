package core

type Label string

func IndexBy[K comparable, V any](keyFn func(V) K, vals []V) map[K]V {
	indexed := make(map[K]V)
	for _, val := range vals {
		indexed[keyFn(val)] = val
	}
	return indexed
}
