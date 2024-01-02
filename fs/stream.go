package fs

import "github.com/okke/funkygo/fu"

type Stream[T any] func() (T, Stream[T])

func Peek[T any](stream Stream[T]) (T, Stream[T]) {
	value, next := stream()
	return value, func() (T, Stream[T]) {
		return value, next
	}
}

func PeekN[T any](stream Stream[T], n int) (Stream[T], Stream[T]) {

	values := make([]T, 0, n)

	for i := 0; i < n; i++ {
		if value, stream := stream(); stream == nil {
			return FromSlice(values), FromSlice(values)
		} else {
			values = append(values, value)
		}
	}

	return FromSlice(values), Sequence(FromSlice(values), stream)
}

func PeekUntil[T any](stream Stream[T], until func(T) bool) (Stream[T], Stream[T]) {

	values := []T{}

	var value T
	for value, stream = stream(); stream != nil && !until(value); value, stream = stream() {
		values = append(values, value)
	}

	if stream == nil {
		return FromSlice(values), FromSlice(values)
	}

	return FromSlice(values), Sequence(FromSlice(append(values, value)), stream)
}

func TakeN[T any](stream Stream[T], n int) (Stream[T], Stream[T]) {

	values := make([]T, 0, n)

	for i := 0; i < n; i++ {
		if value, stream := stream(); stream == nil {
			return FromSlice(values), Empty[T]()
		} else {
			values = append(values, value)
		}
	}

	return FromSlice(values), stream
}

func HasMore[T any](stream Stream[T]) bool {
	_, next := stream()
	return next != nil
}

func Each[T any](stream Stream[T], callback func(T) error) error {

	for value, next := stream(); next != nil; value, next = next() {
		if err := callback(value); err != nil {
			return err
		}
	}
	return nil
}

func Count[T any](stream Stream[T]) int {
	count := 0
	Each(stream, func(x T) error {
		count++
		return nil
	})
	return count
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

	Each(next, func(value T) error {
		result = reducer(result, value)
		return nil
	})

	return result
}

func ReduceInto[I, O any](stream Stream[I], initial O, reducer func(O, I) O) O {

	result := initial

	Each(stream, func(value I) error {
		result = reducer(result, value)
		return nil
	})

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

func Sequence[T any](streams ...Stream[T]) Stream[T] {
	current, streamOfStreams := FromSlice(streams)()
	if streamOfStreams == nil {
		return Empty[T]()
	}
	return sequenceOfStreams(current, streamOfStreams)
}

func sequenceOfStreams[T any](current Stream[T], streamOfStreams Stream[Stream[T]]) Stream[T] {

	return func() (T, Stream[T]) {

		value, next := current()

		for next == nil {
			current, streamOfStreams = streamOfStreams()
			if streamOfStreams == nil {
				return fu.Zero[T](), nil
			}

			value, next = current()
		}

		return value, sequenceOfStreams(next, streamOfStreams)
	}
}

func ToSet[T comparable](stream Stream[T]) fu.Set[T] {
	set := make(fu.Set[T])

	Each(stream, func(value T) error {
		set = set.Add(value)
		return nil
	})

	return set
}
