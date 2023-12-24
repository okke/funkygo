package main

import "context"

func SliceMap[I any, O any](s []I, f Transformer[I, O]) ([]O, error) {
	ctx := context.TODO()
	result := []O{}
	for _, x := range s {
		transformed, err := f(ctx, x)
		if err != nil {
			return result, err
		}
		result = append(result, transformed)
	}
	return result, nil
}
