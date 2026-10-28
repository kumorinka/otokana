// Package deck: membangun deck dan mengacaknya dengan Fisher-Yates.
//
// Go murni, tanpa import mygo.
package deck

import "otokana/internal/kana"

// Settings are the five choices of Setup.
type Settings struct {
	Hiragana, Katakana               bool
	Monographs, Diacritics, Digraphs bool
}

// Default is hiragana monographs.
var Default = Settings{Hiragana: true, Monographs: true}

// Card is one kana in one script.
type Card struct {
	ID     string // "<script>:<romaji>": し and シ are two cards
	Entry  *kana.Entry
	Script kana.Script
	Order  int // place in the unshuffled deck, for listing cards in gojūon order
	// Pitch is how many semitones above the voicebank's base note the
	// card plays, below it when negative.
	Pitch  int
	Done   bool
	Hinted bool // its hint was shown
}

// Kana returns the kana the card shows.
func (c *Card) Kana() string { return c.Entry.Kana(c.Script) }

// Scripts returns the chosen scripts.
func (s Settings) Scripts() []kana.Script {
	var out []kana.Script
	if s.Hiragana {
		out = append(out, kana.Hiragana)
	}
	if s.Katakana {
		out = append(out, kana.Katakana)
	}
	return out
}

// Groups returns the groups of the chosen variants, in building order.
// Digraphs + diacritics come only with both diacritics and digraphs.
func (s Settings) Groups() []kana.Group {
	var out []kana.Group
	if s.Monographs {
		out = append(out, kana.Mono)
	}
	if s.Diacritics {
		out = append(out, kana.Dia)
	}
	if s.Digraphs {
		out = append(out, kana.Di)
	}
	if s.Diacritics && s.Digraphs {
		out = append(out, kana.DiDia)
	}
	return out
}

// Build returns the cards of s, unshuffled: a card for every entry of the
// chosen groups in every chosen script. It is empty without a script or a
// variant.
func Build(s Settings) []Card {
	scripts := s.Scripts()
	var cards []Card
	for _, g := range s.Groups() {
		for i := range kana.Table {
			e := &kana.Table[i]
			if e.Group != g {
				continue
			}
			for _, sc := range scripts {
				cards = append(cards, Card{ID: sc.String() + ":" + e.Romaji, Entry: e, Script: sc, Order: len(cards)})
			}
		}
	}
	return cards
}

// Shuffle shuffles cards with Fisher–Yates, from the last element back,
// with intN returning a number in [0, n), such as rand.IntN.
func Shuffle(cards []Card, intN func(n int) int) {
	for i := len(cards) - 1; i > 0; i-- {
		j := intN(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}
