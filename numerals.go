package bsdate

import (
	"strconv"
	"strings"
)

// NepaliDigits formats n using Devanagari digits (०–९).
//
//	NepaliDigits(2082) // "२०८२"
func NepaliDigits(n int) string {
	s := strconv.Itoa(n)
	var b strings.Builder
	b.Grow(len(s) * 3) // each Devanagari digit is 3 bytes in UTF-8
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune('०' + (r - '0'))
		} else {
			b.WriteRune(r) // keep the minus sign
		}
	}
	return b.String()
}
