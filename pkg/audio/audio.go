package audio

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AudioDevice struct{}

func (ad *AudioDevice) Close() {
	rl.CloseAudioDevice()
}

func NewAudioDevice() *AudioDevice {
	rl.InitAudioDevice()
	return &AudioDevice{}
}

// -- Sound --

type Sound = rl.Sound

func LoadSound(fileName string) Sound {
	return rl.LoadSound(fileName)
}

func UnloadSound(sound Sound) {
	rl.UnloadSound(sound)
}

func (ad *AudioDevice) PlaySound(sound Sound) {
	rl.PlaySound(sound)
}

func (ad *AudioDevice) StopSound(sound Sound) {
	rl.StopSound(sound)
}

func (ad *AudioDevice) PauseSound(sound Sound) {
}

func (ad *AudioDevice) ResumeSound(sound Sound) {
	rl.ResumeSound(sound)
}
