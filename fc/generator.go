package fc

import (
	"errors"
	"time"

	"github.com/okke/funkygo/fu"
)

type Generator[T any] func() (T, error)

var STOP_GENERATING error = errors.New("stop generating")

func StopGeneratingWhen(when bool, finalizers ...func()) error {
	if when {
		for _, finalizer := range finalizers {
			finalizer()
		}
		return STOP_GENERATING
	}
	return nil
}

type GeneratorOptions struct {
	topicOptions []fu.Option[TopicOptions]
	frequency    int // per second
}

func GenerateWithTopicOptions(topicOptions ...fu.Option[TopicOptions]) fu.Option[GeneratorOptions] {
	return func(options *GeneratorOptions) {
		options.topicOptions = append(options.topicOptions, topicOptions...)
	}
}

func Frequency(frequency int) fu.Option[GeneratorOptions] {
	return func(options *GeneratorOptions) {
		options.frequency = frequency
	}
}

func Generate[T any](generator Generator[T], options ...fu.Option[GeneratorOptions]) Subscriber[T] {

	opts := fu.With(&GeneratorOptions{
		topicOptions: []fu.Option[TopicOptions]{},
		frequency:    1,
	}, options...)

	interval := time.Duration(float64(time.Second) / float64(opts.frequency))
	ticker := time.NewTicker(interval)
	pub, sub := Topic[T](opts.topicOptions...)

	go func() {
		for {
			select {
			case <-ticker.C:
				if value, err := generator(); err != nil {
					if err == STOP_GENERATING {
						ticker.Stop()
						return
					}
				} else {
					pub(value)
				}
			}
		}
	}()

	return sub
}
