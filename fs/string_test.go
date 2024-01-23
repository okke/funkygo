package fs

import (
	"testing"
	"unicode"
)

func TestFromString(t *testing.T) {

	if Count(FromString("")) != 0 {
		t.Error("Expected 0, got", Count(FromString("")))
	}

	if Count(FromString("abc")) != 3 {
		t.Error("Expected 3, got", Count(FromString("abc")))
	}

	if c := Count(Filter(FromString("aBcDeFg"), unicode.IsUpper)); c != 3 {
		t.Error("Expected 3, got", c)
	}

}
