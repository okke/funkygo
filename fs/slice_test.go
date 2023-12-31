package fs

import (
	"reflect"
	"testing"
)

func TestFromAndToSlice(t *testing.T) {

	slice := []int{1, 2, 3}
	stream := FromSlice(slice)
	result := ToSlice(stream)
	if !reflect.DeepEqual(result, slice) {
		t.Errorf("Expected %v, got %v", slice, result)
	}
}
