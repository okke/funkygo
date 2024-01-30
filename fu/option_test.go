package fu

import (
	"strings"
	"testing"
)

func TestConstruct(t *testing.T) {

	type testStruct struct {
		A string
		B int
	}

	withA := func(s string) Option[testStruct] {
		return func(t *testStruct) {
			t.A = strings.ToUpper(s)
		}
	}

	withB := func(i int) Option[testStruct] {
		return func(t *testStruct) {
			t.B = i
		}
	}

	constructor := func(options ...Option[testStruct]) *testStruct {
		return Construct(options...)
	}

	result := constructor(withA("foo"), withB(42))

	if result.A != "FOO" {
		t.Errorf("Expected FOO, got %s", result.A)
	}

	if result.B != 42 {
		t.Errorf("Expected 42, got %d", result.B)
	}

	result = Construct(withB(43))

	if result.B != 43 {
		t.Errorf("Expected 43, got %d", result.B)
	}

	if result.A != "" {
		t.Errorf("Expected empty string, got %s", result.A)
	}
}

func TestWith(t *testing.T) {

	type testStruct struct {
		A string
	}

	withA := func(s string) Option[testStruct] {
		return func(t *testStruct) {
			t.A = strings.ToUpper(s)
		}
	}

	var s testStruct

	With(&s, withA("foo"))

	if s.A != "FOO" {
		t.Errorf("Expected FOO, got %s", s.A)
	}
}
