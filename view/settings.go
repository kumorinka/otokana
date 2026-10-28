package view

import (
	"github.com/egoist/mygo/ui"

	"otokana/internal/store"
)

var sliders = ui.MustParseSVG([]byte(`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M4 7h10M18 7h2M4 17h2M10 17h10"/><circle cx="16" cy="7" r="2"/><circle cx="8" cy="17" r="2"/></svg>`))

// settingsButton opens the Settings screen, at the right of the header.
func (a *App) settingsButton(c *ui.Context) {
	b := ui.ButtonBase(c).Label("Settings").Size(44, 44).Radius(6).Center().TextColor(TeksRedup).
		Cursor(ui.CursorPointer).OnClick(func() { a.settingsOpen = true })
	focusRing(b, 6)
	b.Children(func() { ui.Icon(c, sliders).Size(20, 20) })
}

// settingsScreen is a column of the settings, saved as they change (PRD
// section 4.9, without a mockup).
func (a *App) settingsScreen(c *ui.Context) {
	p := &a.Prefs
	ui.Column(c.Key("settings")).Fill().Padding(32, 64, 40).Children(func() {
		header(c, func() { a.settingsOpen = false }, nil)
		ui.Column(c).Grow(1).Width(560).Justify(ui.Center).Gap(36).Children(func() {
			pageTitle(c, "Settings", 44)
			ui.Column(c).Gap(14).Children(func() {
				ui.Text(c, "When to play sounds").FontSize(14).FontWeight(500).TextColor(TeksRedup)
				ui.RadioGroup(c, func() {
					for _, o := range []struct {
						on            store.PlayOn
						label, detail string
					}{
						{store.OnCorrect, "When answered right", "The kana you answer plays, as a sign it was right."},
						{store.OnShow, "When a card shows", "Hear the kana first, then type its romaji."},
						{store.OnBoth, "Both", "The kana you answer plays, then the next one."},
					} {
						choice(c, ui.RadioBase(c, &p.PlayOn, o.on).OnChange(a.savePrefs), p.PlayOn == o.on, 28, "", 0, o.label, o.detail)
					}
				}).Label("When to play sounds").Gap(6)
			})
			if a.VB != nil {
				baseNote(c, a.VB)
			} else {
				ui.Column(c).Gap(8).Children(func() {
					ui.Text(c, "Base note").FontSize(14).FontWeight(500).TextColor(TeksRedup)
					ui.Row(c).Gap(4).TextColor(TeksRedup).Children(func() {
						ui.Text(c, "No voicebank yet.").FontSize(15)
						link(c, "Import", 15, 44, 6).OnClick(func() { a.settingsOpen, a.importing = false, true })
					})
				})
			}
			if a.prefsErr != nil {
				ui.Text(c, "Settings not saved: "+a.prefsErr.Error()).FontSize(14).TextColor(PenaMerah)
			}
		})
	})
}
