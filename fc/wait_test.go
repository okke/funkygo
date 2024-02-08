package fc

import (
	"runtime"
	"testing"
)

func TestWaitN(t *testing.T) {

	count := 0
	WaitN(2, func(done func()) {
		go func() {
			count++
			done()
			runtime.Gosched()
			count++
			done()
		}()
	})(3, func(doneAsWell func()) {
		go func() {
			count++
			doneAsWell()
			runtime.Gosched()
			count++
			doneAsWell()
			runtime.Gosched()
			count++
			doneAsWell()
		}()
	})

	if count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}
}
