package audio

import (
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// assets/
//   ├── audio/
//   │   ├── sfx/
//   │   │   ├── sword_swing.wav
//   │   │   ├── sword_hit.wav
//   │   │   └── ...
//   │   └── music/
//   │       ├── title_theme.mp3
//   │       ├── dungeon_theme.mp3
//   │       └── ...

const (
	MASTER_VOL = 0.1
	MUSIC_BASE = 0.3
	VFX_BASE   = 0.7
)

// SoundType represents different categories of sounds
type SoundType string

const (
	SFX   SoundType = "sfx"
	MUSIC SoundType = "music"
)

// SoundRequest represents a request to play a sound
type SoundRequest struct {
	Name      string
	Type      SoundType
	Volume    float32
	Pitch     float32
	Loop      bool
	StopMusic bool // Used for music transitions
}

// Add this struct to store sound settings
type SoundSettings struct {
	baseVolume float32
	sound      rl.Sound
}

// SoundManager manages all game audio
type SoundManager struct {
	sounds     map[string]SoundSettings
	music      map[string]rl.Music
	soundChan  chan SoundRequest
	volumes    map[SoundType]float32
	masterVol  float32
	mutex      sync.RWMutex
	isRunning  bool
	currentBGM string
}

// NewSoundManager creates a new sound manager instance
func NewSoundManager() *SoundManager {
	rl.InitAudioDevice()

	sm := &SoundManager{
		sounds:    make(map[string]SoundSettings),
		music:     make(map[string]rl.Music),
		soundChan: make(chan SoundRequest, 100), // Buffer for multiple sound requests
		volumes: map[SoundType]float32{
			SFX:   VFX_BASE,
			MUSIC: MUSIC_BASE,
		},
		masterVol: MASTER_VOL,
		isRunning: true,
	}

	sm.loadSounds()
	go sm.processSoundRequests()
	return sm
}

func (sm *SoundManager) loadSounds() {
	// Combat sounds
	sm.LoadSound("sword_swing", "assets/audio/sfx/sword_swing.mp3", 0.2)
	sm.LoadSound("sword_hit", "assets/audio/sfx/sword_hit.mp3", 0.3)
	sm.LoadSound("damage", "assets/audio/sfx/damage.mp3", 0.5)
	sm.LoadSound("death", "assets/audio/sfx/death.mp3", 0.7)
	// sm.LoadSound("enemy_death", "assets/audio/sfx/enemy_death.mp3")

	// Music tracks
	sm.LoadMusic("title_theme", "assets/audio/music/bg.mp3")
	sm.LoadMusic("dungeon_theme", "assets/audio/music/gameplay.mp3")
	// sm.LoadMusic("boss_theme", "assets/audio/music/bg.mp3")
}

// LoadSound loads a single sound effect
func (sm *SoundManager) LoadSound(name, path string, baseVolume float32) {
	sound := rl.LoadSound(path)

	if sound.Stream.Buffer == nil {
		panic("[ERROR] cant load sound " + name + " at : " + path)
	}

	sm.mutex.Lock()
	sm.sounds[name] = SoundSettings{
		baseVolume: baseVolume,
		sound:      sound,
	}
	sm.mutex.Unlock()
}

// LoadMusic loads a music track
func (sm *SoundManager) LoadMusic(name, path string) {
	music := rl.LoadMusicStream(path)

	if music.Stream.Buffer == nil {
		panic("[ERROR] cant load music " + name + " at : " + path)
	}

	sm.mutex.Lock()
	sm.music[name] = music
	sm.mutex.Unlock()
}

// processSoundRequests handles incoming sound requests in a separate goroutine
func (sm *SoundManager) processSoundRequests() {
	for sm.isRunning {
		select {
		case req := <-sm.soundChan:
			// fmt.Printf("%+v\n", sm.sounds)
			// fmt.Printf("%+v\n", sm.music)
			sm.mutex.RLock()
			switch req.Type {
			case SFX:
				if soundSettings, exists := sm.sounds[req.Name]; exists {
					// Multiply by the sound's base volume
					finalVolume := req.Volume * soundSettings.baseVolume * sm.volumes[SFX] * sm.masterVol
					rl.SetSoundVolume(soundSettings.sound, finalVolume)
					rl.SetSoundPitch(soundSettings.sound, req.Pitch)
					rl.PlaySound(soundSettings.sound)
				}
			case MUSIC:
				if req.StopMusic {
					// Stop current music if playing
					if sm.currentBGM != "" {
						if curr, exists := sm.music[sm.currentBGM]; exists {
							rl.StopMusicStream(curr)
						}
					}
				}
				if music, exists := sm.music[req.Name]; exists {
					rl.SetMusicVolume(music, req.Volume*sm.volumes[MUSIC]*sm.masterVol)
					rl.PlayMusicStream(music)
					sm.currentBGM = req.Name
				}
			}
			sm.mutex.RUnlock()
		default:
			// Update current music stream
			if sm.currentBGM != "" {
				sm.mutex.RLock()
				if music, exists := sm.music[sm.currentBGM]; exists {
					rl.UpdateMusicStream(music)
				}
				sm.mutex.RUnlock()
			}
			time.Sleep(time.Millisecond) // Prevent CPU spinning
		}
	}
}

// RequestSound sends a sound request through the channel
func (sm *SoundManager) RequestSound(name string, volume, pitch float32) {
	println("requested sound")
	sm.soundChan <- SoundRequest{
		Name:   name,
		Type:   SFX,
		Volume: volume,
		Pitch:  pitch,
	}
}

// RequestMusic sends a music request through the channel
func (sm *SoundManager) RequestMusic(name string, stopCurrent bool) {
	sm.soundChan <- SoundRequest{
		Name:      name,
		Type:      MUSIC,
		Volume:    1.0,
		StopMusic: stopCurrent,
		Loop:      true,
	}
}

// SetVolume sets the volume for a specific sound type
func (sm *SoundManager) SetVolume(sType SoundType, volume float32) {
	sm.mutex.Lock()
	sm.volumes[sType] = volume
	sm.mutex.Unlock()

	// Update current music volume if needed
	if sType == MUSIC && sm.currentBGM != "" {
		if music, exists := sm.music[sm.currentBGM]; exists {
			rl.SetMusicVolume(music, volume*sm.masterVol)
		}
	}
}

// SetMasterVolume sets the master volume
func (sm *SoundManager) SetMasterVolume(volume float32) {
	sm.mutex.Lock()
	sm.masterVol = volume
	sm.mutex.Unlock()

	// Update current music volume
	if sm.currentBGM != "" {
		if music, exists := sm.music[sm.currentBGM]; exists {
			rl.SetMusicVolume(music, sm.volumes[MUSIC]*volume)
		}
	}
}

// Unload cleans up all audio resources
func (sm *SoundManager) Unload() {
	sm.isRunning = false
	time.Sleep(time.Millisecond * 100) // Give time for goroutine to finish

	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Stop and unload current music
	if sm.currentBGM != "" {
		if music, exists := sm.music[sm.currentBGM]; exists {
			rl.StopMusicStream(music)
		}
	}

	// Unload all music
	for _, music := range sm.music {
		rl.UnloadMusicStream(music)
	}

	// Unload all sounds
	for _, soundSettings := range sm.sounds {
		rl.UnloadSound(soundSettings.sound)
	}

	rl.CloseAudioDevice()
}
