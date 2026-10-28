// Package format: format waktu, accuracy, dan nama deck.
//
// Go murni, tanpa import mygo.
package format

import (
	"fmt"
	"math"
	"strings"

	"otokana/internal/deck"
)

// Clock shows ms as MM:SS, or H:MM:SS from an hour on.
func Clock(ms int64) string {
	h, m, s := split(ms)
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// Duration writes ms out, as "2 min 48 s", leaving out zero parts.
func Duration(ms int64) string {
	h, m, s := split(ms)
	var parts []string
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%d h", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%d min", m))
	}
	if s > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d s", s))
	}
	return strings.Join(parts, " ")
}

func split(ms int64) (h, m, s int64) {
	sec := max(ms, 0) / 1000
	return sec / 3600, sec / 60 % 60, sec % 60
}

// Accuracy is the percentage of cards answered without a hint, rounded.
func Accuracy(cards, hinted int) int {
	if cards <= 0 {
		return 0
	}
	return int(math.Round(100 * float64(cards-hinted) / float64(cards)))
}

// DeckName names the deck of s, as "Hiragana basic".
func DeckName(s deck.Settings) string {
	var scripts, variants []string
	for _, sc := range s.Scripts() {
		scripts = append(scripts, sc.String())
	}
	for _, v := range []struct {
		on   bool
		name string
	}{{s.Monographs, "basic"}, {s.Diacritics, "dakuten"}, {s.Digraphs, "yōon"}} {
		if v.on {
			variants = append(variants, v.name)
		}
	}
	name := list(scripts) + " " + list(variants)
	return strings.ToUpper(name[:1]) + name[1:]
}

// list joins words as English does: "a", "a and b", "a, b, and c".
func list(words []string) string {
	switch len(words) {
	case 0:
		return ""
	case 1:
		return words[0]
	case 2:
		return words[0] + " and " + words[1]
	}
	return strings.Join(words[:len(words)-1], ", ") + ", and " + words[len(words)-1]
}
