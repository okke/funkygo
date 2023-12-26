package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestPipe(t *testing.T) {

	input := 3

	// apply := T[int, int](func(x int) int { return x + 1 })
	// output := Compose[int, int](input, apply)
	output, err := Pipe[int, int](context.TODO(), input,
		IntoF(T(func(x int) float64 { return float64(x) + 1 })),
		IntoF(T(func(x float64) int { return int(x) + 1 })),
	)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if output != 5 {
		t.Errorf("Expected 4, got %d", output)
	}

}

func TestPipeWithError(t *testing.T) {

	input := 42

	output, err := Pipe[int, int](context.TODO(), input,
		IntoF(T(func(x int) float64 { return float64(x) + 1 })),
		//IntoF(func(f float64) (float64, error) { return 0, errors.New("oh no") }),
		IntoAnyF[float64, float64](func(ctx context.Context, f float64) (float64, error) { return 0, errors.New("oh no") }),
		IntoF(T(func(x float64) int { return int(x) + 1 })),
	)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// output should be default value of an int which is 0
	//
	if output != 0 {
		t.Errorf("Expected 0, got %d", output)
	}
}

func TestPipeWithWrongTransformer(t *testing.T) {
	input := 33

	_, err := Pipe[int, int](context.TODO(), input,
		IntoAnyF[int, int](0),
	)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

}

func TestPipeWithNoResults(t *testing.T) {
	input := 33

	output, err := Pipe[int, int](context.TODO(), input,
		IntoAnyF[int, int](func(ctx context.Context, x int) {}), // do nothing, return nothing
	)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if output != 0 {
		t.Errorf("Expected 0, got %d", output)
	}

}

func TestPipe2(t *testing.T) {

	input := []int{1, 2, 3, 4}

	adder := func(in []int) []int {
		out := []int{}
		for _, x := range in {
			out = append(out, x+1)
		}
		return out
	}

	output, err := Pipe[[]int, []int](
		context.TODO(), input,
		IntoAnyF[[]int, []int](sliceFilter[int], func(x int) bool { return x%3 == 0 }),
		IntoF(T(adder)),
		IntoAnyF[[]int, []int](sliceFilterWithoutContext[int], func(x int) bool { return x%2 == 0 }),
	)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(output) != 1 {
		t.Errorf("Expected 2, got %d", len(output))
	}
}

func TestBind(t *testing.T) {

	input := 42

	output, err := Pipe[int, string](context.TODO(), input,
		Bind("Step1", T(func(x int) int { return x + 1 })),
		Bind("Step2", T(func(x int) int { return x * 2 })),
		Bind("Step3", T(func(x int) int { return x - 1 })),
		IntoF(func(ctx context.Context, x int) (string, error) {
			return fmt.Sprintf("%d-%d-%d", ctx.Value("Step1"), ctx.Value("Step2"), ctx.Value("Step3")), nil
		}),
	)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if output != "43-86-85" {
		t.Errorf("Expected \"43-86-85\", got %s", output)
	}

}

func TestPromiseShouldBeMemoized(t *testing.T) {

	input := 42

	output, err := Pipe[int, int](context.TODO(), input,
		IntoF(Promise(T(func(x int) int { return x + 1 }))),
		IntoF[PromisedValue[int], PromisedValue[int]](T(func(x PromisedValue[int]) PromisedValue[int] {
			actualX, _ := x()
			t.Log(actualX)
			return x
		})),
		IntoF[PromisedValue[int], int](T(func(x PromisedValue[int]) int {
			actualX, _ := x()
			t.Log(actualX)
			return actualX
		})),
	)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if output != 43 {
		t.Errorf("Expected 43, got %d", output)
	}

}

func TestAs(t *testing.T) {

	type Recipe struct {
		Name       string
		Difficulty int
	}

	input := "soup"

	output, err := Pipe[string, *Recipe](context.TODO(), input,
		Bind("Name", T(strings.ToUpper)),
		Bind("Difficulty", Promise(T(func(s string) int { return len(s) }))),
		As(&Recipe{}),
	)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if output.Name != "SOUP" {
		t.Errorf("Expected \"SOUP\", got %s", output.Name)
	}
	if output.Difficulty != 4 {
		t.Errorf("Expected 4, got %d", output.Difficulty)
	}

}

func sliceFilter[T any](ctx context.Context, s []T, f func(T) bool) []T {

	result := []T{}

	for _, x := range s {
		if f(x) {
			result = append(result, x)
		}
	}

	return result
}

func sliceFilterWithoutContext[T any](s []T, f func(T) bool) []T {

	result := []T{}

	for _, x := range s {
		if f(x) {
			result = append(result, x)
		}
	}

	return result
}
