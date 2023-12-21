package main

import "reflect"

type Transformer[I any, O any] func(I) O
type TransformerWithArguments[I any, O any] func(I, ...any) O

func Pipe[I any, O any](initial I, steps ...Transformer[any, any]) O {

	var result any = initial
	for _, step := range steps {
		result = step(result)
	}
	return result.(O)
}

func IntoF[I any, O any](t Transformer[I, O]) Transformer[any, any] {

	return func(x any) any {

		functionValue := reflect.ValueOf(t)

		// Check if the variable is a function
		if functionValue.Kind() == reflect.Func {
			// Prepare arguments for the function
			args := []reflect.Value{
				reflect.ValueOf(x),
			}

			// Call the function with arguments
			result := functionValue.Call(args)
			return result[0].Interface().(O)
		}

		panic("Value is not a function")
	}
}

func IntoAnyF[I any, O any](transformerFunc any, arguments ...any) Transformer[any, any] {
	return IntoF(CurryOutTransformerArguments(CreateTransformerWithArguments[I, O](transformerFunc), arguments...))
}

func CreateTransformerWithArguments[I, O any](f any) TransformerWithArguments[I, O] {
	return func(i I, arguments ...any) O {
		functionValue := reflect.ValueOf(f)

		// Check if the variable is a function
		if functionValue.Kind() != reflect.Func {
			panic("Value is not a function")
		}
		args := SliceMap(arguments, func(x any) reflect.Value {
			return reflect.ValueOf(x)
		})

		allArgs := append([]reflect.Value{reflect.ValueOf(i)}, args...)

		// Call the function with arguments
		result := functionValue.Call(allArgs)
		return result[0].Interface().(O)
	}
}

func CurryOutTransformerArguments[I, O any](f TransformerWithArguments[I, O], p ...any) Transformer[I, O] {
	return func(i I) O {
		return f(i, p...)
	}
}
