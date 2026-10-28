package view

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"otokana/internal/kana"
	"otokana/internal/notes"
	"otokana/internal/session"
	"otokana/internal/store"
)

var upload = ui.MustParseSVG([]byte(`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M12 16V5M7.5 9.5L12 5l4.5 4.5"/><path d="M5 19h14"/></svg>`))

// importState is where an import is.
type importState struct {
	picking     bool // a dialog shows
	busy        bool // importing
	done, total int  // samples copied, of total
	err         error
}

// pick shows the dialog choosing a voicebank's folder (dir) or zip file,
// and imports the one chosen.
func (a *App) pick(dir bool) {
	if a.imp.picking || a.imp.busy {
		return
	}
	a.imp.picking = true
	a.Go(func() {
		p, err := a.Pick(dir)
		a.Update(func() {
			a.imp.picking = false
			switch {
			case err != nil:
				a.imp = importState{err: err}
			case p != "":
				a.load(p)
			}
		})
	})
}

// load imports the voicebank at p in the background, which becomes the
// active one.
func (a *App) load(p string) {
	if a.imp.busy {
		return
	}
	a.imp = importState{busy: true}
	a.Go(func() {
		vb, err := a.Store.Import(p, func(done, total int) {
			a.Update(func() { a.imp.done, a.imp.total = done, total })
		})
		a.Update(func() {
			a.imp = importState{err: err}
			if err == nil {
				a.VB, a.Prefs.Voicebank = vb, vb.ID
				a.savePrefs()
			}
		})
	})
}

func (a *App) importScreen(c *ui.Context) {
	closeImport := func() { a.importing = false }
	ui.Column(c.Key("import")).Fill().Padding(32, 64, 40).Children(func() {
		header(c, closeImport, func() { a.settingsButton(c) })
		ui.Row(c).Grow(1).Gap(56).Children(func() {
			ui.Column(c).Width(560).Gap(32).Children(func() {
				ui.Column(c).Gap(12).Children(func() {
					pageTitle(c, "Import a voicebank", 44)
					ui.Text(c, "Choose a UTAU voicebank folder or zip file. Its files are read on this device and are not uploaded anywhere.").
						FontSize(16).LineHeight(1.7).TextColor(TeksRedup).MaxWidth(480)
				})
				a.dropZone(c)
				ui.Text(c, "Voicebanks are not bundled with the app, so use one you have the right to use. Without a voicebank, practice still works, just without sound.").
					FontSize(14).LineHeight(1.7).TextColor(TeksRedup).MaxWidth(500)
				ui.Row(c).Children(func() {
					link(c, "Skip, practice without sound", 16, 52, 10).OnClick(closeImport)
				})
			})
			if a.VB != nil {
				ui.Column(c).Grow(1).Gap(20).Padding(0, 0, 0, 56).BorderWidth(0, 0, 0, 1).BorderColor(GarisPemisah).
					Children(func() { a.activeVoicebank(c) })
			}
		})
	})
}

// dropZone takes a folder or zip dropped on it, or chosen with its
// buttons, and shows how importing it goes.
func (a *App) dropZone(c *ui.Context) {
	zone := ui.Column(c).Height(220).Radius(6).Border(1.5, GarisMeja).BorderStyle(ui.BorderDashed).Center().Gap(18)
	if files := zone.DroppedFiles(); len(files) > 0 {
		a.load(files[0])
	}
	if zone.FileDragOver() {
		zone.BorderColor(Hijau)
	}
	im := a.imp
	title, detail := "Drag a folder or .zip file here", "It must contain oto.ini and wav files"
	switch {
	case im.busy:
		title, detail = "Importing voicebank…", ""
		if im.total > 0 {
			detail = fmt.Sprintf("Copying sample %d of %d", im.done, im.total)
		}
	case im.err != nil:
		title, detail = "The voicebank could not be imported", sentence(im.err.Error())
	}
	zone.Children(func() {
		if !im.busy {
			ui.Icon(c, upload).Size(28, 28).TextColor(TeksRedup)
		}
		ui.Column(c).AlignItems(ui.Center).Gap(6).Padding(0, 24).Children(func() {
			ui.Text(c, title).FontSize(17).FontWeight(500).TextAlign(ui.Center)
			if detail != "" {
				ui.Text(c, detail).FontSize(14).TextColor(TeksRedup).TextAlign(ui.Center)
			}
		})
		if !im.busy {
			ui.Row(c).Gap(12).Children(func() {
				primaryButton(c, "Choose folder").OnClick(func() { a.pick(true) })
				secondaryButton(c, "Choose zip file").OnClick(func() { a.pick(false) })
			})
		}
	})
}

