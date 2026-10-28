package voicebank

import (
	"cmp"
	"errors"
	"path"
	"slices"
	"strings"
)

// Voicebank is what the oto.ini files and the character file of a source
// say.
type Voicebank struct {
	Name   string
	Credit string // the author, and their web site, from the character file
	Otos   int    // oto.ini files read
	Bad    int    // their lines that could not be read
	// Entries holds every line read, of the shallowest oto.ini first, the
	// same alias more than once as OpenUtau keeps them.
	Entries []Entry
}

// ErrNoOto is returned for sources without an oto.ini.
var ErrNoOto = errors.New("no oto.ini inside")

// Read reads the oto.ini files of src, in its subfolders too, and its
// character.txt, and character.yaml over it, the shallowest of each.
// progress, when not nil, gets the oto.ini files read so far and their
// number.
func Read(src *Source, progress func(done, total int)) (*Voicebank, error) {
	byDepth := func(a, b string) int {
		return cmp.Or(cmp.Compare(strings.Count(a, "/"), strings.Count(b, "/")), strings.Compare(a, b))
	}
	named := func(base string) []string {
		var out []string
		for _, f := range src.Files {
			if strings.EqualFold(path.Base(f), base) {
				out = append(out, f)
			}
		}
		slices.SortFunc(out, byDepth)
		return out
	}
	otos := named("oto.ini")
	if len(otos) == 0 {
		return nil, ErrNoOto
	}
	vb := &Voicebank{Name: src.Name}
	var author, web string
	for _, c := range []struct {
		file string
		yaml bool
	}{{"character.txt", false}, {"character.yaml", true}} {
		files := named(c.file)
		if len(files) == 0 {
			continue
		}
		b, err := src.ReadFile(files[0])
		if err != nil {
			return nil, err
		}
		n, a, w := parseCharacter(Decode(b), c.yaml)
		vb.Name, author, web = cmp.Or(n, vb.Name), cmp.Or(a, author), cmp.Or(w, web)
	}
	vb.Credit = strings.Join(slices.DeleteFunc([]string{author, web}, func(s string) bool { return s == "" }), ", ")
	for i, o := range otos {
		b, err := src.ReadFile(o)
		if err != nil {
			return nil, err
		}
		entries, bad := ParseOto(Decode(b), path.Dir(o))
		vb.Entries = append(vb.Entries, entries...)
		vb.Otos++
		vb.Bad += bad
		if progress != nil {
			progress(i+1, len(otos))
		}
	}
	return vb, nil
}

// Load opens the folder or zip file at p, reads it and closes it.
func Load(p string, progress func(done, total int)) (*Voicebank, error) {
	src, err := OpenSource(p)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return Read(src, progress)
}
