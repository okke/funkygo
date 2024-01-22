package fs

import "github.com/okke/funkygo/fu"

func While[T any](start T, while func(T) bool, next func(T) T) Stream[T] {

	value := start

	return func() (T, Stream[T]) {

		if !while(value) {
			return fu.Zero[T](), nil
		}

		return value, While(next(value), while, next)
	}
}

func Range[T fu.Number](start T, end T, step T) Stream[T] {
	if start < end {
		return While(start, fu.Lte(end), fu.Add(step))
	} else {
		return While(start, fu.Gte(end), fu.Add(step))
	}
}

func Endless[T any](value T) Stream[T] {
	return func() (T, Stream[T]) {
		return value, Endless(value)
	}
}

func EndlessIncrement[T fu.Number](initial T, step T) Stream[T] {
	return func() (T, Stream[T]) {
		return initial, EndlessIncrement(initial+step, step)
	}
}
