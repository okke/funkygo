package fs

import (
	"github.com/okke/funkygo/fu"
)

type Stream[T any] func() (T, Stream[T])

func Peek[T any](stream Stream[T]) (T, Stream[T]) {

	value, next := stream()
	if next == nil {
		return value, nil
	}

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
		var value T
		if value, stream = stream(); stream == nil {
			return FromSlice(values), Empty[T]()
		} else {
			values = append(values, value)
		}
	}

	return FromSlice(values), stream
}

func TakeUntil[T any](stream Stream[T], until func(T) bool) (Stream[T], Stream[T]) {

	values := []T{}

	var value T
	for value, stream = stream(); stream != nil && !until(value); value, stream = stream() {
		values = append(values, value)
	}

	if stream == nil {
		return FromSlice(values), Empty[T]()
	}

	return FromSlice(values), Sequence(FromValue(value), stream)
}

func Each[T any](stream Stream[T], callback func(T)) {

	for value, next := stream(); next != nil; value, next = next() {
		callback(value)
	}
}

func EachUntil[T any](stream Stream[T], until func(T) bool, callback func(T)) Stream[T] {

	for {
		var value T
		value, stream = stream()
		if stream == nil {
			return stream
		}

		if until(value) {
			return func() (T, Stream[T]) { return value, stream }
		}

		callback(value)
	}
}

func EachUntilError[T any](stream Stream[T], callback func(T) error) (Stream[T], error) {

	for {
		var value T
		value, stream = stream()
		if stream == nil {
			return stream, nil
		}

		if err := callback(value); err != nil {
			return func() (T, Stream[T]) { return value, stream }, err
		}
	}
}

func Count[T any](stream Stream[T]) int {
	count := 0
	Each(stream, func(x T) {
		count++
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

func FlatMap[I, O any](stream Stream[I], mapper func(I) Stream[O]) Stream[O] {

	return func() (O, Stream[O]) {
		value, next := stream()
		if next == nil {
			return fu.Zero[O](), nil
		}

		mapped := mapper(value)
		mappedIsEmpty, mapped := IsEmpty(mapped)
		if mappedIsEmpty {
			return FlatMap(next, mapper)()
		}

		mappedStream, moreToMap := mapped()

		return mappedStream, TwoStreams(moreToMap, FlatMap(next, mapper))
	}
}

func Reduce[T any](stream Stream[T], reducer func(T, T) T) T {

	result, next := stream()
	if next == nil {
		return result
	}

	Each(next, func(value T) {
		result = reducer(result, value)
	})

	return result
}

func ReduceInto[I, O any](stream Stream[I], initial O, reducer func(O, I) O) O {

	result := initial

	Each(stream, func(value I) {
		result = reducer(result, value)
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

func ToPointers[T any](stream Stream[T]) Stream[*T] {
	isEmpty, stream := IsEmpty(stream)
	if isEmpty {
		return func() (*T, Stream[*T]) {
			return nil, nil
		}
	}
	return Map(stream, func(x T) (*T, error) {
		return &x, nil
	})
}

func TwoStreams[T any](stream1 Stream[T], stream2 Stream[T]) Stream[T] {
	isEmpty1, stream1 := IsEmpty(stream1)
	if isEmpty1 {
		return stream2
	}

	isEmpty2, stream2 := IsEmpty(stream2)
	if isEmpty2 {
		return stream1
	}

	return func() (T, Stream[T]) {
		value, next1 := stream1()
		return value, TwoStreams(next1, stream2)
	}
}

func Sequence[T any](streams ...Stream[T]) Stream[T] {

	streamOfStreams := Map(Filter(FromSlice(streams), func(x Stream[T]) bool { return x != nil }),
		func(s Stream[T]) (Stream[*T], error) {
			return ToPointers(s), nil
		})

	if streamOfStreams == nil {
		return Empty[T]()
	}

	current, streamOfStreams := streamOfStreams()

	return Map(Filter(sequenceOfPointerStreams(current, streamOfStreams), func(x *T) bool {
		return x != nil
	}), func(x *T) (T, error) {
		return *x, nil
	})
}

func sequenceOfPointerStreams[T any](current Stream[*T], streamOfStreams Stream[Stream[*T]]) Stream[*T] {

	return func() (*T, Stream[*T]) {

		isEmpty, current := IsEmpty(current)

		if isEmpty {
			current, streamOfStreams = streamOfStreams()
			if streamOfStreams == nil {
				return nil, nil
			}
			return nil, sequenceOfPointerStreams(current, streamOfStreams)
		}

		value, next := current()
		return value, sequenceOfPointerStreams(next, streamOfStreams)
	}
}

func Append[T any](stream Stream[T], values ...T) Stream[T] {
	return Sequence(stream, FromSlice(values))
}

func Prepend[T any](stream Stream[T], values ...T) Stream[T] {
	return Sequence(FromSlice(values), stream)
}

func ChopN[T any](stream Stream[T], amount int) Stream[[]T] {

	return func() ([]T, Stream[[]T]) {

		result := make([]T, 0, amount)

		for i := 0; i < amount; i++ {
			var value T
			value, stream = stream()
			if stream == nil {
				return result, Empty[[]T]()
			}
			result = append(result, value)
		}

		return result, ChopN(stream, amount)
	}
}

func Chop[T any](stream Stream[T], shouldChop func(T, T) bool) Stream[[]T] {

	return func() ([]T, Stream[[]T]) {

		result := []T{}

		for {
			current, more := stream()
			if more == nil {
				return result, Empty[[]T]()
			}
			result = append(result, current)
			next, more := Peek(more)
			if more == nil {
				return result, Empty[[]T]()
			}

			if shouldChop(current, next) {
				return result, Chop(more, shouldChop)
			}

			stream = more
		}
	}
}

func ToSet[T comparable](stream Stream[T]) fu.Set[T] {
	set := make(fu.Set[T])

	Each(stream, func(value T) {
		set = set.Add(value)
	})

	return set
}
