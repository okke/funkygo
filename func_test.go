package main

import "testing"

func TestPipe(t *testing.T) {

	input := 3

	// apply := T[int, int](func(x int) int { return x + 1 })
	// output := Compose[int, int](input, apply)
	output := Pipe[int, int](input,
		IntoF(func(x int) float64 { return float64(x) + 1 }),
		IntoF(func(x float64) int { return int(x) + 1 }))

	if output != 5 {
		t.Errorf("Expected 4, got %d", output)
	}

}

func TestPipe2(t *testing.T) {

	input := []int{1, 2, 3, 4}

	// f := C1(sliceFilter, func(x int) bool { return x%2 == 0 })
	f := CurryOutTransformerArguments(CreateTransformerWithArguments[[]int, []int](sliceFilter[int]), func(x int) bool { return x%2 == 0 })

	output := Pipe[[]int, []int](
		input,
		IntoF(f),
		IntoAnyF[[]int, []int](sliceFilter[int], func(x int) bool { return x%2 == 0 }),
	)

	if len(output) != 2 {
		t.Errorf("Expected 2, got %d", len(output))
	}
}

func C1[I, O, A any](f func(I, A) O, p A) Transformer[I, O] {
	return func(i I) O {
		return f(i, p)
	}
}

func CP1[I, O, A any](f func(I, A) O, p A) Transformer[any, any] {
	return IntoF(C1[I, O, A](f, p))
}

func sliceFilter[T any](s []T, f func(T) bool) []T {

	result := []T{}

	for _, x := range s {
		if f(x) {
			result = append(result, x)
		}
	}

	return result
}
