package fs

import (
	"strings"

	"github.com/okke/funkygo/fu"
)

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

		stream, _ := EachUntil(stream,
			fu.Eq('\n'),
			fu.Safe(func(r rune) {
				if r != '\r' {
					sb.WriteRune(r)
				}
			}))

		if stream != nil {
			_, stream = stream()
		}

		return sb.String(), Runes2Lines(stream)
	}
}
