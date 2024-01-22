package fu

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

func Eq[T Number](v T) func(T) bool {
	return func(n T) bool {
		return n == v
	}
}

func Ne[T Number](v T) func(T) bool {
	return func(n T) bool {
		return n != v
	}
}

func Gt[T Number](v T) func(T) bool {
	return func(n T) bool {
		return n > v
	}
}

func Gte[T Number](v T) func(T) bool {
	return func(n T) bool {
		return n >= v
	}
}

func Lt[T Number](v T) func(T) bool {
	return func(n T) bool {
		return n < v
	}
}

func Lte[T Number](v T) func(T) bool {
	return func(n T) bool {
		return n <= v
	}
}

func Add[T Number](v T) func(T) T {
	return func(n T) T {
		return n + v
	}
}

func Subtract[T Number](v T) func(T) T {
	return func(n T) T {
		return n - v
	}
}

func Multiply[T Number](v T) func(T) T {
	return func(n T) T {
		return n * v
	}
}

func Divide[T Number](v T) func(T) T {
	return func(n T) T {
		return n / v
	}
}

func Mod[T constraints.Integer](v T) func(T) T {
	return func(n T) T {
		return n % v
	}
}

func Increment[T Number]() func(T) T {
	return Add[T](1)
}

func Decrement[T Number]() func(T) T {
	return Subtract[T](1)
}
