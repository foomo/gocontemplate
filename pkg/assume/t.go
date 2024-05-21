package assume

// T type conversion helper
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
