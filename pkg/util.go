package pkg

func Nullable[T comparable](v T) *T {
	var empty T

	if v == empty {
		return &v
	}

	return nil
}
