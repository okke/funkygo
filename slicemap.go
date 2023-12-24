package main

func SliceMap[I any, O any](s []I, f Transformer[I, O]) ([]O, error) {
	result := []O{}
	for _, x := range s {
		transformed, err := f(x)
		if err != nil {
			return result, err
		}
		result = append(result, transformed)
	}
	return result, nil
}
