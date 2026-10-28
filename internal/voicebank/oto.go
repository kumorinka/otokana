// Package voicebank: decode teks, parser oto.ini, loader folder dan zip, resolver.
//
// Go murni, tanpa import mygo.
package voicebank

import (
	"bytes"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
)

// Decode returns text as UTF-8 when it is valid UTF-8, without its byte
// order mark, and decoded from Shift-JIS, as many Japanese voicebanks
// have it, otherwise.
func Decode(b []byte) string {
	b = bytes.TrimPrefix(b, []byte("\xEF\xBB\xBF"))
	if utf8.Valid(b) {
		return string(b)
	}
	s, err := japanese.ShiftJIS.NewDecoder().Bytes(b)
	if err != nil {
		return strings.ToValidUTF8(string(b), "�")
	}
	return string(s)
}

// Entry is a line of an oto.ini: a sample, in milliseconds of its file.
type Entry struct {
	Alias string
	File  string // slash path of the wav in the source
	// Offset is where the sample starts. Cutoff is where it ends: that far
	// before the end of the file when positive, and its length from Offset
	// when negative.
	Offset, Consonant, Cutoff, Preutterance, Overlap float64
}

// Duration returns how long the sample plays in a file fileMs long: -Cutoff
// when negative, what is left of the file between Offset and Cutoff
// otherwise, within the file.
func (e Entry) Duration(fileMs float64) float64 {
	left := max(fileMs-max(e.Offset, 0), 0)
	d := left - e.Cutoff
	if e.Cutoff < 0 {
		d = -e.Cutoff
	}
	return min(max(d, 0), left)
}

// ParseOto reads the lines of an oto.ini, file.wav=alias,offset,consonant,
// cutoff,preutterance,overlap, whose files are in dir. An empty alias is
// the file's name without its extension, and missing numbers are 0, as
// OpenUtau reads them. It returns the entries and how many lines it could
// not read; blank lines and those starting with # are skipped.
func ParseOto(text, dir string) (entries []Entry, bad int) {
	for line := range strings.Lines(text) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		wav, rest, ok := strings.Cut(line, "=")
		wav = strings.TrimSpace(strings.ReplaceAll(wav, `\`, "/"))
		if !ok || wav == "" {
			bad++
			continue
		}
		fields := strings.Split(rest, ",")
		e := Entry{Alias: strings.TrimSpace(fields[0]), File: path.Join(dir, wav)}
		if e.Alias == "" {
			e.Alias = strings.TrimSuffix(path.Base(wav), path.Ext(wav))
		}
		nums := []*float64{&e.Offset, &e.Consonant, &e.Cutoff, &e.Preutterance, &e.Overlap}
		for i, n := range nums {
			if i+1 >= len(fields) {
				break
			}
			s := strings.TrimSpace(fields[i+1])
			if s == "" {
				continue
			}
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				ok = false
				break
			}
			*n = v
		}
		if !ok {
			bad++
			continue
		}
		entries = append(entries, e)
	}
	return entries, bad
}

// parseCharacter reads the name and the author and web site of a
// character.txt, lines of key=value or key:value (Kasane Teto's has
// web:), or of a character.yaml, key: value only (yaml true).
func parseCharacter(text string, yaml bool) (name, author, web string) {
	for line := range strings.Lines(text) {
		seps := "=:"
		if yaml {
			seps = ":"
		}
		i := strings.IndexAny(line, seps)
		if i < 0 || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			continue
		}
		k, v := line[:i], line[i+1:]
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "name":
			name = v
		case "author":
			author = v
		case "web":
			web = v
		}
	}
	return name, author, web
}
