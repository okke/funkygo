package fc

type Promised[T any] func() T

func Promise[T any](f func() T) Promised[T] {
	promised := make(chan T, 1)
	go func() {
		result := f()
		promised <- result
	}()

	return Memoize(TakePtr(promised))
}
