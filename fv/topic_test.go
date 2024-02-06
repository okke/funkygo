package fv

import "testing"

func TestTopic(t *testing.T) {
	pub, sub := Topic[int]()

	total := 0
	done := make(chan struct{})
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
