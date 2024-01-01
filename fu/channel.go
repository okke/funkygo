package fu

func C[T any](values ...T) chan T {

	channel := make(chan T, len(values))
	go func() {
		for _, value := range values {
			channel <- value
		}
		close(channel)
	}()
	return channel
}
