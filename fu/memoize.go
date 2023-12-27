package fu

func MemoizeWithError[T any](f func() (T, error)) func() (T, error) {

	memoized := false
	var result T
	var err error
	return func() (T, error) {
		if memoized {
			return result, err
		}
		result, err = f()
		memoized = true
		return result, err
	}
}

func Memoize[T any](f func() T) func() T {

	memoized := false
	var result T
	return func() T {

		if memoized {
			return result
		}
		result = f()
		memoized = true
		return result
	}
}
