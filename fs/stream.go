package fs

import "funcgo/fu"

type Stream[T any] func() (T, Stream[T])

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

func Each[T any](stream Stream[T], callback func(T) error) error {

	for e, s := stream(); s != nil; e, s = s() {
		if err := callback(e); err != nil {
			return err
		}
	}
	return nil
}

func Filter[T any](stream Stream[T], filter func(T) bool) Stream[T] {

	return func() (T, Stream[T]) {

		if value, next := stream(); next == nil {
			return fu.Zero[T](), nil
		} else if filter(value) {
			return value, Filter(next, filter)
		} else {
			return Filter(next, filter)()
		}

	}
}

func Map[I, O any](stream Stream[I], mapper func(I) O) Stream[O] {

	return func() (O, Stream[O]) {
		value, next := stream()
		if next == nil {
			return fu.Zero[O](), nil
		}
		return mapper(value), Map(next, mapper)
	}
}

func ToSet[T comparable](stream Stream[T]) fu.Set[T] {
	set := make(fu.Set[T])

	for {
		value, next := stream()
		if next == nil {
			return set
		}
		set = set.Add(value)
		stream = next
	}

}
