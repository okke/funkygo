package fs

import (
	"errors"
	"fmt"
	"funcgo/fu"
	"math"
	"math/rand"
	"testing"
)

func TestPeek(t *testing.T) {

	stream := FromSlice([]int{1, 2, 3})

	if peeked, _ := Peek(stream); peeked != 1 {
		t.Errorf("Expected 1, got %d", peeked)
	}

	if peeked, _ := Peek(stream); peeked != 1 {
		t.Errorf("Expected 1, got %d", peeked)
	}

	if peeked, _ := Peek(FromSlice([]int{})); peeked != 0 {
		t.Errorf("Expected 0, got %d", peeked)
	}
}

func TestPeekWithChannel(t *testing.T) {

	stream := FromChannel(fu.C(1, 2, 3))

	for i := 0; i < 10; i++ {
		var peeked int
		if peeked, stream = Peek(stream); peeked != 1 {
			t.Errorf("Expected 1, got %d", peeked)
		}
	}

	if reduced := ReduceInto(stream, 0, func(x, y int) int { return x + y }); reduced != 6 {
		t.Errorf("Expected 6, got %d", reduced)
	}

}

func TestPeekN(t *testing.T) {

	peeked, stream := PeekN(FromChannel(fu.C(1, 2, 3)), 2)

	if count := Count(peeked); count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}

	if count := Count(stream); count != 3 {
		t.Errorf("Expected 1, got %d", count)
	}

	var peekedValue int
	if peekedValue, peeked = peeked(); peekedValue != 1 {
		t.Errorf("Expected 1, got %d", peekedValue)
	}

	var unPeekedValue int
	if unPeekedValue, stream = stream(); unPeekedValue != 1 {
		t.Errorf("Expected 1, got %d", unPeekedValue)
	}
}

func TestPeekNWithNotEnoughToPeek(t *testing.T) {

	peeked, stream := PeekN(FromChannel(fu.C(1, 2, 3)), 5)

	if count := Count(peeked); count != 3 {
		t.Errorf("Expected 2, got %d", count)
	}

	if count := Count(stream); count != 3 {
		t.Errorf("Expected 1, got %d", count)
	}
}

func TestHasMore(t *testing.T) {

	stream := FromSlice([]int{1, 2, 3})
	if !HasMore(stream) {
		t.Errorf("Expected true, got %v", HasMore(stream))
	}

	if more := HasMore(FromSlice([]int{})); more {
		t.Errorf("Expected false, got %v", more)
	}
}

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

func TestMatchFirst(t *testing.T) {

	if matched, rest := MatchFirst(FromSlice([]int{1, 2, 3}), 1, 2, 3); !matched {
		t.Errorf("Expected 1, 2 and 3, got %v", rest)
	}

	if matched, rest := MatchFirst(FromSlice([]int{1, 2, 3}), 5); matched {
		t.Errorf("Did not expect 5, got %v", rest)
	}
}

func TestDistinct(t *testing.T) {

	distinct := Distinct(FromSlice([]int{1, 2, 2, 1, 3}))
	if matched, _ := MatchFirst(distinct, 1, 2, 3); !matched {
		t.Errorf("Expected 1, 2 and 3, got %v", distinct)
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

func TestFindFirst(t *testing.T) {
	if found, _ := FindFirst(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x%2 == 0
	}); found != 2 {
		t.Errorf("Expected 2, got %d", found)
	}

	if found, rest := FindFirst(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x == 5
	}); !(found == 5 && rest != nil) {
		t.Errorf("Expected 5 and not nil, got %d, %v", found, rest)
	}

	if _, rest := FindFirst(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) bool {
		return x > 10
	}); rest != nil {
		t.Errorf("Expected nil, got %v", rest)
	}
}

func TestLimit(t *testing.T) {

	if count := Count(Limit(FromSlice([]int{1, 2, 3, 4, 5}), 3)); count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}
}

func TestSequence(t *testing.T) {

	stream := Sequence(FromSlice([]int{1, 2, 3}), FromSlice([]int{4, 5, 6}))

	count := ReduceInto(stream, 0, func(acc int, x int) int {
		return acc + x
	})

	if count != 21 {
		t.Errorf("Expected 21, got %d", count)
	}
}

func TestSequenceWithEmptyStreams(t *testing.T) {

	if count := Count(Sequence(Empty[int]())); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	if count := Count(Sequence(Empty[int](), Empty[int](), Empty[int]())); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestSequenceWithMultipleEmptyAndNonEmptyStream(t *testing.T) {

	stream := Sequence(Empty[int](), FromSlice([]int{1, 2, 3}), Empty[int](), Empty[int](), FromSlice([]int{4, 5, 6}), Empty[int]())

	if count := ReduceInto(stream, 0, func(acc int, x int) int {
		return acc + x
	}); count != 21 {
		t.Errorf("Expected 21, got %d", count)
	}

}

// ------------------------------------

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
