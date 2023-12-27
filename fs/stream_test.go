package fs

import (
	"math"
	"math/rand"
	"testing"
)

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

func TestMap(t *testing.T) {

	stream := Map(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) int {
		return x * 2
	})

	set := ToSet(stream)

	if !set.ContainsAll(2, 4, 6, 8, 10) {
		t.Errorf("Expected 2, 4, 6, 8 and 10, got %v", set)
	}
}
