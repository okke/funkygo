package fv

import (
	"github.com/okke/funkygo/fs"
	"github.com/okke/funkygo/fu"
)

/*
A publisher is a function taking a message and publishing it to the topic
*/
type Publisher[T any] func(T)

/*
A subscriber is a function taking a handler for messages and returning a function
which can be called to unsubscribe from the topic
*/
type Subscriber[T any] func(func(T)) func()

type TopicOptions struct {
	topicBufSize      int
	subscriberBufSize int
}

func TopicBufSize(size int) fu.Option[TopicOptions] {
	return func(options *TopicOptions) {
		options.topicBufSize = size
	}
}

func SubscriberBufSize(size int) fu.Option[TopicOptions] {
	return func(options *TopicOptions) {
		options.subscriberBufSize = size
	}
}

/*
Topic creates a topic and returns the publish and the subscribe functions
*/
func Topic[T any](options ...fu.Option[TopicOptions]) (Publisher[T], Subscriber[T]) {

	opts := fu.With(&TopicOptions{
		topicBufSize:      256,
		subscriberBufSize: 256,
	}, options...)

	topicChannel := make(chan *T, opts.topicBufSize)
	subscribeChannels := []chan *T{}
	unsubscribeChannel := make(chan struct{}, 16)

	synchronized := NewMutex()

	go func() {
		for {
			dispatch(synchronized, topicChannel, subscribeChannels)
		}
	}()

	go func() {
		for {
			<-unsubscribeChannel
			synchronized(func() {
				subscribeChannels = fs.ToSlice(fs.Filter(fs.FromSlice(subscribeChannels), func(c chan *T) bool { return c != nil }))
			})
		}
	}()

	return func(t T) {
			topicChannel <- &t
		},

		// subscribe
		//
		func(f func(T)) func() {

			doneChannel := make(chan struct{}, 16)
			subChannel := make(chan *T, opts.subscriberBufSize)

			var subscriberIndex int
			synchronized(func() {
				subscriberIndex = len(subscribeChannels)
				subscribeChannels = append(subscribeChannels, subChannel)
			})

			go func() {
				for {
					select {
					case <-doneChannel:
						return
					case msg := <-subChannel:
						if msg != nil {
							go f(*msg)
						}
					default:
						// do nothing
					}
				}
			}()

			// unsubscribe
			//
			return func() {

				synchronized(func() {
					close(subscribeChannels[subscriberIndex])
					subscribeChannels[subscriberIndex] = nil
					doneChannel <- struct{}{}
					unsubscribeChannel <- struct{}{}
				})
			}
		}
}

func dispatch[T any](synchronized Mutex, topicChannel chan *T, subscribers []chan *T) {

	select {
	case msg := <-topicChannel:
		if msg == nil {
			return
		}

		synchronized(func() {
			for _, subChannel := range subscribers {
				if subChannel != nil {
					subChannel <- msg
				}
			}
		})
	default:
		return
	}

}
