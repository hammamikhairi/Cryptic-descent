package audio

import rl "github.com/gen2brain/raylib-go/raylib"

// SoundManager manages game sounds and effects
type SoundManager struct {
    BiwaSound rl.Sound // Placeholder for the biwa sound effect
}

// NewSoundManager initializes a new sound manager and loads sounds
func NewSoundManager() *SoundManager {
    rl.InitAudioDevice() // Initialize the audio device

    // Load the biwa sound from assets (ensure the path is correct)
    biwa := rl.LoadSound("assets/biwa.ogg")

    return &SoundManager{BiwaSound: biwa}
}

// PlayBiwaSound plays the biwa sound effect
func (s *SoundManager) PlayBiwaSound() {
    rl.PlaySound(s.BiwaSound)
}

// Unload sounds and close the audio device
func (s *SoundManager) Unload() {
    rl.UnloadSound(s.BiwaSound)
    rl.CloseAudioDevice()
}
