// Package audio: membaca WAV, memotong segmen, fade, resample, dan memutar lewat oto.
//
// Go murni, tanpa import mygo.
package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Sound is mono 16-bit PCM.
type Sound struct {
	Rate    int // samples a second
	Samples []int16
}

// Ms returns the sound's length in milliseconds.
func (s *Sound) Ms() float64 {
	if s.Rate <= 0 {
		return 0
	}
	return float64(len(s.Samples)) * 1000 / float64(s.Rate)
}

// The WAV formats of a fmt chunk.
const (
	formatPCM        = 1
	formatFloat      = 3
	formatExtensible = 0xFFFE
)

var errNotWAV = errors.New("not a WAV file")

// ReadWAV reads a RIFF WAV of 16-bit integer PCM, mono or stereo, which it
// mixes to mono, at any sample rate. Other formats are refused with an
// error that says which they are. A data chunk cut short is read as far
// as it goes.
func ReadWAV(b []byte) (*Sound, error) {
	if len(b) < 12 || string(b[0:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		return nil, errNotWAV
	}
	var (
		channels, bits int
		rate           int
		haveFmt        bool
	)
	for p := 12; p+8 <= len(b); {
		id, size := string(b[p:p+4]), int(binary.LittleEndian.Uint32(b[p+4:p+8]))
		body := b[p+8 : min(p+8+size, len(b))]
		switch id {
		case "fmt ":
			if len(body) < 16 {
				return nil, errNotWAV
			}
			format := binary.LittleEndian.Uint16(body[0:2])
			channels = int(binary.LittleEndian.Uint16(body[2:4]))
			rate = int(binary.LittleEndian.Uint32(body[4:8]))
			bits = int(binary.LittleEndian.Uint16(body[14:16]))
			if format == formatExtensible && len(body) >= 26 {
				format = binary.LittleEndian.Uint16(body[24:26]) // the sub-format GUID's first bytes
			}
			switch {
			case format == formatFloat:
				return nil, fmt.Errorf("%d-bit float WAV is not supported, only 16-bit PCM", bits)
			case format != formatPCM:
				return nil, fmt.Errorf("WAV format %#x is not supported, only 16-bit PCM", format)
			case bits != 16:
				return nil, fmt.Errorf("%d-bit WAV is not supported, only 16-bit PCM", bits)
			case channels != 1 && channels != 2:
				return nil, fmt.Errorf("%d-channel WAV is not supported, only mono or stereo", channels)
			case rate <= 0:
				return nil, errNotWAV
			}
			haveFmt = true
		case "data":
			if !haveFmt {
				return nil, errNotWAV
			}
			frames := len(body) / (2 * channels)
			s := &Sound{Rate: rate, Samples: make([]int16, frames)}
			for i := range frames {
				l := int16(binary.LittleEndian.Uint16(body[2*channels*i:]))
				if channels == 2 {
					r := int16(binary.LittleEndian.Uint16(body[2*channels*i+2:]))
					l = int16((int32(l) + int32(r)) / 2)
				}
				s.Samples[i] = l
			}
			return s, nil
		}
		p += 8 + size + size%2 // chunks are padded to even sizes
	}
	return nil, fmt.Errorf("%w: no audio data", errNotWAV)
}
