package fc

func TakePtr[T any](c <-chan T) func() *T {
	return func() *T {
		var result T
		result = <-c
		return &result
	}
}
