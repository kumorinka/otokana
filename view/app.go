package view

import (
	"log"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"otokana/internal/audio"
	"otokana/internal/deck"
	"otokana/internal/format"
	"otokana/internal/notes"
	"otokana/internal/session"
	"otokana/internal/store"
)

// App is the state the window shows: the session, the import, and what
// animates.
type App struct {
	S *session.Session

	// What main.go connects to the system. Pick shows the dialog choosing
	// a voicebank's folder (dir) or zip file, and returns the path chosen,
	// "" when canceled; it blocks. Update runs fn on the main thread, then
	// draws a frame. Go runs fn in the background.
	Pick   func(dir bool) (string, error)
	Update func(fn func())
	Go     func(fn func())
	// Store keeps the settings, Prefs, and the imported voicebank, VB,
	// whose sounds Out plays.
	Store *store.Store
	Prefs store.Settings
	VB    *store.Voicebank
	Out   Player

	// The Settings and Import screens show over Setup and Results, the
	// Settings one over the Import one.
	settingsOpen bool
	prefsErr     error // why the settings were not saved
	importing    bool
	imp          importState

	// The sheets torn off by correct answers, turning over the binding,
	// the oldest first.
	torn []tornSheet
	// When the input last became wrong, for the answer line's shake.
	wrongAt time.Time
}

// tornSheet is a sheet torn off, as it showed: with its hint or not.
type tornSheet struct {
	card deck.Card
	hint bool
	at   time.Time
}

// Player plays sounds one after another: an audio.Out.
type Player interface{ Play(s ...*audio.Sound) }

// NewApp returns the app on Setup.
func NewApp() *App { return &App{S: session.New()} }

// Load reads the settings and the active voicebank from the store.
func (a *App) Load() error {
	var err error
	if a.Prefs, err = a.Store.Settings(); err != nil {
		return err
	}
	a.VB, err = a.Store.Active()
	return err
}

// savePrefs keeps the settings.
func (a *App) savePrefs() {
	a.prefsErr = a.Store.SaveSettings(a.Prefs)
}

// sound returns the sound of the kana of romaji, nil without one.
func (a *App) sound(romaji string) *audio.Sound { return a.VB.Sound(romaji) }

// note returns the name of the note the card plays at, "" when the
// voicebank's base note is not known (PRD section 7).
func (a *App) note(card *deck.Card) string {
	if a.VB == nil {
		return ""
	}
	base, ok := notes.Parse(a.VB.Base)
	if !ok {
		return ""
	}
	return notes.Name(base + card.Pitch)
}

// playSound plays the sounds of the cards' kana that have one, one after
// another, each at its card's note when it has one, else as sampled. Nil
// cards play nothing.
func (a *App) playSound(cards ...*deck.Card) {
	var sounds []*audio.Sound
	for _, card := range cards {
		if card == nil {
			continue
		}
		s := a.sound(card.Entry.Romaji)
		if s == nil {
			continue
		}
		if a.note(card) != "" {
			s = s.Shift(card.Pitch)
		}
		sounds = append(sounds, s)
	}
	if len(sounds) > 0 && a.Out != nil {
		a.Out.Play(sounds...)
	}
}

// start plays the deck that begin starts, with the sound of its first card
// when cards play as they show.
func (a *App) start(begin func(now time.Time), now time.Time) {
	begin(now)
	if a.S.Phase == session.Playing && a.Prefs.PlayOn != store.OnCorrect {
		a.playSound(a.S.Active())
	}
}

// record adds the session, just ended, to the history.
func (a *App) record() {
	s := a.S
	hinted := s.Hinted()
	r := store.Record{Date: s.Ended, Deck: format.DeckName(s.Settings), Cards: len(s.Deck), Hinted: hinted,
		Ms: s.Elapsed(s.Ended).Milliseconds(), Accuracy: format.Accuracy(len(s.Deck), hinted)}
	if err := a.Store.AddRecord(r); err != nil {
		log.Printf("history: %v", err)
	}
}

// View builds the screen of the session's phase, or the Settings or
// Import screen over it.
func (a *App) View(c *ui.Context) {
	c.SetTheme(theme)
	switch {
	case a.S.Phase == session.Playing:
		a.play(c)
	case a.settingsOpen:
		a.settingsScreen(c)
	case a.importing:
		a.importScreen(c)
	case a.S.Phase == session.Results:
		a.results(c)
	default:
		a.setup(c)
	}
}

// input takes the keys and text typed while playing.
func (a *App) input(ev ui.InputEvent) bool {
	switch ev.Kind {
	case ui.InputText:
		hint := a.S.Hint // before the answer hides it
		done, wrong := a.S.Type(ev.Text, ev.Time)
		if done != nil {
			a.torn = append(slices.DeleteFunc(a.torn, func(t tornSheet) bool { return ev.Time.Sub(t.at) >= flipTime }),
				tornSheet{*done, hint, ev.Time})
			// The card answered as it is, the next one as it shows.
			var cards []*deck.Card
			if a.Prefs.PlayOn != store.OnShow {
				cards = append(cards, done)
			}
			if a.Prefs.PlayOn != store.OnCorrect {
				cards = append(cards, a.S.Active())
			}
			a.playSound(cards...)
			if a.S.Phase == session.Results {
				a.record()
			}
		}
		if wrong {
			a.wrongAt = ev.Time
		}
		return true
	case ui.InputKeyDown:
		if ev.Key != ui.KeyBackspace {
			return false
		}
		if ev.Mods&(ui.Ctrl|ui.Alt|ui.Super) != 0 {
			a.S.Clear()
		} else {
			a.S.Backspace()
		}
		return true
	}
	return false
}
