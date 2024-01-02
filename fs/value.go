package fs

import "github.com/okke/funkygo/fu"

func FromValue[T any](value T) Stream[T] {
	return func() (T, Stream[T]) {
		return value, Empty[T]()
	}
}

func Empty[T any]() Stream[T] {
	return func() (T, Stream[T]) {
		return fu.Zero[T](), nil
	}
}
