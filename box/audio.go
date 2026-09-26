package box

import (
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

const sampleRate = beep.SampleRate(44100)

var music *beep.Ctrl
var speakerReady bool

// Audio holds a decoded MP3 stream ready for playback.
type Audio struct {
	stream beep.StreamSeekCloser
	format beep.Format
}

func ensureSpeaker() {
	if speakerReady {
		return
	}
	speaker.Init(sampleRate, sampleRate.N(100*time.Millisecond))
	speakerReady = true
}

// LoadSound opens and decodes an MP3 file.
func LoadSound(path string) (*Audio, error) {
	var f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	var stream, format, decodeErr = mp3.Decode(f)
	if decodeErr != nil {
		f.Close()
		return nil, decodeErr
	}
	return &Audio{stream, format}, nil
}

func (s *Audio) resample(inner beep.Streamer) beep.Streamer {
	if s.format.SampleRate == sampleRate {
		return inner
	}
	return beep.Resample(4, s.format.SampleRate, sampleRate, inner)
}

// PlaySound plays s once from the beginning, non-blocking.
func PlaySound(a *Audio) {
	ensureSpeaker()
	a.stream.Seek(0)
	speaker.Play(a.resample(a.stream))
}

// PlayMusic plays s in an infinite loop, replacing any currently playing music.
func PlayMusic(a *Audio) {
	ensureSpeaker()
	StopMusic()
	a.stream.Seek(0)
	var loop, _ = beep.Loop2(a.stream)
	var ctrl = &beep.Ctrl{Streamer: a.resample(loop)}
	speaker.Lock()
	music = ctrl
	speaker.Unlock()
	speaker.Play(ctrl)
}

// StopMusic stops the currently playing music, if any.
func StopMusic() {
	if music == nil {
		return
	}
	speaker.Lock()
	music.Paused = true
	music = nil
	speaker.Unlock()
}
