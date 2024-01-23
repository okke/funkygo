package fs

func FromString(s string) Stream[rune] {
	return func() (rune, Stream[rune]) {
		if len(s) == 0 {
			return 0, nil
		}
		return rune(s[0]), FromString(s[1:])
	}
}
