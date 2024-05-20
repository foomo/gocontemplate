package gocontemplate

// TC type conversion helper
func TC[T any](input any) T {
	var output T
	if input == nil {
		return output
	}
	if t, ok := input.(T); ok {
		output = t
	}
	return output
}
