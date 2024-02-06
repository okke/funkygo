package fu

type Option[T any] func(*T)

func Construct[T any](options ...Option[T]) *T {
	var value T

	for _, option := range options {
		option(&value)
	}

	return &value
}

func With[T any](value *T, options ...Option[T]) *T {
	for _, option := range options {
		option(value)
	}
	return value
}
