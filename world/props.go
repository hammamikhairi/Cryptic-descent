package world

import (
	"crydes/helpers"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Prop represents an interactive or static item in the game.
type Prop struct {
	ID          int                // Unique identifier for the prop
	Position    rl.Vector2         // Position of the prop
	Size        rl.Vector2         // Size of the prop for collision detection
	Scale       float32            // Scale for rendering
	Visible     bool               // Visibility of the prop
	CurrentAnim *helpers.Animation // Current animation to display if animated
	Animation   *helpers.Animation // Map of animations
	IsAnimated  bool               // Flag indicating if the prop is animated
	Color       rl.Color           // Base color for the prop (e.g., for shading effects)

	// Optional properties to control prop behavior.
	Rotation float32 // Rotation in degrees
	Opacity  float32 // Opacity for rendering (0.0 to 1.0)
	Friction float32 // Friction to apply when interacting with other objects
}

// NewProp initializes and returns a new Prop instance.
func NewProp(id int, x, y float32, scale float32, size rl.Vector2, animations *helpers.Animation, isAnimated bool) *Prop {
	return &Prop{
		ID:          id,
		Position:    rl.NewVector2(x, y),
		Size:        size,
		Scale:       scale,
		Visible:     true,
		IsAnimated:  isAnimated,
		Animation:   animations,
		CurrentAnim: animations, // Set a default animation if available
		Color:       rl.White,   // Default color
		Rotation:    0,
		Opacity:     1.0,
		Friction:    1.0,
	}
}

// Update handles animation and other dynamic properties.
func (p *Prop) Update(refreshRate float32) {
	if !p.Visible {
		return
	}

	// If animated, update animation frames

	// helpers.DEBUG("Updating prop %d", p.ID)
	// helpers.DEBUG("IsAnimated: %t", p.IsAnimated)
	// helpers.DEBUG("CurrentAnim: %v", p.CurrentAnim)

	if p.IsAnimated && p.CurrentAnim != nil {
		// helpers.DEBUG("====Updating animation for prop %d", p.ID)
		p.UpdateAnimation(refreshRate)
	}
}

// UpdateAnimation updates the current animation frame of the prop.
func (p *Prop) UpdateAnimation(refreshRate float32) {
	if p.CurrentAnim == nil {
		return
	}

	p.CurrentAnim.Timer += refreshRate
	if p.CurrentAnim.Timer >= p.CurrentAnim.FrameTime {
		p.CurrentAnim.CurrentFrame = (p.CurrentAnim.CurrentFrame + 1) % len(p.CurrentAnim.Frames)
		p.CurrentAnim.Timer = 0
	}
}

// Render draws the prop based on its properties and state.
func (p *Prop) Render() {
	// if !p.Visible {
	// 	return
	// }

	// Determine color with opacity
	finalColor := rl.Fade(p.Color, p.Opacity)

	if p.IsAnimated && p.CurrentAnim != nil {
		rl.DrawTextureEx(p.CurrentAnim.Frames[p.CurrentAnim.CurrentFrame], p.Position, p.Rotation, p.Scale, finalColor)
	} else {
		rl.DrawTextureEx(p.Animation.Frames[0], p.Position, p.Rotation, p.Scale, finalColor)
	}
}

// SetPosition updates the position of the prop.
func (p *Prop) SetPosition(x, y float32) {
	p.Position = rl.NewVector2(x, y)
}

// SetVisibility toggles the visibility of the prop.
func (p *Prop) SetVisibility(visible bool) {
	p.Visible = visible
}

// RandomizePosition places the prop at a random position within the given bounds.
func (p *Prop) RandomizePosition(bounds rl.Rectangle) {
	x := bounds.X + rand.Float32()*(bounds.Width-p.Size.X*p.Scale)
	y := bounds.Y + rand.Float32()*(bounds.Height-p.Size.Y*p.Scale)
	p.Position = rl.NewVector2(x, y)
}

// CheckCollision checks if the prop collides with another rectangle.
func (p *Prop) CheckCollision(target rl.Rectangle) bool {
	propRect := rl.NewRectangle(p.Position.X, p.Position.Y, p.Size.X*p.Scale, p.Size.Y*p.Scale)
	return rl.CheckCollisionRecs(propRect, target)
}

// ApplyFriction reduces the velocity of the prop based on its friction coefficient.
func (p *Prop) ApplyFriction(velocity rl.Vector2) rl.Vector2 {
	velocity.X *= p.Friction
	velocity.Y *= p.Friction
	return velocity
}
