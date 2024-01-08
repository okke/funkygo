package fu

import (
	"errors"
	"testing"
)

func TestTryPositive(t *testing.T) {

	var found string
	Try(func() (string, error) {
		return "soup", nil
	}).OnError(func(err error) {
		t.Error(err)
	}).OnSuccess(func(s string) {
		found = s
	})

	if found != "soup" {
		t.Errorf("Expected soup, got %s", found)
	}
}

func TestTryNegative(t *testing.T) {

	var found string
	Try(func() (string, error) {
		return "", errors.New("soup")
	}).OnError(func(err error) {
		found = err.Error()
	}).OnSuccess(func(s string) {
		t.Error("Should not be called")
	})

	if found != "soup" {
		t.Errorf("Expected soup, got %s", found)
	}
}

func TestTryAndReturn(t *testing.T) {

	wasPositive := false
	positive := func() (string, error) {

		return Try(func() (string, error) {
			return "soup", nil
		}).OnSuccess(func(s string) {
			wasPositive = true
		}).Return()
	}

	if !wasPositive {
		t.Error("Expected wasPositive to be true")
	}

	if positiveResult, err := positive(); err != nil {
		t.Error(err)
	} else {
		if positiveResult != "soup" {
			t.Errorf("Expected soup, got %s", positiveResult)
		}
	}

	negative := func() (string, error) {

		return Try(func() (string, error) {
			return "", errors.New("soup")
		}).Return()
	}

	if _, err := negative(); err == nil {
		t.Error("Expected error")
	}
}

func TestOptional(t *testing.T) {

	type testStruct struct {
		value string
	}

	Optional[testStruct](nil).Do(func(s *testStruct) {
		t.Errorf("Expected nil, got %v", s)
	})

	test := &testStruct{
		value: "soup",
	}
	done := false
	Optional(test).Do(func(s *testStruct) {
		done = true
	})

	if !done {
		t.Error("Expected done to be true")
	}

	done = false
	Optional[testStruct](nil).Do(func(s *testStruct) {
		done = true
	})

	if done {
		t.Error("Expected done to be false")
	}

	done = false

	Optional[testStruct](nil).Or(func() *testStruct {
		return &testStruct{
			value: "soup",
		}
	}).Do(func(s *testStruct) {
		done = s.value == "soup"
	})

	noString := Optional[string](nil).Or(func() *string {
		return nil
	})

	if noString.Exists() {
		t.Error("Expected noString to be false")
	}
}
