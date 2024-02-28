package fc

import "github.com/okke/funkygo/fu"

func Memoize[T any](f func() *T) func() T {

	synchronize := NewMutex()
	result := fu.Optional[T](nil)

	return func() T {
		synchronize(func() {
			result = result.WhenNil(f)
		})
		sure, _ := result()
		return *sure
	}
}