// activeVoicebank shows the active voicebank, and the sheet of what its
// kana play: their own sample on paper, another's in a dashed frame, and
// none in a dim box.
func (a *App) activeVoicebank(c *ui.Context) {
	vb := a.VB
	ui.Column(c).Gap(6).Children(func() {
		ui.Text(c, "Active voicebank").FontSize(14).FontWeight(500).TextColor(TeksRedup)
		kanaText(c, vb.Name, 26).LineHeight(1.3).MaxLines(2)
		if vb.Credit != "" {
			ui.Text(c, "Credit: "+vb.Credit).FontSize(14).TextColor(TeksRedup)
		}
	})
	baseNote(c, vb)
	standIns, missing := 0, 0
	ui.Row(c).Width(15*28 + 14*4).Wrap().Gap(4).Children(func() {
		for i := range kana.Table {
			k := &kana.Table[i]
			size := float32(15)
			if len([]rune(k.Hiragana)) > 1 {
				size = 11
			}
			s, ok := vb.Samples[k.Romaji]
			switch {
			case !ok:
				missing++
				practiceBox(c, 28, false, k.Hiragana, size).TextColor(TeksRedup)
			case !s.Own(k.Romaji):
				standIns++
				standInBox(c, 28, k.Hiragana, size)
			default:
				practiceBox(c, 28, true, k.Hiragana, size)
			}
		}
	})
	note := func(s string) ui.Element { return ui.Text(c, s).FontSize(14).LineHeight(1.6).TextColor(TeksRedup) }
	ui.Column(c).Gap(8).Children(func() {
		ui.Row(c).Children(func() {
			ui.Text(c, fmt.Sprintf("%d of %d kana", vb.Own(), len(kana.Table))).FontSize(14).LineHeight(1.6)
			note(" have their own sample.")
		})
		if standIns > 0 {
			ui.Row(c).Gap(10).Children(func() {
				standInBox(c, 20, "ぢ", 11)
				note("Dashed: a stand-in sound, as ぢ uses じ's.")
			})
		}
		if missing > 0 {
			ui.Row(c).Gap(10).Children(func() {
				practiceBox(c, 20, false, "", 11)
				note(fmt.Sprintf("Dim: %d kana without sound.", missing))
			})
		}
		if n := len(vb.Skipped); n > 0 {
			note(fmt.Sprintf("%d samples skipped, such as %s.", n, vb.Skipped[0])).MaxWidth(476)
		}
	})
}

// baseField is what the base note's field holds: the text typed, and why
// it was not saved.
type baseField struct {
	text string
	err  error
}

// baseNote is the field telling the note vb's samples are sung at, which
// the cards are labeled and pitched by. It saves the note as it is typed;
// a text that is not a note leaves it unknown.
func baseNote(c *ui.Context, vb *store.Voicebank) {
	box := ui.Column(c).MaxWidth(476)
	f := ui.Local(box, vb.ID, func() baseField { return baseField{text: vb.Base} })
	msg := ""
	if _, ok := notes.Parse(f.text); !ok && strings.TrimSpace(f.text) != "" {
		msg = "Not a note name, such as E4 or F#3."
	}
	if f.err != nil {
		msg = "Not saved: " + f.err.Error()
	}
	box.Children(func() {
		ui.Field(c, "Base note", func() {
			// The message beside the input, not below it, keeps the panel from
			// moving as the note is typed.
			ui.Row(c).Gap(14).Children(func() {
				in := ui.TextInput(c, &f.text).Width(88).Placeholder("e.g. E4").InputMethod(false).FontSize(15).TextColor(Teks).
					OnChange(func() {
						base := ""
						if midi, ok := notes.Parse(f.text); ok {
							base = notes.Name(midi)
						}
						if base != vb.Base {
							vb.Base = base
							f.err = vb.Save()
						}
					})
				if msg != "" {
					in.Error(msg)
					ui.Text(c, msg).TextColor(PenaMerah).Shrink(1)
				}
			})
		}).FontSize(14).TextColor(TeksRedup).
			Description(fmt.Sprintf("The note %s's samples are sung at. Once set, each card plays up to %d semitones above or below it, with its note name in the card's corner.", vb.Name, session.PitchRange))
	})
}

// sentence capitalizes s and ends it with a period.
func sentence(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToUpper(s[:1]) + s[1:]
	if !strings.HasSuffix(s, ".") {
		s += "."
	}
	return s
}
