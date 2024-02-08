package fc

import (
	"runtime"
	"testing"

	"github.com/okke/funkygo/fs"
)

func TestWaitN(t *testing.T) {

	count := 0
	WaitN(2, func(done func()) {
		go func() {
			fs.Each(fs.RangeN(1), func(x int) error {
				go func() {
					runtime.Gosched()
					count++
					done()
				}()
				return nil
			})
		}()
	})(3, func(doneAsWell func()) {
		go func() {
			fs.Each(fs.RangeN(2), func(x int) error {
				go func() {
					runtime.Gosched()
					count++
					doneAsWell()
				}()
				return nil
			})
		}()
	})

	if count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}
}
