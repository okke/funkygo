package fs

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"unicode"

	"github.com/okke/funkygo/fu"
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

func TestPeekOnEmptyStream(t *testing.T) {

	if _, stream := Peek(FromSlice([]int{})); stream != nil {
		t.Errorf("Expected nil, got %v", stream)
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

func TestPeekUntil(t *testing.T) {

	peeked, stream := PeekUntil(FromChannel(fu.C(1, 1, 1, 2, 3)), func(x int) bool {
		return x != 1
	})

	if count := Count(peeked); count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}

	if count := Count(stream); count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}
}

func TestTakeN(t *testing.T) {

	taken, stream := TakeN(FromChannel(fu.C(1, 2, 3)), 2)

	if count := Count(taken); count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}

	if count := Count(stream); count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}

}

func TestTakeNByTakingTooMuch(t *testing.T) {

	taken, stream := TakeN(FromChannel(fu.C(1, 2, 3)), 5)

	if count := Count(taken); count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}

	if count := Count(stream); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestTakeUntil(t *testing.T) {

	taken, stream := TakeUntil(FromChannel(fu.C(1, 1, 1, 2, 3)), func(x int) bool {
		return x != 1
	})

	if count := Count(taken); count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}

	if count := Count(stream); count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}

func TestTakeNAfterTakeUntil(t *testing.T) {
	ignored, _ := TakeUntil(FromArgs('a', 'b', 'C'), unicode.IsUpper)
	firstThreeIgnored, _ := TakeN(ignored, 3) // can only take two

	if count := Count(firstThreeIgnored); count != 2 {
		t.Error("Expected 2, got", count)
	}
}

func TestTakeUntilByTakingEverything(t *testing.T) {

	taken, stream := TakeUntil(FromChannel(fu.C(1, 1, 1, 2, 3)), func(x int) bool {
		return false
	})

	if count := Count(taken); count != 5 {
		t.Error("Expected 5, got", count)
	}

	if count := Count(stream); count != 0 {
		t.Error("Expected 0, got", count)
	}
}

