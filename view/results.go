package view

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"otokana/internal/format"
)

// results shows how many cards were answered without their hint, of all.
func (a *App) results(c *ui.Context) {
	s := a.S
	cards, hinted := len(s.Deck), s.Hinted()
	c.OnShortcut(0, ui.KeyEnter, s.PlayAgain)
	ui.Column(c.Key("results")).Fill().Padding(32, 64, 40).Children(func() {
		header(c, nil, func() { a.settingsButton(c) })
		ui.Column(c).Grow(1).MaxWidth(640).Justify(ui.Center).Gap(44).Children(func() {
			ui.Column(c).Gap(16).Children(func() {
				pageTitle(c, "Finished", 48)
				ui.Row(c).AlignItems(ui.Baseline).Gap(14).Children(func() {
					kanaText(c, fmt.Sprintf("%d/%d", cards-hinted, cards), 96).FontFeatures("tnum")
					ui.Text(c, "correct").FontSize(20)
				})
				ui.Text(c, fmt.Sprintf("%s · %s · %d%%", format.DeckName(s.Settings),
					format.Duration(s.Elapsed(c.Now()).Milliseconds()), format.Accuracy(cards, hinted))).
					FontSize(16).TextColor(TeksRedup)
			})
			ui.Row(c).Children(func() { enterHint(c, "Play again", "Press", "to play again", s.PlayAgain) })
		})
	})
}
