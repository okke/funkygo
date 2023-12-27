package fp

import (
	"context"
	"errors"
	"fmt"
	"funcgo/fs"
	"funcgo/fu"
	"reflect"
)

type Transformer[I any, O any] func(context.Context, I) (O, error)
type SimpleTransformer[I any, O any] func(I) O
type TransformerWithArguments[I any, O any] func(context.Context, I, ...any) (O, error)

type PromisedValue[O any] func() (O, error)
type BoundValue[O any] func() (string, O)

var _CONTEXT_TYPE = reflect.TypeOf((*context.Context)(nil)).Elem()

var _CONTEXT_KEYS struct{} = struct{}{}

func T[I, O any](t SimpleTransformer[I, O]) Transformer[I, O] {
	return func(ctx context.Context, i I) (O, error) {
		return t(i), nil
	}
}

func Promise[I, O any](t Transformer[I, O]) Transformer[I, PromisedValue[O]] {

	return func(ctx context.Context, i I) (PromisedValue[O], error) {
		channel := make(chan PromisedValue[O], 1)
		go func() {
			result, err := t(ctx, i)
			channel <- func() (O, error) {
				return result, err
			}
		}()
		return fu.MemoizeWithError(func() (O, error) {
			result := <-channel
			return result()
		}), nil
	}
}

func Pipe[I any, O any](ctx context.Context, initial I, steps ...Transformer[any, any]) (O, error) {

	var result any = initial
	for _, step := range steps {
		stepResult, err := step(ctx, result)
		if err != nil {
			return fu.Zero[O](), err
		}

		ctx, result = bindIntoContext(ctx, stepResult)
	}
	return result.(O), nil
}

func Demux[I any, O any](ctx context.Context, initial I, steps ...Transformer[any, any]) (fs.Stream[O], error) {

	return fs.Map(fs.FromSlice(steps), func(t Transformer[any, any]) (O, error) {
		result, err := t(ctx, initial)
		return result.(O), err
	}), nil
}

func bindIntoContext(ctx context.Context, value any) (context.Context, any) {

	functionValue := reflect.ValueOf(value)
	if functionValue.Kind() != reflect.Func {
		return ctx, value
	}

	if functionValue.Type().NumIn() != 0 {
		return ctx, value
	}

	if functionValue.Type().NumOut() != 2 {
		return ctx, value
	}

	if functionValue.Type().Out(0).Kind() != reflect.String {
		return ctx, value
	}

	result := functionValue.Call([]reflect.Value{})
	name := fmt.Sprintf("%s", result[0].Interface())
	actualValue := result[1].Interface()

	// store values in own mapping for later usage
	//
	mapping := ctx.Value(_CONTEXT_KEYS)
	if mapping == nil {
		mapping = map[string]any{}
		ctx = context.WithValue(ctx, _CONTEXT_KEYS, mapping)
	}
	mapping.(map[string]any)[name] = actualValue

	// and store them in context as well
	//
	return context.WithValue(ctx, name, actualValue), actualValue
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

func Bind[I any, O any](name string, t Transformer[I, O]) Transformer[any, any] {

	apply := IntoF(t)
	return func(ctx context.Context, x any) (any, error) {
		result, err := apply(ctx, x)
		if err != nil {
			return fu.Zero[O](), err
		}
		return bind[O](name, result.(O)), nil
	}
}

func bind[O any](name string, value O) BoundValue[O] {
	return func() (string, O) {
		return name, value
	}
}

func As[O any](into O) Transformer[any, any] {

	return func(ctx context.Context, x any) (any, error) {

		mapping := ctx.Value(_CONTEXT_KEYS)
		if mapping == nil {
			return into, nil
		}

		for key, value := range mapping.(map[string]any) {

			// if value is a function, call it and store the result in the mapping
			//
			if reflect.TypeOf(value).Kind() == reflect.Func {

				functionValue := reflect.ValueOf(value)
				functionResult := functionValue.Call([]reflect.Value{})

				mapping.(map[string]any)[key] = functionResult[0].Interface()
			}
		}

		return fu.Map2Struct(mapping.(map[string]any), into)
	}
}

func result2TransformerResult[O any](result []reflect.Value) (O, error) {
	if len(result) == 0 {
		return fu.Zero[O](), nil
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
			return fu.Zero[O](), errors.New("transformer is not a function")
		}

		args := fs.ToSlice(fs.Map(fs.FromSlice(arguments), func(argValue any) (reflect.Value, error) {
			return reflect.ValueOf(argValue), nil
		}))

		allArgs := append([]reflect.Value{reflect.ValueOf(i)}, args...)

		// prepend the context if the transformer function accepts it
		//
		if functionValue.Type().NumIn() > 0 && functionValue.Type().In(0) == _CONTEXT_TYPE {
			allArgs = append([]reflect.Value{reflect.ValueOf(ctx)}, allArgs...)
		}

		return result2TransformerResult[O](functionValue.Call(allArgs))
	}
}

func curryOutTransformerArguments[I, O any](f TransformerWithArguments[I, O], p ...any) Transformer[I, O] {
	return func(ctx context.Context, i I) (O, error) {
		return f(ctx, i, p...)
	}
}
