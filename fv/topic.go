package fv

type Publisher[T any] func(T)
type Subscriber[T any] func(func(T))

func Topic[T any]() (Publisher[T], Subscriber[T]) {

	channel := make(chan T)
	subscribers := make([]chan T, 0)

	go func() {

		for {
			msg := <-channel
			for _, subChannel := range subscribers {
				subChannel <- msg
			}
		}
	}()

	return func(t T) {
			channel <- t
		},
		func(f func(T)) {
			subChannel := make(chan T)
			subscribers = append(subscribers, subChannel)

			go func() {
				for {
					f(<-subChannel)
				}
			}()
		}
}
