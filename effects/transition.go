package effects

import rl "github.com/gen2brain/raylib-go/raylib"

// Transition represents a fog-like screen transition
type Transition struct {
	active    bool    // Indicates if a transition is in progress
	opacity   float32 // Current opacity of the fog effect
	fadeSpeed float32 // Speed of the fade effect
	fadingOut bool    // Indicates if the transition is fading out or in
}

// NewTransition creates a new transition effect
func NewTransition(fadeSpeed float32) *Transition {
	return &Transition{
		active:    false,
		opacity:   0,
		fadeSpeed: fadeSpeed,
		fadingOut: false,
	}
}

// Start initiates the transition effect
func (t *Transition) Start() {
	t.active = true
	t.fadingOut = false
	t.opacity = 0 // Reset opacity
}

// Update handles the transition logic (fading in and out)
func (t *Transition) Update() {
	if !t.active {
		return
	}

	// Increase or decrease opacity based on the fading direction
	if t.fadingOut {
		t.opacity -= t.fadeSpeed
	} else {
		t.opacity += t.fadeSpeed
	}

	// Check if the transition has completed
	if t.opacity >= 1.0 {
		t.fadingOut = true // Start fading out
	} else if t.opacity <= 0 {
		t.active = false // Transition complete
		t.fadingOut = false
	}
}

// Render draws the fog effect if active
func (t *Transition) Render() {
	if !t.active {
		return
	}

	// Draw a semi-transparent rectangle over the entire screen
	fogColor := rl.Fade(rl.Gray, t.opacity) // Use opacity to set transparency
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), fogColor)
}
