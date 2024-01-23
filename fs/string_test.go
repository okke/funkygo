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

func TestRunes2Lines(t *testing.T) {

	non := ToSlice(Runes2Lines(nil))
	if len(non) != 0 {
		t.Error("Expected 0, got", len(non))
	}

	lines := ToSlice(Runes2Lines(FromString("abc\ndef")))
	if lines[0] != "abc" {
		t.Error("Expected 'abc', got", lines[0])
	}
	if lines[1] != "def" {
		t.Error("Expected 'def', got", lines[1])
	}

	lines = ToSlice(Runes2Lines(FromString("abc\r\ndef")))
	if lines[0] != "abc" {
		t.Error("Expected 'abc', got", lines[0])
	}
	if lines[1] != "def" {
		t.Error("Expected 'def', got", lines[1])
	}

	lines = ToSlice(Runes2Lines(FromString("abc\rdef")))
	if lines[0] != "abcdef" {
		t.Error("Expected 'abcdef', got", lines[0])
	}

	lines = ToSlice(Runes2Lines(FromString("\n")))
	if lines[0] != "" {
		t.Error("Expected '', got", lines[0])
	}

	lines = ToSlice(Runes2Lines(FromString("\n\n")))
	if lines[1] != "" {
		t.Error("Expected '', got", lines[0])
	}

}
