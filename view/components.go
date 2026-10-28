package view

import "github.com/egoist/mygo/ui"

// kanaText is kana, or hint romaji, in the handwriting font.
func kanaText(c *ui.Context, s string, size float32) ui.Element {
	return ui.Text(c, s).Font(FontKana).FontWeight(600).FontSize(size).LineHeight(1)
}

// practiceBox is a box of a Japanese writing workbook: a frame, of paper
// when on, see-through and dim when off.
func practiceBox(c *ui.Context, size float32, on bool, kana string, fontSize float32) ui.Element {
	frame := GarisKotak
	b := ui.Box(c).Size(size, size).Center()
	if on {
		b.Background(Kertas).TextColor(Tinta)
	} else {
		frame = GarisKotakRedup
	}
	b.Border(1.5, frame)
	// Attached, the text stays centered when its lines, the font's own
	// height, are taller than the box, which Center would align at the top.
	// The kana passes the pointer to the box, and to what holds it.
	return b.Children(func() { kanaText(c, kana, fontSize).Attach(ui.AnchorCenter, ui.AnchorCenter).PassThrough() })
}

// standInBox is the practice box of a kana that plays another's sound, as
// ぢ plays じ: paper, its frame dashed.
func standInBox(c *ui.Context, size float32, kana string, fontSize float32) ui.Element {
	return practiceBox(c, size, true, kana, fontSize).BorderStyle(ui.BorderDashed)
}

// pageTitle is a page's title. MyGo makes lines no shorter than the font's
// own height, 1.45 em for Klee One; the negative margins give it the 1.2 em
// the mockups have.
func pageTitle(c *ui.Context, s string, size float32) ui.Element {
	return kanaText(c, s, size).Margin(-0.125*size, 0)
}

// link is a button of underlined green text.
func link(c *ui.Context, label string, size, height, padX float32) ui.Element {
	b := ui.ButtonBase(c).Height(height).Padding(0, padX).Radius(4).Cursor(ui.CursorPointer)
	focusRing(b, 4)
	return b.Children(func() { ui.Text(c, label).FontSize(size).FontWeight(500).TextColor(Hijau).Underline() })
}

// sheetLayers returns how many sheet edges show under the top one: 5 for
// a full deck, thinning out to none on the last card.
func sheetLayers(remaining, total int) int {
	if remaining <= 1 || total <= 1 {
		return 0
	}
	n := (5*(remaining-1) + total - 2) / (total - 1) // ceil
	return min(n, remaining-1, 5)
}

// pad styles e as the pad of paper: the binding strip on top, the edges of
// layers sheets below it, and its shadow on the desk.
func pad(e ui.Element, layers int) ui.Element {
	// MyGo paints later shadows over earlier ones, unlike CSS: the drop
	// shadow first, then the edges from the deepest up.
	e.Radius(3).Background(Kertas).BorderWidth(12, 0, 0, 0).BorderColor(Pengikat).Shadow(0, 30, 40, -12, BayanganBlok)
	for i := layers - 1; i >= 0; i-- {
		k := float32(i + 1)
		e.Shadow(0, 3*k, 0, -k, TepiKertas[i])
	}
	return e
}

// focusRing replaces MyGo's ring on e with the design's: 2 DIPs of green,
// 3 DIPs away.
func focusRing(e ui.Element, radius float32) ui.Element {
	visible := e.FocusVisible()
	return e.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if visible {
			p.Stroke(ui.Rect{X: r.X - 5, Y: r.Y - 5, W: r.W + 10, H: r.H + 10}, Hijau, radius+5, 2)
		}
	})
}

// keycap shows a key of a shortcut hint.
func keycap(c *ui.Context, key string) ui.Element {
	return ui.Row(c).Height(24).Padding(0, 8).Radius(4).BorderWidth(1, 1, 2, 1).BorderColor(GarisMeja).
		Children(func() { ui.Text(c, key).FontSize(12).TextColor(TeksKeycap) })
}

// primaryButton is a button of paper.
func primaryButton(c *ui.Context, label string) ui.Element {
	b := ui.ButtonBase(c).Height(52).Padding(0, 28).Radius(6).Background(Kertas).Cursor(ui.CursorPointer)
	focusRing(b, 6)
	return b.Children(func() { ui.Text(c, label).FontSize(17).FontWeight(700).TextColor(Tinta) })
}

// secondaryButton is underlined text.
func secondaryButton(c *ui.Context, label string) ui.Element {
	b := ui.ButtonBase(c).Height(52).Padding(0, 10).Radius(6).Cursor(ui.CursorPointer)
	focusRing(b, 6)
	return b.Children(func() { ui.Text(c, label).FontSize(17).FontWeight(500).Underline() })
}

var chevronLeft = ui.MustParseSVG([]byte(`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M15 5l-7 7 7 7"/></svg>`))

// header is the logo, after a back button when back is not nil, with what
// right builds at the other end.
func header(c *ui.Context, back func(), right func()) ui.Element {
	return ui.Row(c).Height(44).Justify(ui.SpaceBetween).Children(func() {
		ui.Row(c).Gap(6).Children(func() {
			if back != nil {
				b := ui.ButtonBase(c).Label("Back").Size(44, 44).Margin(0, 0, 0, -10).Radius(6).Center().
					TextColor(TeksRedup).Cursor(ui.CursorPointer).OnClick(back)
				focusRing(b, 6)
				b.Children(func() { ui.Icon(c, chevronLeft).Size(20, 20) })
			}
			ui.Row(c).Gap(10).Children(func() {
				practiceBox(c, 30, true, "音", 19)
				ui.Text(c, "OtoKana").Font(FontKana).FontWeight(600).FontSize(22)
			})
		})
		if right != nil {
			right()
		}
	})
}

// enterHint is a hint of what Enter does, as "Press [Enter] to start", and
// a button labeled label doing it too, for a mouse or a finger.
func enterHint(c *ui.Context, label, before, after string, fn func()) ui.Element {
	b := ui.ButtonBase(c).Label(label).Height(44).Padding(0, 10).Margin(0, 0, 0, -10).Gap(8).Radius(6).
		Cursor(ui.CursorPointer).TextColor(TeksRedup).OnClick(fn)
	if b.Hovered() {
		b.TextColor(Teks)
	}
	focusRing(b, 6)
	return b.Children(func() {
		ui.Text(c, before).FontSize(14)
		keycap(c, "Enter")
		ui.Text(c, after).FontSize(14)
	})
}
