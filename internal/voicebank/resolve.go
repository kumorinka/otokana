package voicebank

import "otokana/internal/kana"

// Sample is the oto entry a kana plays.
type Sample struct {
	Entry
	// From is the romaji of the kana whose alias it is: the kana's own, or
	// for ぢ, づ, ぢゃ, ぢゅ and ぢょ, that of じ, ず, じゃ, じゅ and じょ,
	// which sound the same.
	From string
}

// Own reports whether the sample is the kana's own rather than one standing
// in for it.
func (s Sample) Own(romaji string) bool { return s.From == romaji }

// standIns are the kana another's sample stands in for, by romaji.
var standIns = map[string]string{"di": "ji", "du": "zu", "dya": "ja", "dyu": "ju", "dyo": "jo"}

// Resolve returns the samples of the table's kana, by romaji: the first
// entry aliased as the kana's hiragana, else as its romaji, else as - and
// its hiragana (the first of a VCV bank's), else what stands in for it.
// Kana without any are left out; multi-pitch banks give their plain
// aliases alone.
func Resolve(vb *Voicebank) map[string]Sample {
	byAlias := map[string]Entry{}
	for _, e := range vb.Entries {
		if _, ok := byAlias[e.Alias]; !ok {
			byAlias[e.Alias] = e
		}
	}
	own := func(k *kana.Entry) (Entry, bool) {
		for _, alias := range []string{k.Hiragana, k.Romaji, "- " + k.Hiragana} {
			if e, ok := byAlias[alias]; ok {
				return e, true
			}
		}
		return Entry{}, false
	}
	out := map[string]Sample{}
	for i := range kana.Table {
		k := &kana.Table[i]
		if e, ok := own(k); ok {
			out[k.Romaji] = Sample{Entry: e, From: k.Romaji}
		}
	}
	for romaji, from := range standIns {
		if _, ok := out[romaji]; !ok {
			if s, ok := out[from]; ok {
				out[romaji] = s
			}
		}
	}
	return out
}
