package slices

// Map applies given function to every value of slice
func Map[S ~[]T, T, M any](s S, fn func(T) M) []M {
	if s == nil {
		return []M(nil)
	}

	if len(s) == 0 {
		return make([]M, 0)
	}

	res := make([]M, len(s))
	for i, v := range s {
		res[i] = fn(v)
	}

	return res
}
