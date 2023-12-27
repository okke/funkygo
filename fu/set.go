package fu

type Set[T comparable] map[T]struct{}

func (set Set[T]) Add(value T) Set[T] {
	set[value] = struct{}{}
	return set
}

func (set Set[T]) Contains(value T) bool {
	_, ok := set[value]
	return ok
}

func (set Set[T]) ContainsAll(values ...T) bool {
	for _, value := range values {
		if !set.Contains(value) {
			return false
		}
	}
	return true
}

func (set Set[T]) ContainsNone(values ...T) bool {
	for _, value := range values {
		if set.Contains(value) {
			return false
		}
	}
	return true
}
