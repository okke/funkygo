package fu

import "testing"

func TestC(t *testing.T) {

	channel := C(1, 2, 3)

	if c, _ := <-channel; c != 1 {
		t.Errorf("Expected 1, got %d", c)
	}
	if c, _ := <-channel; c != 2 {
		t.Errorf("Expected 2, got %d", c)
	}
	if c, _ := <-channel; c != 3 {
		t.Errorf("Expected 3, got %d", c)
	}

	if _, ok := <-channel; ok {
		t.Errorf("Expected closed channel")
	}
}
