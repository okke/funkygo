package fv

import "testing"

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

	pub, sub := Topic[int](TopicBufSize(16), SubscriberBufSize(16))

	calls := 0
	done := make(chan struct{}, 16)
	unsubscribe := sub(func(x int) {
		calls++
		done <- struct{}{}
	})

	sub(func(x int) {
		calls++
		done <- struct{}{}
	})

	pub(42)
	<-done
	<-done

	if calls != 2 {
		t.Errorf("Expected 1, got %d", calls)
	}

	unsubscribe()

	pub(42)

	<-done
	select {
	case <-done:
		t.Errorf("Expected only one call to be finished")
	default:
		// ok
	}

	if calls != 3 {
		t.Errorf("Expected 3 call, got %d", calls)
	}

}
