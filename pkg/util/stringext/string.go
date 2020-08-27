package stringext

import "strings"

func ExistPrefix(pre []string, s string)  bool {
	for _, p := range pre {
		if strings.HasPrefix(s, p) {
			return true
		}
	}

	return false
}

func Reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

func TrimPreVal(val string, pre string) string {
	val = strings.ReplaceAll(val, "\"", "")
	return val[len(pre):]
}

func Split(r rune) bool {
	return r == '-' || r == '_'
}