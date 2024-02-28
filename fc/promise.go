package fc

import (
	"github.com/okke/funkygo/fu"
)

type Promised[T any] func() T

func Promise[T any](f func() T) Promised[T] {
	promised := make(chan T, 1)
	go func() {
		result := f()
		promised <- result
	}()

	synchronize := NewMutex()
	result := fu.Optional[T](nil)

	return func() T {
		synchronize(func() {
			result = result.WhenNil(TakePtr(promised))
		})
		sure, _ := result()
		return *sure
	}
}
