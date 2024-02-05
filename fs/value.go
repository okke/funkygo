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

func IsEmpty[T any](stream Stream[T]) (bool, Stream[T]) {
	if stream == nil {
		return true, nil
	}
	value, next := stream()
	return next == nil, func() (T, Stream[T]) {
		return value, next
	}
}
