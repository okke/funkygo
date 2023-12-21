package main

func SliceMap[I any, O any](s []I, f Transformer[I, O]) []O {
	result := []O{}
	for _, x := range s {
		result = append(result, f(x))
	}
	return result
}
