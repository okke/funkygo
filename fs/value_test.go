package fs

import "testing"

func TestEmpty(t *testing.T) {
	if count := Count(Empty[int]()); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestFromValue(t *testing.T) {
	if count := Count(FromValue(1)); count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}
}
