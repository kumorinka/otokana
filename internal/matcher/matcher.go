// Package matcher: menilai ketikan: benar, tunggu, atau salah.
//
// Go murni, tanpa import mygo.
package matcher

import "strings"

// Result is how typed input compares with a card's answers.
type Result int

const (
	Correct Result = iota // the input is one of the answers
	Pending               // the input begins one of them
	Wrong                 // neither
)

// Match compares input with aliases, Correct before Pending.
func Match(input string, aliases []string) Result {
	for _, a := range aliases {
		if input == a {
			return Correct
		}
	}
	for _, a := range aliases {
		if strings.HasPrefix(a, input) {
			return Pending
		}
	}
	return Wrong
}
