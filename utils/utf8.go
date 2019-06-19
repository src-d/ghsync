package utils

import (
	"strings"
	"unicode/utf8"
)

// UTF8String returns valid utf8 string by removing incorrect characters
func UTF8String(s string) string {
	// postgres doesn't support NULL characters
	s = strings.Replace(s, "\x00", "", -1)

	if utf8.ValidString(s) {
		return s
	}

	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				continue
			}
		}
		v = append(v, r)
	}

	return string(v)
}
