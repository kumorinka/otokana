package audio

import (
	"bytes"
	"log"
	"slices"
	"sync"

	"github.com/ebitengine/oto/v3"
)

// Rate is the sample rate sounds play at, which imports resample to.
const Rate = 44100

// Out plays sounds on the audio device, or nothing while it is opening and
// when there is none.
type Out struct {
	mu      sync.Mutex
	ctx     *oto.Context
	players []*oto.Player // kept until they finish
}

// Open opens the audio device in the background. The app runs on without
// sound when it cannot.
func Open() *Out {
	o := &Out{}
	go func() {
		ctx, ready, err := oto.NewContext(&oto.NewContextOptions{SampleRate: Rate, ChannelCount: 1, Format: oto.FormatSignedInt16LE})
		if err == nil {
			<-ready
			err = ctx.Err()
		}
		if err != nil {
			log.Printf("audio: no device, playing without sound: %v", err)
			return
		}
		o.mu.Lock()
		o.ctx = ctx
		o.mu.Unlock()
	}()
	return o
}

// joined returns the samples of the sounds at Rate, faded, one after
// another, GapMs apart.
func joined(sounds []*Sound) []int16 {
	var samples []int16
	for _, s := range sounds {
		if s == nil || len(s.Samples) == 0 {
			continue
		}
		if s.Rate != Rate {
			s = s.Resample(Rate)
		}
		if len(samples) > 0 {
			samples = append(samples, make([]int16, Rate*GapMs/1000)...)
		}
		samples = append(samples, s.faded()...)
	}
	return samples
}

// GapMs is the silence between sounds played one after another, which
// keeps two kana from sounding as one word.
const GapMs = 150

// Play plays the sounds one after another, each faded in and out, over
// what plays already. Nil and empty sounds play nothing.
func (o *Out) Play(sounds ...*Sound) {
	samples := joined(sounds)
	if o == nil || len(samples) == 0 {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.ctx == nil {
		return
	}
	o.players = slices.DeleteFunc(o.players, func(p *oto.Player) bool { return !p.IsPlaying() })
	p := o.ctx.NewPlayer(bytes.NewReader(pcm(samples)))
	p.Play()
	o.players = append(o.players, p)
}
