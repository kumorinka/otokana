// Package store: pengaturan, riwayat, dan sampel di folder data yang diberikan main.go.
//
// Go murni, tanpa import mygo.
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"otokana/internal/audio"
	"otokana/internal/kana"
	"otokana/internal/voicebank"
)

// Store keeps the app's data in Dir: settings.json, history.json, and the
// samples of voicebanks in voicebanks/<id>/.
type Store struct{ Dir string }

// Settings are settings.json: the active voicebank, and what the Settings
// screen sets.
type Settings struct {
	Voicebank string `json:"voicebank,omitempty"` // the active voicebank's ID
	PlayOn    PlayOn `json:"playOn,omitempty"`
}

// PlayOn is when cards play their sounds.
type PlayOn string

const (
	OnCorrect PlayOn = ""     // when answered right, the default
	OnShow    PlayOn = "show" // when they show
	OnBoth    PlayOn = "both"
)

// Settings returns settings.json, the defaults when there is none.
func (s *Store) Settings() (Settings, error) {
	var st Settings
	err := readJSON(filepath.Join(s.Dir, "settings.json"), &st)
	return st, err
}

// SaveSettings writes settings.json.
func (s *Store) SaveSettings(st Settings) error {
	return writeJSON(filepath.Join(s.Dir, "settings.json"), st)
}

// Record is a session played to its end, in history.json.
type Record struct {
	Date     time.Time `json:"date"` // when it ended
	Deck     string    `json:"deck"` // its name, as Results shows it
	Cards    int       `json:"cards"`
	Hinted   int       `json:"hinted"` // cards whose hint was shown
	Ms       int64     `json:"ms"`
	Accuracy int       `json:"accuracy"` // percent, as Results shows it
}

// AddRecord adds r to the end of history.json.
func (s *Store) AddRecord(r Record) error {
	p := filepath.Join(s.Dir, "history.json")
	var history []Record
	if err := readJSON(p, &history); err != nil {
		return err // not to write over the history it could not read
	}
	return writeJSON(p, append(history, r))
}

// Voicebank is an imported voicebank: the samples its kana play, cut from
// its files, and voicebank.json about them.
type Voicebank struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Credit  string            `json:"credit,omitempty"`
	Samples map[string]Sample `json:"samples"` // by romaji
	// Base is the note the samples are sung at, as E4, which the user
	// tells; "" when not known.
	Base string `json:"base,omitempty"`
	// Skipped says which samples the import left out, and why.
	Skipped []string `json:"skipped,omitempty"`

	dir   string
	cache map[string]*audio.Sound
}

// Sample is where a kana's sound comes from, and its file.
type Sample struct {
	Alias    string  `json:"alias"`
	Source   string  `json:"source"`   // the wav in the voicebank
	Offset   float64 `json:"offset"`   // ms into Source
	Duration float64 `json:"duration"` // ms
	From     string  `json:"from"`     // romaji of the kana whose alias it is
	File     string  `json:"file"`     // in the voicebank's folder
}

// Own reports whether the sample is the kana's own rather than one
// standing in for it.
func (s Sample) Own(romaji string) bool { return s.From == romaji }

// Own returns how many kana have samples of their own.
func (v *Voicebank) Own() int {
	n := 0
	for romaji, s := range v.Samples {
		if s.Own(romaji) {
			n++
		}
	}
	return n
}

// Sound returns the sound of the kana, read once, or nil when the
// voicebank has none, or is nil.
func (v *Voicebank) Sound(romaji string) *audio.Sound {
	if v == nil {
		return nil
	}
	s, ok := v.Samples[romaji]
	if !ok {
		return nil
	}
	if snd, ok := v.cache[s.File]; ok {
		return snd
	}
	var snd *audio.Sound
	if b, err := os.ReadFile(filepath.Join(v.dir, s.File)); err == nil {
		snd, _ = audio.ReadWAV(b)
	}
	if v.cache == nil {
		v.cache = map[string]*audio.Sound{}
	}
	v.cache[s.File] = snd
	return snd
}

// ErrNoSamples is returned for a voicebank none of whose aliases is a kana.
var ErrNoSamples = errors.New("no oto.ini alias matches a kana")

