package fs

import "funcgo/fu"

type Stream[T any] func() (T, Stream[T])

func Peek[T any](stream Stream[T]) T {
	value, _ := stream()
	return value
}

func HasMore[T any](stream Stream[T]) bool {
	_, next := stream()
	return next != nil
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

func Distinct[T comparable](stream Stream[T]) Stream[T] {

	set := fu.Set[T]{}

	return Filter(stream, func(x T) bool {

		if set.Contains(x) {
			return false
		}
		set.Add(x)
		return true
	})

}

func MatchFirst[T comparable](stream Stream[T], values ...T) (bool, Stream[T]) {

	var found T
	next := stream
	for _, value := range values {
		if found, next = next(); found != value {
			return false, stream
		}
	}
	return true, next
}

func Map[I, O any](stream Stream[I], mapper func(I) (O, error)) Stream[O] {

	return func() (O, Stream[O]) {
		value, next := stream()
		if next == nil {
			return fu.Zero[O](), nil
		}
		if mapped, err := mapper(value); err != nil {
			return fu.Zero[O](), nil
		} else {
			return mapped, Map(next, mapper)
		}
	}
}

func Reduce[T any](stream Stream[T], reducer func(T, T) T) T {

	result, next := stream()
	if next == nil {
		return result
	}

	for value, next := next(); next != nil; value, next = next() {
		result = reducer(result, value)
	}

	return result
}

func ReduceInto[I, O any](stream Stream[I], initial O, reducer func(O, I) O) O {

	result := initial

	for value, next := stream(); next != nil; value, next = next() {
		result = reducer(result, value)
	}

	return result
}

func FindFirst[T any](stream Stream[T], predicate func(T) bool) (T, Stream[T]) {

	for value, next := stream(); next != nil; value, next = next() {
		if predicate(value) {
			return value, next
		}
	}
	return fu.Zero[T](), nil
}

func Limit[T any](stream Stream[T], limit int) Stream[T] {

	return func() (T, Stream[T]) {
		if limit == 0 {
			return fu.Zero[T](), nil
		}
		value, next := stream()
		return value, Limit(next, limit-1)
	}
}

func ToSet[T comparable](stream Stream[T]) fu.Set[T] {
	set := make(fu.Set[T])

	for value, next := stream(); next != nil; value, next = next() {
		set = set.Add(value)
	}

	return set
}
