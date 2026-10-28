// Package notes: nama nada, nomor MIDI, dan frekuensi.
//
// Go murni, tanpa import mygo.
package notes

import (
	"math"
	"strconv"
	"strings"
)

var names = [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// Name returns the name of the MIDI note, with sharps, as UTAU writes it:
// E4 for 64, C#4 for 61, C-1 for 0.
func Name(midi int) string {
	octave, step := midi/12, midi%12
	if step < 0 { // below C-1, a pitch shift away from it
		octave, step = octave-1, step+12
	}
	return names[step] + strconv.Itoa(octave-1)
}

// Parse returns the MIDI note of a name as Name gives, in either case,
// with # or ♯ for a sharp and b or ♭ for a flat, and whether it is a note
// of 0 to 127.
func Parse(name string) (int, bool) {
	s := strings.TrimSpace(name)
	if s == "" {
		return 0, false
	}
	letter := strings.IndexByte("CDEFGAB", strings.ToUpper(s[:1])[0])
	if letter < 0 {
		return 0, false
	}
	midi := [7]int{0, 2, 4, 5, 7, 9, 11}[letter]
	s = s[1:]
	for _, acc := range []struct {
		sign string
		step int
	}{{"#", 1}, {"♯", 1}, {"b", -1}, {"♭", -1}} {
		if rest, ok := strings.CutPrefix(s, acc.sign); ok {
			midi, s = midi+acc.step, rest
			break
		}
	}
	octave, err := strconv.Atoi(s)
	midi += (octave + 1) * 12
	if err != nil || midi < 0 || midi > 127 {
		return 0, false
	}
	return midi, true
}

// Freq returns the frequency of the MIDI note in Hz, A4 (69) being 440.
func Freq(midi int) float64 { return 440 * math.Pow(2, float64(midi-69)/12) }
