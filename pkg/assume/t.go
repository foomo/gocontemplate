package assume

// T returns input asserted to type T. It returns the zero value of T if
// input is nil or does not hold a value of type T, and never panics.
func T[T any](input any) T {
	var output T
	if input == nil {
		return output
	}

	if t, ok := input.(T); ok {
		output = t
	}

	return output
}
