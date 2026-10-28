// Package session: mesin state sesi latihan.
//
// Go murni, tanpa import mygo.
package session

import (
	"math/rand/v2"
	"time"

	"otokana/internal/deck"
	"otokana/internal/matcher"
)

// Lock is how long input is ignored after a correct answer.
const Lock = 150 * time.Millisecond

// PitchRange is how many semitones at most a card plays from the
// voicebank's base note: further, a resampled voice sounds a chipmunk's.
const PitchRange = 3

// Phase is the screen the session is on.
type Phase int

const (
	Setup Phase = iota
	Playing
	Results
)

// Session is a game, from Setup through Results.
type Session struct {
	Phase     Phase
	Settings  deck.Settings
	Deck      []deck.Card
	Input     string
	Hint      bool // the active card shows its romaji
	Started   time.Time
	LockUntil time.Time
	Ended     time.Time
	// Shuffle shuffles decks; rand.IntN when nil.
	Shuffle func(n int) int
}

// New returns a session on Setup with the default settings.
func New() *Session { return &Session{Settings: deck.Default} }

// CanStart reports whether the settings give a deck.
func (s *Session) CanStart() bool { return len(deck.Build(s.Settings)) > 0 }

// Start builds and shuffles the deck of the settings and plays it from
// now. It does nothing when the deck is empty.
func (s *Session) Start(now time.Time) {
	cards := deck.Build(s.Settings)
	if len(cards) == 0 {
		return
	}
	intN := s.Shuffle
	if intN == nil {
		intN = rand.IntN
	}
	deck.Shuffle(cards, intN)
	for i := range cards {
		cards[i].Pitch = intN(2*PitchRange+1) - PitchRange
	}
	s.Phase, s.Deck = Playing, cards
	s.Input, s.Hint = "", false
	s.Started, s.LockUntil, s.Ended = now, time.Time{}, time.Time{}
}

// Active returns the active card: the last one not done, or nil.
func (s *Session) Active() *deck.Card {
	for i := len(s.Deck) - 1; i >= 0; i-- {
		if !s.Deck[i].Done {
			return &s.Deck[i]
		}
	}
	return nil
}

// Remaining returns how many cards are not done.
func (s *Session) Remaining() int {
	n := 0
	for _, c := range s.Deck {
		if !c.Done {
			n++
		}
	}
	return n
}

// Hinted returns how many cards had their hint shown, the cards Results
// counts as not right.
func (s *Session) Hinted() int {
	n := 0
	for _, c := range s.Deck {
		if c.Hinted {
			n++
		}
	}
	return n
}

// Elapsed returns the time played: until now while playing, until the
// last card on Results.
func (s *Session) Elapsed(now time.Time) time.Duration {
	switch s.Phase {
	case Playing:
		return now.Sub(s.Started)
	case Results:
		return s.Ended.Sub(s.Started)
	}
	return 0
}

// Result returns how the input compares with the active card.
func (s *Session) Result() matcher.Result {
	c := s.Active()
	if c == nil {
		return matcher.Pending
	}
	return matcher.Match(s.Input, c.Entry.Aliases)
}

// Type adds the letters a to z of text to the input, in lower case, and
// judges the input after each: a correct answer finishes the card, clears
// the input and locks it, a wrong one keeps the input, and counts for
// nothing. Text typed during the lock is dropped. It returns the card
// finished, if any, and whether a letter made the input wrong.
func (s *Session) Type(text string, now time.Time) (done *deck.Card, wrong bool) {
	for _, r := range text {
		if s.Phase != Playing || now.Before(s.LockUntil) {
			break
		}
		if 'A' <= r && r <= 'Z' {
			r += 'a' - 'A'
		}
		if r < 'a' || r > 'z' {
			continue
		}
		c := s.Active()
		s.Input += string(r)
		switch s.Result() {
		case matcher.Correct:
			c.Done = true
			done = c
			s.Input, s.Hint = "", false
			s.LockUntil = now.Add(Lock)
			if s.Active() == nil {
				s.Phase, s.Ended = Results, now
			}
		case matcher.Wrong:
			wrong = true
		}
	}
	return done, wrong
}

// Backspace deletes the last letter of the input.
func (s *Session) Backspace() {
	if s.Input != "" {
		s.Input = s.Input[:len(s.Input)-1]
	}
}

// Clear empties the input.
func (s *Session) Clear() { s.Input = "" }

// ShowHint shows the romaji of the active card, until it is answered, and
// marks the card hinted.
func (s *Session) ShowHint() {
	c := s.Active()
	if s.Phase != Playing || c == nil {
		return
	}
	s.Hint = true
	c.Hinted = true
}

// Stop ends a game without results and goes back to Setup.
func (s *Session) Stop() {
	s.Phase, s.Deck, s.Input, s.Hint = Setup, nil, "", false
}

// PlayAgain goes back to Setup from Results, keeping the settings.
func (s *Session) PlayAgain() { s.Stop() }
