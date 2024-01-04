package fs

import (
	"reflect"
	"testing"

	"github.com/okke/funkygo/fu"
)

func TestWhile(t *testing.T) {

	if count := Count(While(0, fu.Lt(5), fu.Add(1))); count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}
}

func TestRange(t *testing.T) {

	if count := Count(Range(1, 5, 1)); count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}

	slice := ToSlice(Range(0, 9, 1))

	if expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}; !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v, got %v", expected, slice)
	}
}

func TestEndless(t *testing.T) {

	if count := Count(Limit(Endless(1), 1_000_000)); count != 1_000_000 {
		t.Errorf("Expected 1, got %d", count)
	}
}

func TestEndlessIncrement(t *testing.T) {

	slice := ToSlice(Limit(EndlessIncrement(0, 1), 1_000_000))

	if len(slice) != 1_000_000 {
		t.Errorf("Expected 1_000_000, got %d", len(slice))
	}

	if slice[0] != 0 {
		t.Errorf("Expected 0, got %d", slice[0])
	}

	if slice[len(slice)-1] != 999_999 {
		t.Errorf("Expected 999_999, got %d", slice[len(slice)-1])
	}

}
