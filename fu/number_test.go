package fu

import "testing"

func TestNumericOperations(t *testing.T) {

	if !Eq(1)(1) {
		t.Error("Expected 1 == 1")
	}

	if !Ne(1)(2) {
		t.Error("Expected 1 != 2")
	}

	if !Gt(1)(2) {
		t.Error("Expected 2 > 1")
	}

	if !Gte(1)(1) {
		t.Error("Expected 1 >= 1")
	}

	if !Lt(2)(1) {
		t.Error("Expected 1 < 2")
	}

	if !Lte(1)(1) {
		t.Error("Expected 1 <= 1")
	}

	if Add(1)(2) != 3 {
		t.Error("Expected 2 + 1 = 3")
	}

	if Subtract(2)(1) != -1 {
		t.Error("Expected 1 - 2 = -1")
	}

	if Multiply(2)(3) != 6 {
		t.Error("Expected 3 * 2 = 6")
	}

	if Divide(2)(6) != 3 {
		t.Error("Expected 6 / 2 = 3")
	}

	if Mod(3)(10) != 1 {
		t.Error("Expected 10 % 3 = 1")
	}

	if Increment[int]()(3) != 4 {
		t.Error("Expected 3 + 1 = 4")
	}

	if Decrement[int]()(3) != 2 {
		t.Error("Expected 3 - 1 = 2")
	}
}
