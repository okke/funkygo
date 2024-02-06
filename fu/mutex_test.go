package fu

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
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
				log.Println("add " + fmt.Sprintf("%d", i))
				for j := 0; j < i; j++ {
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
