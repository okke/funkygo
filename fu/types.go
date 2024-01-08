package fu

type LeftOrRight bool

const (
	IsLeft  LeftOrRight = false
	IsRight LeftOrRight = true
)

type CanError[T any] func() (T, error)
type Either[L, R any] func() (L, R, LeftOrRight)
type ValueOrError[T any] Either[T, error]
type OptionalValue[T any] func() (*T, bool)
type Pointer[T any] func() *T

func Left[L, R any](left L) Either[L, R] {
	return func() (L, R, LeftOrRight) {
		return left, Zero[R](), IsLeft
	}
}

func Right[L, R any](right R) Either[L, R] {
	return func() (L, R, LeftOrRight) {
		return Zero[L](), right, IsRight
	}
}

func valueOrError[T any](v Either[T, error]) ValueOrError[T] {
	return func() (T, error, LeftOrRight) {
		return v()
	}
}

func Try[T any](f CanError[T]) ValueOrError[T] {
	if result, err := f(); err != nil {
		return valueOrError(Right[T, error](err))
	} else {
		return valueOrError(Left[T, error](result))
	}
}

func (v ValueOrError[T]) OnSuccess(f func(T)) ValueOrError[T] {
	actual, _, leftOrRight := v()
	if leftOrRight == IsLeft {
		f(actual)
	}
	return v
}

func (v ValueOrError[T]) OnError(f func(error)) ValueOrError[T] {
	_, actual, leftOrRight := v()
	if leftOrRight == IsRight {
		f(actual)
	}
	return v
}

func (v ValueOrError[T]) Return() (T, error) {
	result, err, _ := v()
	return result, err
}

func Optional[T any](value *T) OptionalValue[T] {
	return func() (*T, bool) {
		return value, value != nil
	}
}

func OptionalP[T any](value Pointer[T]) OptionalValue[T] {
	return func() (*T, bool) {
		v := value()
		return v, v != nil
	}
}

func (o OptionalValue[T]) Exists() bool {
	_, exists := o()
	return exists
}

func (o OptionalValue[T]) Do(with func(*T)) {

	if actual, exists := o(); exists {
		with(actual)
	}
}

func (o OptionalValue[T]) Or(fill func() *T) OptionalValue[T] {

	if _, exists := o(); exists {
		return o
	} else {
		return Optional(fill())
	}
}

func Ptr[T any](value T) Pointer[T] {
	return func() *T {
		return &value
	}
}

func Nil[T any]() Pointer[T] {
	return func() *T {
		return nil
	}
}
