package pointer

func ValueOrEmpty[T any](v *T) T {
	var empty T

	if v != nil {
		return *v
	}

	return empty
}
