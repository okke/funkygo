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

	if positiveResult, err := positive(); err != nil {
		t.Error(err)
	} else {
		if positiveResult != "soup" {
			t.Errorf("Expected soup, got %s", positiveResult)
		}
	}

	if !wasPositive {
		t.Error("Expected positive soup")
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

type testStructForOptional struct {
	value string
}

func TestOptionalNilDo(t *testing.T) {

	Optional[testStructForOptional](nil).Do(func(s *testStructForOptional) {
		t.Errorf("Expected nil, got %v", s)
	})
}

func TestOptionalDo(t *testing.T) {

	test := &testStructForOptional{
		value: "soup",
	}
	done := false
	Optional(test).Do(func(s *testStructForOptional) {
		done = true
	})

	if !done {
		t.Error("Expected done to be true")
	}
}

func TestOptionalOr(t *testing.T) {

	done := false

	Optional[string](nil).Or(Ptr("soup")).Do(func(s *string) {
		done = *s == "soup"
	})

	if !done {
		t.Error("Expected done to be true")
	}
}

func TestOptionalNotExecutedOr(t *testing.T) {

	notDone := false

	OptionalP(Ptr("sauce")).Or(Ptr("soup")).Do(func(s *string) {
		notDone = *s == "sauce"
	})

	if !notDone {
		t.Error("Expected done to be true")
	}
}

func TestOptionalExists(t *testing.T) {

	noString := Optional[string](nil).Or(Nil[string]())

	if noString.Exists() {
		t.Error("Expected noString to be false")
	}
}
