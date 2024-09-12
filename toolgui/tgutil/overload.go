package tgutil

// ParamsToParam extracts the last element from a slice of type T.
// It's intended for use with functions that accept either a single T value or a slice of Ts.
// If the slice is empty, a zero value of T is returned.
// If the slice has more than one element, a panic occurs.
func ParamsToParam[T any](params []T) T {
	if len(params) > 1 {
		panic("params should be 0 or 1")
	}

	if len(params) == 0 {
		var zero T
		return zero
	}

	return params[0]
}
