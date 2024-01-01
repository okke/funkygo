package fs

import "testing"

func TestFromChannel(t *testing.T) {

	channel := make(chan int, 10)
	channel <- 1
	channel <- 2
	channel <- 3
	close(channel)

	stream := FromChannel(channel)

	count := 0
	Each(stream, func(x int) error {
		count += x
		return nil
	})

	if count != 6 {
		t.Errorf("Expected 6, got %d", count)
	}

}

func TestToChannel(t *testing.T) {
	channel := make(chan int, 10)
	stream := FromSlice([]int{1, 2, 3, 4, 5})

	ToChannel(stream, channel)

	count := 0
	for v, hasValue := <-channel; hasValue; v, hasValue = <-channel {
		count += v
	}

	if count != 15 {
		t.Errorf("Expected 15, got %d", count)
	}
}
