package view

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"otokana/internal/deck"
	"otokana/internal/kana"
)

var speaker = ui.MustParseSVG([]byte(`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M4 9v6h4l5 4V5L8 9H4z"/><path d="M16.5 8.5a5 5 0 0 1 0 7"/></svg>`))

func (a *App) setup(c *ui.Context) {
	s := &a.S.Settings
	c.OnShortcut(0, ui.KeyEnter, func() { a.start(a.S.Start, c.Now()) })
	w, _ := c.Size()
	ui.Column(c.Key("setup")).Fill().Padding(32, 64, 40).Children(func() {
		header(c, nil, func() { a.settingsButton(c) })
		ui.Row(c).Grow(1).Gap(96).Children(func() {
			ui.Column(c).Width(520).Gap(36).Children(func() {
				ui.Column(c).Gap(12).Children(func() {
					pageTitle(c, "Kana practice", 48)
					ui.Text(c, "Type the romaji of each kana that shows up.").
						FontSize(16).LineHeight(1.7).TextColor(TeksRedup).MaxWidth(440)
				})
				group(c, "Scripts", func() {
					option(c, &s.Hiragana, "か", 32, "Hiragana", "")
					option(c, &s.Katakana, "カ", 32, "Katakana", "")
				})
				group(c, "Variants", func() {
					option(c, &s.Monographs, "か", 32, "Basic", "46 kana")
					option(c, &s.Diacritics, "が", 32, "Dakuten", "25 kana")
					option(c, &s.Digraphs, "きゃ", 24, "Yōon", "21 kana")
				})
				ui.Row(c).Height(44).Gap(10).TextColor(TeksRedup).Children(func() {
					ui.Icon(c, speaker).Size(18, 18)
					label := "Import"
					if vb := a.VB; vb != nil {
						label = "Change"
						ui.Row(c).Shrink(1).Children(func() {
							ui.Text(c, vb.Name).FontSize(15).TextColor(Teks).SingleLine().Shrink(1)
							ui.Text(c, fmt.Sprintf(", %d of %d kana with sound", len(vb.Samples), len(kana.Table))).FontSize(15).Shrink(0)
						})
					} else {
						ui.Text(c, "No voicebank yet, practicing without sound.").FontSize(15)
					}
					link(c, label, 15, 44, 6).OnClick(func() { a.importing = true })
				})
				ui.Row(c).Children(func() {
					start := enterHint(c, "Start practice", "Press", "to start", func() { a.start(a.S.Start, c.Now()) })
					if !a.S.CanStart() {
						start.Disabled(true).Opacity(0.4)
					}
				})
			})
			// The preview, which narrow windows leave out.
			ui.Column(c).Grow(1).AlignItems(ui.Center).Gap(40).Hidden(w < 1180).Children(func() {
				preview(c, *s)
			})
		})
	})
}

// group is a labeled row of options.
func group(c *ui.Context, label string, options func()) {
	ui.Column(c).Gap(14).Children(func() {
		ui.Text(c, label).FontSize(14).FontWeight(500).TextColor(TeksRedup)
		ui.Row(c).Gap(28).Children(options)
	})
}

// option toggles a setting, shown as a practice box beside its label.
func option(c *ui.Context, on *bool, sample string, sampleSize float32, label, count string) {
	choice(c, ui.ToggleBase(c, on), *on, 56, sample, sampleSize, label, count)
}

// choice styles e, a control that is on or off, as a practice box of size
// with sample in it, on paper when on, beside its label and detail.
func choice(c *ui.Context, e ui.Element, on bool, size float32, sample string, sampleSize float32, label, detail string) {
	e.Gap(14).Padding(4, 10, 4, 4).Radius(4).Cursor(ui.CursorPointer)
	if on {
		e.TextColor(Teks)
	} else {
		e.TextColor(TeksRedup)
	}
	focusRing(e, 4)
	e.Children(func() {
		practiceBox(c, size, on, sample, sampleSize)
		ui.Column(c).Gap(2).Shrink(1).Children(func() {
			ui.Text(c, label).FontSize(17).FontWeight(500)
			if detail != "" {
				ui.Text(c, detail).FontSize(13).TextColor(TeksRedup)
			}
		})
	})
}

// preview shows the pad with a kana of the choice, and the deck's size.
func preview(c *ui.Context, s deck.Settings) {
	cards := deck.Build(s)
	sample := ""
	if len(cards) > 0 {
		// か, が or きゃ, of the first script and variant chosen.
		romaji := map[kana.Group]string{kana.Mono: "ka", kana.Dia: "ga", kana.Di: "kya"}[s.Groups()[0]]
		for _, c := range cards {
			if c.Entry.Romaji == romaji {
				sample = c.Kana()
				break
			}
		}
	}
	size := float32(160)
	if len([]rune(sample)) > 1 {
		size = 100
	}
	pad(ui.Box(c).Size(300, 344).Center(), sheetLayers(len(cards), len(cards))).Children(func() {
		practiceBox(c, 220, true, sample, size)
	})
	caption := fmt.Sprintf("%d cards, shuffled", len(cards))
	if len(cards) == 0 {
		caption = "Choose at least one script and one variant"
	}
	ui.Text(c, caption).FontSize(14).TextColor(TeksRedup)
}
