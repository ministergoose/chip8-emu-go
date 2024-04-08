package hardware

import "github.com/gen2brain/beeep"

type Sound interface {
	Beep()
}

type SoundStd struct {
}

func NewSoundStd() *SoundStd {
	return &SoundStd{}
}

func (s *SoundStd) Beep() {
	go func() {
		err := beeep.Beep(beeep.DefaultFreq, 16)
		if err != nil {
			panic(err)
		}
	}()
}
