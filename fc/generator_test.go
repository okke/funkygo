package fc

import (
	"testing"
)

func TestGenerate(t *testing.T) {

	counter := 0
	done := make(chan struct{}, 16)
	sub := Generate(func() (int, error) {
		counter++
		return counter, StopGeneratingWhen(counter > 20, func() {
			done <- struct{}{}
		})
	}, Frequency(1000 /* times every second */), GenerateWithTopicOptions(TopicBufSize(1), SubscriberBufSize(1)))

	received := 0
	sub(func(x int) {
		received++
	})

	<-done

	if received != 20 {
		t.Errorf("Expected 21, got %d", received)
	}

}
