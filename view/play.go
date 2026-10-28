package view

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"otokana/internal/deck"
	"otokana/internal/format"
	"otokana/internal/matcher"
)

// The motions of PRD section 4.9: the sheet torn off a correct answer
// turning over the binding, and the answer line shaking at a wrong one.
const (
	flipTime  = 300 * time.Millisecond
	shakeTime = 240 * time.Millisecond
)

func (a *App) play(c *ui.Context) {
	s := a.S
	now := c.Now()
	reduce := c.Preferences().ReduceMotion
	c.OnShortcut(0, ui.KeyEscape, s.Stop)
	// The screen takes the typing, with input methods off, so that a
	// Japanese one does not turn romaji into kana.
	ui.Column(c.Key("play")).Fill().Padding(32, 64, 32).Focusable().AutoFocus().FocusRing(false).
		InputMethod(false).TextCaret(ui.Rect{}).HandleInput(a.input).Children(func() {
		header(c, nil, func() {
			stop := ui.ButtonBase(c).Height(44).Padding(0, 8).Gap(8).Radius(6).TextColor(TeksRedup).
				Cursor(ui.CursorPointer).InputMethod(false).TextCaret(ui.Rect{}).HandleInput(a.input).OnClick(s.Stop)
			focusRing(stop, 6)
			stop.Children(func() {
				ui.Text(c, "Stop").FontSize(14)
				keycap(c, "Esc")
			})
		}).Opacity(0.45)
		ui.Column(c).Grow(1).Center().Gap(44).Children(func() {
			a.sheet(c, now, reduce)
			a.answer(c, now, reduce)
		})
		ui.Row(c).Justify(ui.Center).Gap(32).Opacity(0.7).Children(func() {
			ui.Row(c).Gap(8).Children(func() {
				keycap(c, "Backspace")
				ui.Text(c, "delete letter").FontSize(13).TextColor(TeksRedup)
			})
			ui.Row(c).Gap(8).Children(func() {
				keycap(c, "Ctrl")
				ui.Text(c, "+").FontSize(13).TextColor(TeksRedup)
				keycap(c, "Backspace")
				ui.Text(c, "clear").FontSize(13).TextColor(TeksRedup)
			})
		})
	})
}

// sheet is the pad with the active card on top, which a click turns to its
// hint. Like the other controls of the screen, it takes the typing too
// while it has the focus.
func (a *App) sheet(c *ui.Context, now time.Time, reduce bool) {
	s := a.S
	ui.Box(c).Size(320, 380).Children(func() {
		top := ui.ButtonBase(c).Label("Show hint").KeepFocus().Size(320, 380).Cursor(ui.CursorPointer).
			InputMethod(false).TextCaret(ui.Rect{}).HandleInput(a.input).OnClick(s.ShowHint)
		pad(top, sheetLayers(s.Remaining(), len(s.Deck)))
		focusRing(top, 3)
		top.Children(func() { face(c, s.Active(), s.Hint) })
		// Over the sheet, apart from its hint: the note the card plays at,
		// and its sound again.
		card := s.Active()
		if note := a.note(card); note != "" {
			ui.Text(c, note).Absolute().Top(12 + 20).Left(42).PassThrough().
				FontSize(14).FontWeight(500).TextColor(TintaRedup)
		}
		if a.sound(card.Entry.Romaji) != nil {
			again := ui.ButtonBase(c).Label("Play sound").KeepFocus().Size(44, 44).Radius(6).Center().
				Absolute().Top(12 + 6).Right(30).TextColor(TintaRedup).Cursor(ui.CursorPointer).
				InputMethod(false).TextCaret(ui.Rect{}).HandleInput(a.input).OnClick(func() { a.playSound(card) })
			focusRing(again, 6)
			again.Children(func() { ui.Icon(c, speaker).Size(20, 20).PassThrough() })
		}
		// The sheets torn off turn over the next one, the oldest on top.
		for i := len(a.torn) - 1; i >= 0; i-- {
			if age := now.Sub(a.torn[i].at); age < flipTime {
				a.flip(c, &a.torn[i], float32(age)/float32(flipTime), reduce)
				c.AnimationFrame()
			}
		}
	})
}

var flipEase = ui.CubicBezier(0.25, 0.1, 0.25, 1)

