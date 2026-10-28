// Package view holds OtoKana's screens and components.
package view

import "github.com/egoist/mygo/ui"

// The palette of PRD section 4.9. It is fixed: it follows neither the
// system's light or dark appearance nor its accent color.
var (
	Meja       = ui.Hex("#22201C") // window background
	Kertas     = ui.Hex("#EFEBE0") // sheets, selected boxes, primary button
	Tinta      = ui.Hex("#2A2723") // kana and text on paper
	TintaRedup = ui.Hex("#6B665C") // note label and icons on paper
	Teks       = ui.Hex("#E9E4D8") // main text on the desk
	TeksRedup  = ui.Hex("#A39C8E") // secondary text, group labels, icons
	GarisKotak = ui.Hex("#6E8F72") // frame of a practice box
	Hijau      = ui.Hex("#9DB89F") // answer line while typing, caret, links, focus ring
	TintaHijau = ui.Hex("#4F7A56") // hint romaji on paper, large text only
	PenaMerah  = ui.Hex("#E8705F") // typed text and answer line when wrong
	Pengikat   = ui.Hex("#3F5A45") // binding strip on top of the pad

	// Desk lines.
	GarisKotakRedup = ui.Hex("#4E4A43") // frame of an unselected box
	GarisBantuRedup = ui.Hex("#3E3B35") // guide lines of an unselected box
	GarisMeja       = ui.Hex("#57524A") // keycaps, import area
	GarisJawaban    = ui.Hex("#5B564D") // neutral answer line

	// From the mockups: keycap text, and the edges of the sheets under the
	// top one, lighter near it, with the pad's drop shadow.
	TeksKeycap = ui.Hex("#CFC9BC")
	// The guide lines of a box whose sound stands in for another, and the
	// line between the Import screen's columns.
	GarisPemisah = ui.Hex("#3A3731")
	TepiKertas   = []ui.Color{ui.Hex("#DAD4C6"), ui.Hex("#CDC6B6"), ui.Hex("#C0B8A6"), ui.Hex("#B3AA97"), ui.Hex("#A69D89")}
	BayanganBlok = ui.RGBA(0, 0, 0, 0.55)
)

var theme = Theme()

// Font families, which main.go registers from assets/fonts.
const (
	FontKana = "Klee One"            // 600: kana, page titles, logo, hint romaji
	FontUI   = "Zen Kaku Gothic New" // 400, 500, 700: all interface text
)

// Theme returns the theme every frame sets, so that what widgets draw
// keeps to the palette.
func Theme() *ui.Theme {
	t := ui.DarkTheme()
	t.Background = Meja
	t.Surface, t.SurfaceHover, t.SurfacePressed = Meja, GarisBantuRedup, GarisKotakRedup
	t.Border = GarisKotakRedup
	t.Text, t.TextMuted = Teks, TeksRedup
	t.Accent, t.AccentHover, t.AccentPressed, t.AccentText = Kertas, Kertas, Kertas, Tinta
	t.Danger = PenaMerah
	t.Success = Hijau
	t.Focus = Hijau
	t.Selection = ui.RGBA(Hijau.R, Hijau.G, Hijau.B, 0.35)
	t.Inverse, t.InverseText = Kertas, Tinta
	t.Font = FontUI
	t.FontSize = 15
	return t
}
