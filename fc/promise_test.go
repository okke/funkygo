package fc

import "testing"

func TestPromise(t *testing.T) {

	called := 0
	promise := Promise(func() int {
		called++
		return 42
	})

	if actual := promise(); 42 != actual {
		t.Errorf("Expected %d, got %d", 42, actual)
	}

	// call promise again
	//
	promise()

	if called != 1 {
		t.Errorf("Expected 1, got %d", called)
	}
}
