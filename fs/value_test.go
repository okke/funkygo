package fs

import "testing"

func TestFromValue(t *testing.T) {
	if count := Count(FromValue(1)); count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}
}

func TestEmpty(t *testing.T) {
	if count := Count(Empty[int]()); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestIsEmpty(t *testing.T) {
	if isEmpty, _ := IsEmpty(Empty[int]()); !isEmpty {
		t.Errorf("Expected empty stream")
	}

	if isEmpty, _ := IsEmpty(FromArgs[int]()); !isEmpty {
		t.Errorf("Expected empty stream")
	}
}
