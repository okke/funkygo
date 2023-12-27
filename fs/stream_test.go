package fs

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestEach(t *testing.T) {

	count := 0
	Each(FromSlice([]int{}), func(x int) error {
		count += x
		return nil
	})
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	Each(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) error {
		count += x
		if x == 3 {
			return errors.New("oh no")
		}
		return nil
	})

	if count != 6 {
		t.Errorf("Expected 15, got %d", count)
	}

}

func TestFilter(t *testing.T) {

	stream := Filter(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x%2 == 0
	})

	set := ToSet(stream)

	if !set.ContainsAll(2, 4) {
		t.Errorf("Expected 2 and 4, got %v", set)
	}

	if !set.ContainsNone(1, 3, 5) {
		t.Errorf("Expected none of 1, 3 and 5, got %v", set)
	}
}

func TestMap(t *testing.T) {

	stream := Map(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) (int, error) {
		return x * 2, nil
	})

	set := ToSet(stream)

	if !set.ContainsAll(2, 4, 6, 8, 10) {
		t.Errorf("Expected 2, 4, 6, 8 and 10, got %v", set)
	}
}

func TestReduce(t *testing.T) {

	if found := Reduce(FromSlice([]int{}), func(x, y int) int {
		return x + y
	}); found != 0 {
		t.Errorf("Expected 0, got %d", found)
	}

	if found := Reduce(FromSlice([]int{1}), func(x, y int) int {
		return x + y
	}); found != 1 {
		t.Errorf("Expected 1, got %d", found)
	}

	if found := Reduce(FromSlice([]int{1, 2, 3, 4, 5}), func(x, y int) int {
		return x + y
	}); found != 15 {

		t.Errorf("Expected 15, got %d", found)
	}
}

func TestReduceInto(t *testing.T) {

	if found := ReduceInto(FromSlice([]int{1, 2, 3}), "", func(acc string, x int) string {
		return fmt.Sprintf("%s%d", acc, x)
	}); found != "123" {
		t.Errorf("Expected 123, got %s", found)
	}
}

func createRandomSliceOfInts(size int) []int {

	slice := make([]int, size)
	for i := 0; i < size; i++ {

		// generate a random number
		//
		randomNumber := rand.Intn(math.MaxInt)
		slice[i] = randomNumber
	}
	return slice
}

var largeStreamOfInts = FromSlice(createRandomSliceOfInts(1_000_000))

func BenchmarkFilter(b *testing.B) {

	for i := 0; i < b.N; i++ {
		count := 0

		Each(Filter(largeStreamOfInts, func(x int) bool {
			return x%20000 == 0
		}), func(x int) error {
			count++
			return nil
		})

		b.Log("set with length", count)
	}
}
