package fs

import "github.com/okke/funkygo/fu"

func FromChannel[T any](channel <-chan T) Stream[T] {
	return func() (T, Stream[T]) {
		value, ok := <-channel
		if !ok {
			return fu.Zero[T](), nil
		}
		return value, FromChannel(channel)
	}
}

func ToChannel[T any](stream Stream[T], channel chan<- T) {
	for e, s := stream(); s != nil; e, s = s() {
		channel <- e
	}
	close(channel)
}
