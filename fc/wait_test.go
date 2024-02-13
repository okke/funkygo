package fc

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/okke/funkygo/fs"
)

func TestWaitN(t *testing.T) {

	var count int64 = 0
	WaitN(2, func(done Done) {
		go func() {
			fs.Each(fs.RangeN(1), func(x int) {
				go func() {
					runtime.Gosched()
					atomic.AddInt64(&count, 1)
					done()
				}()
			})
		}()
	})(3, func(doneAsWell Done) {
		go func() {
			fs.Each(fs.RangeN(2), func(x int) {
				go func() {
					runtime.Gosched()
					atomic.AddInt64(&count, 1)
					doneAsWell()
				}()
			})
		}()
	})

	if count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}
}

func TestWait(t *testing.T) {

	for doManyTimes := 0; doManyTimes < 10; doManyTimes++ {

		var count int64 = 0
		expected := 500
		waitMore := Wait(func(submitTask TaskSubmitter) {

			for i := 0; i < expected; i++ {
				submitTask(func() {
					runtime.Gosched()
					atomic.AddInt64(&count, 1)
				})
			}
		})

		if int(count) != expected {
			t.Errorf("Expected %d, got %d", expected, count)
		}

		waitMore(func(submitTask TaskSubmitter) {

			for i := 0; i < expected; i++ {
				submitTask(func() {
					runtime.Gosched()
					atomic.AddInt64(&count, 1)
				})
			}
		})

		if int(count) != expected*2 {
			t.Errorf("Expected %d, got %d", expected*2, count)
		}
	}
}
