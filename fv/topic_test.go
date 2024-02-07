package fv

import (
	"fmt"
	"testing"

	"github.com/okke/funkygo/fs"
)

func TestTopic(t *testing.T) {
	pub, sub := Topic[int](TopicBufSize(16), SubscriberBufSize(16))

	total := 0
	done := make(chan struct{}, 16)
	sub(func(x int) {
		total += x
		done <- struct{}{}
	})
	sub(func(x int) {
		total += x
		done <- struct{}{}
	})
	pub(42)

	<-done
	<-done

	if total != 84 {
		t.Errorf("Expected 42, got %d", total)
	}

	pub(-42)

	<-done
	<-done

	if total != 0 {
		t.Errorf("Expected 0, got %d", total)
	}

}

func TestTopicUnsubscribe(t *testing.T) {

	fs.Each(fs.Range(0, 10, 1), func(x int) error {
		pub, sub := Topic[string]()

		calls := 0
		done := make(chan struct{}, 16)
		unsubscribe := sub(func(s string) {
			calls++
			done <- struct{}{}
		})

		sub(func(s string) {
			calls++
			done <- struct{}{}
		})

		pub(fmt.Sprintf("1:%d", x))
		<-done
		<-done

		if calls != 2 {
			t.Fatalf("Expected 1, got %d", calls)
		}

		unsubscribe()

		pub(fmt.Sprintf("2:%d", x))

		<-done

		if calls != 3 {
			t.Fatalf("Expected 3 call, got %d", calls)
		}
		return nil
	})

}
