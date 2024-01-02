package fs

import "github.com/okke/funkygo/fu"

func FromSlice[T any](slice []T) Stream[T] {

	return func() (T, Stream[T]) {
		if len(slice) == 0 {
			return fu.Zero[T](), nil
		}
		return slice[0], FromSlice(slice[1:])
	}
}

func ToSlice[T any](stream Stream[T]) []T {

	result := []T{}
	for e, s := stream(); s != nil; e, s = s() {
		result = append(result, e)
	}
	return result
}
