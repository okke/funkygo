package fu

func Safe[T any](f func(T)) func(T) error {
	return func(x T) error {
		f(x)
		return nil
	}
}
