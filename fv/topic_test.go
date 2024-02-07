package fv

import (
	"fmt"
	"sync"
	"testing"

	"github.com/okke/funkygo/fs"
)

func TestTopic(t *testing.T) {
	pub, sub := Topic[int](TopicBufSize(16), SubscriberBufSize(16))

	total := 0

	wg := sync.WaitGroup{}

	sub(func(x int) {
		total += x
		wg.Done()
	})
	sub(func(x int) {
		total += x
		wg.Done()
	})

	wg.Add(2)
	pub(42)
	wg.Wait()

	if total != 84 {
		t.Errorf("Expected 42, got %d", total)
	}

	wg.Add(2)
	pub(-42)
	wg.Wait()

	if total != 0 {
		t.Errorf("Expected 0, got %d", total)
	}

}

func TestTopicUnsubscribe(t *testing.T) {

	fs.Each(fs.Range(0, 5, 1), func(x int) error {
		pub, sub := Topic[string]()

		wg := sync.WaitGroup{}

		unsubscribe := sub(func(s string) {
			wg.Done()
		})

		sub(func(s string) {
			wg.Done()
		})

		wg.Add(2)
		pub(fmt.Sprintf("1:%d", x))
		wg.Wait()

		unsubscribe()

		wg.Add(1)
		pub(fmt.Sprintf("2:%d", x))
		wg.Wait()

		return nil
	})

}
