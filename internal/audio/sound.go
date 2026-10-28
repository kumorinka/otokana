package audio

import (
	"encoding/binary"
	"math"
)

// Segment returns the part of s lengthMs long from fromMs, within s.
func (s *Sound) Segment(fromMs, lengthMs float64) *Sound {
	at := func(ms float64) int {
		return min(max(int(math.Round(ms*float64(s.Rate)/1000)), 0), len(s.Samples))
	}
	from := at(fromMs)
	to := max(at(fromMs+lengthMs), from)
	return &Sound{Rate: s.Rate, Samples: append([]int16(nil), s.Samples[from:to]...)}
}

// Resample returns s at rate, interpolated linearly: the same sound, in
// more or fewer samples.
func (s *Sound) Resample(rate int) *Sound {
	if rate == s.Rate || len(s.Samples) == 0 {
		return &Sound{Rate: rate, Samples: append([]int16(nil), s.Samples...)}
	}
	n := int(int64(len(s.Samples)) * int64(rate) / int64(s.Rate))
	out := make([]int16, n)
	step := float64(s.Rate) / float64(rate)
	last := len(s.Samples) - 1
	for i := range out {
		pos := float64(i) * step
		j := int(pos)
		if j >= last {
			out[i] = s.Samples[last]
			continue
		}
		f := pos - float64(j)
		out[i] = int16(math.Round(float64(s.Samples[j])*(1-f) + float64(s.Samples[j+1])*f))
	}
	return &Sound{Rate: rate, Samples: out}
}

// Shift returns s the semitones higher, or lower when they are negative,
// and shorter or longer as much, as a playbackRate of 2^(semitones/12)
// does: resampled, and played at s's rate.
func (s *Sound) Shift(semitones int) *Sound {
	if semitones == 0 {
		return s
	}
	out := s.Resample(int(math.Round(float64(s.Rate) / math.Pow(2, float64(semitones)/12))))
	out.Rate = s.Rate
	return out
}

// FadeMs is how long a sound fades in and out as it plays, against clicks.
const FadeMs = 5

// faded returns s's samples faded in and out linearly over FadeMs.
func (s *Sound) faded() []int16 {
	out := append([]int16(nil), s.Samples...)
	n := min(s.Rate*FadeMs/1000, len(out)/2)
	for i := range n {
		g := float64(i) / float64(n)
		out[i] = int16(float64(out[i]) * g)
		out[len(out)-1-i] = int16(float64(out[len(out)-1-i]) * g)
	}
	return out
}

func pcm(samples []int16) []byte {
	b := make([]byte, 0, 2*len(samples))
	for _, v := range samples {
		b = binary.LittleEndian.AppendUint16(b, uint16(v))
	}
	return b
}

// WriteWAV encodes s as a WAV of 16-bit mono PCM, as ReadWAV reads.
func WriteWAV(s *Sound) []byte {
	data := pcm(s.Samples)
	le := binary.LittleEndian
	b := append([]byte("RIFF"), le.AppendUint32(nil, uint32(36+len(data)))...)
	b = append(b, "WAVEfmt "...)
	b = le.AppendUint32(b, 16)
	b = le.AppendUint16(b, formatPCM)
	b = le.AppendUint16(b, 1)
	b = le.AppendUint32(b, uint32(s.Rate))
	b = le.AppendUint32(b, uint32(2*s.Rate))
	b = le.AppendUint16(b, 2)
	b = le.AppendUint16(b, 16)
	b = append(b, "data"...)
	b = le.AppendUint32(b, uint32(len(data)))
	return append(b, data...)
}
