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

// Sound holds a decoded MP3 stream ready for playback.
type Sound struct {
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
func LoadSound(path string) (*Sound, error) {
	var f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	var stream, format, decodeErr = mp3.Decode(f)
	if decodeErr != nil {
		f.Close()
		return nil, decodeErr
	}
	return &Sound{stream, format}, nil
}

func (s *Sound) resample(inner beep.Streamer) beep.Streamer {
	if s.format.SampleRate == sampleRate {
		return inner
	}
	return beep.Resample(4, s.format.SampleRate, sampleRate, inner)
}

// PlaySound plays s once from the beginning, non-blocking.
func PlaySound(s *Sound) {
	ensureSpeaker()
	s.stream.Seek(0)
	speaker.Play(s.resample(s.stream))
}

// PlayMusic plays s in an infinite loop, replacing any currently playing music.
func PlayMusic(s *Sound) {
	ensureSpeaker()
	StopMusic()
	s.stream.Seek(0)
	var ctrl = &beep.Ctrl{Streamer: s.resample(beep.Loop(-1, s.stream))}
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
