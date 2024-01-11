package fu

import "testing"

func TestConstruct(t *testing.T) {

	type testStruct struct {
		A string
		B int
	}

	withA := func(s string) Option[testStruct] {
		return func(t *testStruct) {
			t.A = s
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

	if result.A != "foo" {
		t.Errorf("Expected foo, got %s", result.A)
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
