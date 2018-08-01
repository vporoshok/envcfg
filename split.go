package envcfg

import (
	"strings"
)

func SplitWords(s string) []string {
	var res []string

	chunks := strings.Split(s, "_")
	for _, chunk := range chunks {
		res = append(res, splitChunk(chunk)...)
	}

	return res
}

func splitChunk(chunk string) []string {
	const (
		unknown = iota
		lower
		number
		upper
	)

	isUpper := func(c byte) bool { return c >= 'A' && c <= 'Z' }
	isNumber := func(c byte) bool { return c >= '0' && c <= '9' }

	if len(chunk) == 0 {

		return nil
	}
	var res []string
	word := &strings.Builder{}
	prev := unknown
	for i := len(chunk) - 1; i >= 0; i-- {
		c := chunk[i]
		word.WriteByte(c)

		switch true {
		case isUpper(c):
			if prev == lower {
				res = append(res, reverseString(word.String()))
				word.Reset()
				prev = unknown

				continue
			}

			prev = upper

		case isNumber(c):
			prev = number

		default:
			if prev == upper {
				w := word.String()
				res = append(res, reverseString(w[:len(w)-1]))
				word.Reset()
				word.WriteByte(c)
			}

			prev = lower
		}
	}

	if word.Len() > 0 {
		res = append(res, reverseString(word.String()))
	}

	return reverseSlice(res)
}

func reverseString(s string) string {
	res := &strings.Builder{}
	for i := len(s) - 1; i >= 0; i-- {
		res.WriteByte(s[i])
	}

	return res.String()
}

func reverseSlice(s []string) []string {
	res := make([]string, len(s))
	for i := range res {
		res[i] = s[len(s)-i-1]
	}

	return res
}