func TestEach(t *testing.T) {

	count := 0
	Each(FromSlice([]int{}), func(x int) {
		count += x
	})
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func TestEachUntil(t *testing.T) {

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	count := 0
	stream := EachUntil(FromSlice(slice), fu.Gte(5), func(x int) {
		count++
	})
	if count != 4 {
		t.Error("Expected 4, got", count)
	}
	if s := ToSlice(stream); !reflect.DeepEqual(s, []int{5, 6, 7, 8, 9}) {
		t.Error("Expected [5, 6, 7, 8, 9], got", s)
	}
}

func TestEachUntilError(t *testing.T) {

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	count := 0

	stream, err := EachUntilError(FromSlice(slice), func(x int) error {
		if x == 5 {
			return fmt.Errorf("error")
		}
		count++
		return nil
	})

	if err == nil {
		t.Error("Expected error")
	}

	if count != 4 {
		t.Error("Expected 4, got", count)
	}

	if s := ToSlice(stream); !reflect.DeepEqual(s, []int{5, 6, 7, 8, 9}) {
		t.Error("Expected [5, 6, 7, 8, 9], got", s)
	}
}

func TestEchUntilNotError(t *testing.T) {

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	count := 0
	stream, err := EachUntilError(FromSlice(slice), fu.Safe(func(x int) {
		count++
	}))
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if count != 9 {
		t.Error("Expected 9, got", count)
	}
	if stream != nil {
		t.Error("Expected nil, got", stream)
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
	} else {
		if isEmpty, _ := IsEmpty(rest); !isEmpty {
			t.Errorf("Expected empty stream, got %v", rest)
		}
	}

	if matched, rest := MatchFirst(FromSlice([]int{1, 2, 3}), 1); !matched {
		t.Errorf("Did expect 1")
	} else {
		if c := Count(rest); c != 2 {
			t.Errorf("Expected 2, got %d", c)
		}
	}

	if matched, rest := MatchFirst(FromSlice([]int{1, 2, 3}), 5); matched {
		t.Errorf("Did not expect 5, got %v", rest)
	} else {
		if c := Count(rest); c != 3 {
			t.Errorf("Expected 3, got %d", c)
		}
	}
}

func TestDistinct(t *testing.T) {

	distinct := Distinct(FromSlice([]int{1, 2, 2, 1, 3}))
	if matched, _ := MatchFirst(distinct, 1, 2, 3); !matched {
		t.Errorf("Expected 1, 2 and 3, got %v", distinct)
	}

}

func TestDistinctBy(t *testing.T) {

	type Person struct {
		Name string
	}

	distinct := DistinctBy(FromSlice([]Person{
		{Name: "John"},
		{Name: "Jane"},
		{Name: "John"},
		{Name: "Jane"},
	}), func(x Person) string {
		return x.Name
	})

	if Count(distinct) != 2 {
		t.Errorf("Expected 2, got %v", distinct)
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

func TestFlatMap(t *testing.T) {

	stream := FlatMap(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) Stream[int] {
		return FromSlice([]int{x, x})
	})

	arr := ToSlice(stream)

	if !reflect.DeepEqual(arr, []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}) {
		t.Errorf("Expected [1, 1, 2, 2, 3, 3, 4, 4, 5, 5], got %v", arr)
	}
}

func TestFlatMapWithEmpty(t *testing.T) {

	stream := FlatMap(FromSlice([]int{1, 2, 3, 4, 5}), func(x int) Stream[int] {
		if x%2 == 0 {
			return FromSlice([]int{x, x})
		}
		return Empty[int]()
	})

	arr := ToSlice(stream)

	if !reflect.DeepEqual(arr, []int{2, 2, 4, 4}) {
		t.Errorf("Expected [2, 2, 4, 4], got %v", arr)
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

	empty := Sequence(Empty[int](), Empty[int]())
	if count := Count(empty); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	notSoEmpty := Sequence(FromSlice([]int{1, 2, 3}), Empty[int](), FromSlice([]int{1, 2, 3}))
	if count := Count(notSoEmpty); count != 6 {
		t.Errorf("Expected 3, got %d", count)
	}
}

func TestSequenceWithEmptyStreams(t *testing.T) {

	if count := Count(Sequence(Empty[int]())); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	if count := Count(Sequence(Empty[int](), Empty[int](), Empty[int]())); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}

	slice := ToSlice(Sequence(FromArgs(1, 2, 3), Empty[int](), FromArgs(4, 5, 6)))
	if !reflect.DeepEqual(slice, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Expected [1, 2, 3, 4, 5, 6], got %v", slice)
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

func TestAppend(t *testing.T) {

	stream := Append(FromSlice([]int{1, 2, 3}), 5, 6, 7)

	if count := Count(stream); count != 6 {
		t.Errorf("Expected 6, got %d", count)
	}

	slice := ToSlice(stream)
	if expected := []int{1, 2, 3, 5, 6, 7}; !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v, got %v", expected, slice)
	}
}

func TestPrepend(t *testing.T) {

	stream := Prepend(FromSlice([]int{1, 2, 3}), 5, 6, 7)

	if count := Count(stream); count != 6 {
		t.Errorf("Expected 6, got %d", count)
	}

	slice := ToSlice(stream)
	if expected := []int{5, 6, 7, 1, 2, 3}; !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v, got %v", expected, slice)
	}
}

func TestChopN(t *testing.T) {

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream := ChopN(FromSlice(slice), 3)

	if count := Count(stream); count != 4 {
		t.Errorf("Expected 4, got %d", count)
	}

	chopLengths := ToSlice(Map(ChopN(FromSlice(slice), 3), func(x []int) (int, error) {
		return len(x), nil
	}))

	if !reflect.DeepEqual(chopLengths, []int{3, 3, 3, 1}) {
		t.Errorf("Expected %v, got %v", []int{3, 3, 3, 1}, chopLengths)
	}
}

func TestChop(t *testing.T) {

	slice := []int{10, 11, 12, 20, 21, 22, 30, 31, 32}
	stream := Chop(FromSlice(slice), func(x1, x2 int) bool {
		return x2-x1 > 1
	})

	chop, next := stream()
	if !reflect.DeepEqual(chop, []int{10, 11, 12}) {
		t.Errorf("Expected %v, got %v", []int{10, 11, 12}, chop)
	}

	chop, next = next()
	if !reflect.DeepEqual(chop, []int{20, 21, 22}) {
		t.Errorf("Expected %v, got %v", []int{20, 21, 22}, chop)
	}

	chop, next = next()
	if !reflect.DeepEqual(chop, []int{30, 31, 32}) {
		t.Errorf("Expected %v, got %v", []int{30, 31, 32}, chop)
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
		}), func(x int) {
			count++
		})

		b.Log("set with length", count)
	}
}
