package fc

import (
	"fmt"
	"testing"

	"github.com/okke/funkygo/fs"
)

func TestTopic(t *testing.T) {
	pub, sub := Topic[int](TopicBufSize(16), SubscriberBufSize(16))

	total := 0

	waitMore := WaitN(2, func(done func()) {

		sub(func(x int) {
			total += x
			done()
		})
		sub(func(x int) {
			total += x
			done()
		})

		pub(42)
	})

	if total != 84 {
		t.Errorf("Expected 42, got %d", total)
	}

	waitMore(2, func(done func()) {
		pub(-42)
	})

	if total != 0 {
		t.Errorf("Expected 0, got %d", total)
	}

}

func TestTopicUnsubscribe(t *testing.T) {

	fs.Each(fs.Range(0, 5, 1), func(x int) error {
		pub, sub := Topic[string]()
		var unsubscribe UnSubscriber

		waitMore := WaitN(2, func(done func()) {
			unsubscribe = sub(func(s string) {
				done()
			})

			sub(func(s string) {
				done()
			})

			pub(fmt.Sprintf("1:%d", x))
		})

		unsubscribe()

		waitMore(1, func(done func()) {
			pub(fmt.Sprintf("2:%d", x))
		})

		return nil
	})

}
