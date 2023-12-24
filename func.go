package main

import (
	"context"
	"errors"
	"reflect"
)

type Transformer[I any, O any] func(context.Context, I) (O, error)
type SimpleTransformer[I any, O any] func(I) O
type TransformerWithArguments[I any, O any] func(context.Context, I, ...any) (O, error)

func T[I, O any](t SimpleTransformer[I, O]) Transformer[I, O] {
	return func(ctx context.Context, i I) (O, error) {
		return t(i), nil
	}
}

func Pipe[I any, O any](ctx context.Context, initial I, steps ...Transformer[any, any]) (O, error) {

	var result any = initial
	for _, step := range steps {
		stepResult, err := step(ctx, result)
		if err != nil {
			return *new(O), err
		}
		result = stepResult
	}
	return result.(O), nil
}

func IntoF[I any, O any](t Transformer[I, O]) Transformer[any, any] {

	return func(ctx context.Context, x any) (any, error) {

		functionValue := reflect.ValueOf(t)

		return result2TransformerResult[O](functionValue.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(x),
		}))

	}
}

func result2TransformerResult[O any](result []reflect.Value) (O, error) {
	if len(result) == 0 {
		return *new(O), nil
	} else if len(result) == 2 && !result[1].IsNil() {
		return result[0].Interface().(O), result[1].Interface().(error)
	}
	return result[0].Interface().(O), nil
}

func IntoAnyF[I any, O any](transformerFunc any, arguments ...any) Transformer[any, any] {
	return IntoF(curryOutTransformerArguments(createTransformerWithArguments[I, O](transformerFunc), arguments...))
}

func createTransformerWithArguments[I, O any](f any) TransformerWithArguments[I, O] {
	return func(ctx context.Context, i I, arguments ...any) (O, error) {
		functionValue := reflect.ValueOf(f)

		// Check if the variable is a function
		if functionValue.Kind() != reflect.Func {
			return *new(O), errors.New("transformer is not a function")
		}

		args, _ := SliceMap(arguments, T(func(x any) reflect.Value {
			return reflect.ValueOf(x)
		}))

		allArgs := append([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(i)}, args...)

		return result2TransformerResult[O](functionValue.Call(allArgs))
	}
}

func curryOutTransformerArguments[I, O any](f TransformerWithArguments[I, O], p ...any) Transformer[I, O] {
	return func(ctx context.Context, i I) (O, error) {
		return f(ctx, i, p...)
	}
}