// Import reads the voicebank at p, a folder or zip file, copies the samples
// of its kana into the store, cut to their oto.ini segments and resampled
// to audio.Rate, in place of one of the same name, which SaveSettings then
// makes the active voicebank. The folder or zip is not changed. progress,
// when not nil, gets the samples cut so far and their number.
func (s *Store) Import(p string, progress func(done, total int)) (*Voicebank, error) {
	src, err := voicebank.OpenSource(p)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	read, err := voicebank.Read(src, nil)
	if err != nil {
		return nil, err
	}
	resolved := voicebank.Resolve(read)
	h := fnv.New32a()
	h.Write([]byte(read.Name))
	v := &Voicebank{ID: fmt.Sprintf("%08x", h.Sum32()), Name: read.Name, Credit: read.Credit, Samples: map[string]Sample{}}
	final := filepath.Join(s.Dir, "voicebanks", v.ID)
	// Imported again, it keeps the base note it was told.
	var old Voicebank
	if err := readJSON(filepath.Join(final, "voicebank.json"), &old); err == nil {
		v.Base = old.Base
	}
	tmp := final + ".tmp"
	if err := os.RemoveAll(tmp); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(tmp, 0o755); err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp) // gone once renamed, or after a failure

	var own []*kana.Entry
	for i := range kana.Table {
		k := &kana.Table[i]
		if r, ok := resolved[k.Romaji]; ok && r.Own(k.Romaji) {
			own = append(own, k)
		}
	}
	wavs := map[string]*audio.Sound{} // files several samples share, as a VCV bank's
	for i, k := range own {
		r := resolved[k.Romaji]
		smp, err := cut(src, wavs, r, filepath.Join(tmp, k.Romaji+".wav"))
		if err != nil {
			v.Skipped = append(v.Skipped, fmt.Sprintf("%s (%s): %v", k.Hiragana, r.Alias, err))
		} else {
			smp.File = k.Romaji + ".wav"
			v.Samples[k.Romaji] = smp
		}
		if progress != nil {
			progress(i+1, len(own))
		}
	}
	if len(v.Samples) == 0 {
		if len(v.Skipped) > 0 {
			return nil, fmt.Errorf("no usable sample: %s", v.Skipped[0])
		}
		return nil, ErrNoSamples
	}
	for romaji, r := range resolved {
		if from, ok := v.Samples[r.From]; ok && !r.Own(romaji) {
			v.Samples[romaji] = from // the file of the kana standing in
		}
	}
	if err := writeJSON(filepath.Join(tmp, "voicebank.json"), v); err != nil {
		return nil, err
	}
	if err := os.RemoveAll(final); err != nil {
		return nil, err
	}
	if err := os.Rename(tmp, final); err != nil {
		return nil, err
	}
	v.dir = final
	return v, nil
}

// cut writes the sample r as a WAV at dst.
func cut(src *voicebank.Source, wavs map[string]*audio.Sound, r voicebank.Sample, dst string) (Sample, error) {
	snd, ok := wavs[r.File]
	if !ok {
		b, err := src.ReadFile(r.File)
		if err != nil {
			return Sample{}, err
		}
		if snd, err = audio.ReadWAV(b); err != nil {
			return Sample{}, err
		}
		wavs[r.File] = snd
	}
	d := r.Duration(snd.Ms())
	if d <= 0 {
		return Sample{}, errors.New("the sample is 0 ms long")
	}
	seg := snd.Segment(r.Offset, d).Resample(audio.Rate)
	if err := os.WriteFile(dst, audio.WriteWAV(seg), 0o644); err != nil {
		return Sample{}, err
	}
	return Sample{Alias: r.Alias, Source: r.File, Offset: r.Offset, Duration: d, From: r.From}, nil
}

// Active returns the active voicebank, nil when there is none.
func (s *Store) Active() (*Voicebank, error) {
	st, err := s.Settings()
	if err != nil || st.Voicebank == "" {
		return nil, err
	}
	dir := filepath.Join(s.Dir, "voicebanks", st.Voicebank)
	v := &Voicebank{dir: dir}
	if err := readJSON(filepath.Join(dir, "voicebank.json"), v); err != nil {
		return nil, err
	}
	return v, nil
}

// Save writes v's voicebank.json again, as after its base note changes.
func (v *Voicebank) Save() error { return writeJSON(filepath.Join(v.dir, "voicebank.json"), v) }

// readJSON reads the JSON file at p into v, leaving v as it is when there is
// no file.
func readJSON(p string, v any) error {
	b, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// writeJSON writes v to p through a temporary file, so that a crash
// leaves the old file whole.
func writeJSON(p string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}
