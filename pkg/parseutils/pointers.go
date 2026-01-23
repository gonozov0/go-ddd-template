package parseutils

func ToPointer[T any](value T) *T {
	return &value
}

func FromPointer[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}

	return *value
}