// flip draws a sheet torn off, k of the way through its turn: it swings on
// the binding toward the desk and up over it, shaded as it stands, its back
// showing once past upright, and fades out at the end. With less motion,
// it only fades.
func (a *App) flip(c *ui.Context, t *tornSheet, k float32, reduce bool) {
	// CSS's ease: it lifts at once, then slows as it goes over.
	turn, fade := 165*flipEase(k), 1-max(k-0.4, 0)/0.6
	if reduce {
		turn, fade = 0, 1-k
	}
	stand := float32(math.Sin(float64(turn) * math.Pi / 180)) // 0 lying flat, 1 upright
	// The shadow it casts on the next sheet, below the binding.
	if stand > 0 {
		shade := ui.LinearGradient{From: ui.RGBA(0, 0, 0, 0.22*stand*fade), To: ui.RGBA(0, 0, 0, 0), Angle: 180}
		ui.Box(c).Absolute().Top(12).Left(0).Size(320, 80*stand).PassThrough().
			Draw(func(p *ui.Painter, r ui.Rect) { p.FillGradient(r, shade, 0) })
	}
	ui.Box(c).Absolute().Top(12).Left(0).Size(320, 368).PassThrough().Opacity(fade).Children(func() {
		page := ui.Box(c).Absolute().Fill().PassThrough().Radius(0, 0, 3, 3).TransformOrigin(0.5, 0).RotateX(turn).Perspective(800).
			DrawOver(func(p *ui.Painter, r ui.Rect) { p.Fill(r, ui.RGBA(0, 0, 0, 0.2*stand), 0) })
		if turn >= 90 {
			page.Background(TepiKertas[0]) // its back, blank
			return
		}
		// Its face, as it showed on top.
		page.Background(Kertas).Children(func() {
			face(c, &t.card, t.hint).PassThrough()
			if note := a.note(&t.card); note != "" {
				ui.Text(c, note).Absolute().Top(20).Left(42).PassThrough().FontSize(14).FontWeight(500).TextColor(TintaRedup)
			}
			if a.sound(t.card.Entry.Romaji) != nil {
				ui.Box(c).Absolute().Top(6).Right(30).Size(44, 44).Center().PassThrough().TextColor(TintaRedup).
					Children(func() { ui.Icon(c, speaker).Size(20, 20).PassThrough() })
			}
		})
	})
}

// face is the practice box of a sheet, with the card's kana, or its
// romaji in green for the hint.
func face(c *ui.Context, card *deck.Card, hint bool) ui.Element {
	text, size := card.Kana(), float32(172)
	if len([]rune(text)) > 1 {
		size = 104 // yōon, two kana wide
	}
	if hint {
		text, size = card.Entry.Romaji, 96
	}
	b := practiceBox(c, 236, true, text, size).Absolute().Left(42).Top(64)
	if hint {
		b.TextColor(TintaHijau)
	}
	return b
}

// answer is the line the typing shows on, between the time and the cards
// left.
func (a *App) answer(c *ui.Context, now time.Time, reduce bool) {
	s := a.S
	ms := s.Elapsed(now).Milliseconds()
	c.After(time.Second - time.Duration(ms%1000)*time.Millisecond)
	line, text, caret := GarisJawaban, Teks, Hijau
	switch {
	case s.Result() == matcher.Wrong:
		line, text, caret = PenaMerah, PenaMerah, PenaMerah
	case s.Input != "":
		line = Hijau
	}
	ui.Row(c).Width(440).AlignItems(ui.End).Justify(ui.SpaceBetween).Children(func() {
		ui.Text(c, format.Clock(ms)).Width(80).Padding(0, 0, 10, 0).FontSize(16).TextColor(TeksRedup).FontFeatures("tnum")
		l := ui.Row(c).Size(220, 56).Justify(ui.Center).Gap(3).BorderWidth(0, 0, 2, 0).BorderColor(line)
		if t := now.Sub(a.wrongAt); !reduce && t < shakeTime {
			k := float64(t) / float64(shakeTime)
			l.Translate(float32(3*math.Sin(k*6*math.Pi)*(1-k)), 0)
			c.AnimationFrame()
		}
		l.Children(func() {
			if s.Input != "" {
				ui.Text(c, s.Input).FontSize(30).FontWeight(500).LetterSpacing(1.8).TextColor(text)
			}
			ui.Box(c).Size(2, 30).Background(caret)
		})
		ui.Row(c).Width(80).Padding(0, 0, 10, 0).Justify(ui.End).AlignItems(ui.Baseline).FontFeatures("tnum").Children(func() {
			ui.Textf(c, "%d", s.Remaining()).FontSize(16)
			ui.Textf(c, "/%d", len(s.Deck)).FontSize(13).TextColor(TeksRedup)
		})
	})
}
