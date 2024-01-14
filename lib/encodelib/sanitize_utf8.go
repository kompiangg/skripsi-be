package encodelib

import "unicode/utf8"

func SanitizeUTF8String(str string) string {
	if utf8.ValidString(str) {
		return str
	}

	v := make([]rune, 0, len(str))
	for i, r := range str {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(str[i:])
			if size == 1 {
				continue
			}
		}
		v = append(v, r)
	}
	return string(v)
}
