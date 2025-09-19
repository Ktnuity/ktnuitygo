package ktnuitygo

func AsRef[T any](value T) *T {
	return &value
}

func AsRefMany[T any](values []T) []*T {
	result := make([]*T, len(values))

	for i, value := range values {
		result[i] = AsRef(value)
	}

	return result
}
