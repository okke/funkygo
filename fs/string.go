package fs

import "strings"

func FromString(s string) Stream[rune] {
	return func() (rune, Stream[rune]) {
		if len(s) == 0 {
			return 0, nil
		}
		return rune(s[0]), FromString(s[1:])
	}
}

func Runes2Lines(stream Stream[rune]) Stream[string] {

	if stream == nil {
		return Empty[string]()
	}

	return func() (string, Stream[string]) {
		var sb strings.Builder
		var r rune

		for r, stream = stream(); stream != nil; r, stream = stream() {
			switch r {
			case '\n':
				return sb.String(), Runes2Lines(stream)
			case '\r':
				r2, input := Peek(stream)
				if r2 == '\n' {
					_, input = input()
					return sb.String(), Runes2Lines(input)
				} else {
					sb.WriteRune(r)
				}
			default:
				sb.WriteRune(r)
			}
		}
		return sb.String(), Runes2Lines(stream)
	}
}
