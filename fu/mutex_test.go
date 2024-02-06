package fu

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWithMutex(t *testing.T) {

	var m sync.Mutex
	var builder strings.Builder

	waitGroup := sync.WaitGroup{}

	for i := 0; i < 10; i++ {

		waitGroup.Add(1)
		go func(i int) {
			defer waitGroup.Done()
			WithMutex(&m, func() {
				for j := 0; j < i; j++ {

					// sleep for a little bit
					//
					time.Sleep(1 * time.Millisecond)
					builder.WriteString(fmt.Sprintf("%d", i))
				}
			})
		}(i)
	}
	waitGroup.Wait()

	for i := 0; i < 10; i++ {
		check := ""
		for j := 0; j < i; j++ {
			check += fmt.Sprintf("%d", i)
		}
		if !strings.Contains(builder.String(), check) {
			t.Errorf("Expected %s to contain %s", builder.String(), check)
		}
	}
}
