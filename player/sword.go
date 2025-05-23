package player

import (
	helpers "crydes/helpers"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sword struct {
	Position  rl.Vector2
	Animation *helpers.Animation
	Visible   bool
	Direction string

	Offset rl.Vector2
	Scale  float32
}

// NewSword creates a new sword instance with the given sprite.
func NewSword(offset rl.Vector2, direction string) *Sword {
	// Load the left sprite of the sword and use it as the base frame.

	framesPaths := []string{
		"assets/sword/1.png",
		"assets/sword/2.png",
		"assets/sword/3.png",
		"assets/sword/4.png",
		"assets/sword/5.png",
	}

	frames := make([]rl.Texture2D, 0, len(framesPaths))
	for _, framePath := range framesPaths {
		frame := rl.LoadTexture(framePath)

		if frame.ID == 0 {
			rl.TraceLog(rl.LogError, "Failed to load sword sprite: %s", framePath)
			panic("Failed to load sword sprite")
		}

		frames = append(frames, frame)
	}

	// Create a simple idle sword animation with just one frame.
	swordAnimation := &helpers.Animation{
		Frames:       frames,
		FrameTime:    0.1, // Modify frame time if you have multiple frames.
		CurrentFrame: 0,
		Timer:        0,
		ID:           "sword_idle",
	}

	return &Sword{
		Position:  rl.NewVector2(0, 0),
		Animation: swordAnimation,
		Visible:   false,
		Direction: direction, // This indicates whether the sprite is mirrored for the right direction.
		Offset:    offset,
		Scale:     0.5,
	}
}

func (s *Sword) Update(refreshRate float32, playerPos rl.Vector2, playerDirection string) {

	if !s.Visible {
		return
	}

	s.Position = rl.NewVector2(playerPos.X+s.Offset.X, playerPos.Y+s.Offset.Y)
	s.Direction = playerDirection

	s.Animation.Timer += refreshRate
	if s.Animation.Timer >= s.Animation.FrameTime {
		s.Animation.CurrentFrame = (s.Animation.CurrentFrame + 1) % len(s.Animation.Frames)
		s.Animation.Timer = 0

		if s.Animation.CurrentFrame == len(s.Animation.Frames)-1 {
			s.Visible = false
			s.Animation.CurrentFrame = 0
		}
	}

}

func (s *Sword) Render() {

	// Draw a rectangle outline around the enemy after rendering the texture
	// bounds := s.GetSwordRect()
	// rl.DrawRectangleLines(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), rl.Red)

	if !s.Visible {
		return
	}

	rotation := float32(0)

	if s.Direction == "right" {
		// Mirror the texture vertically
		rotation = 180
		// Adjust the position to compensate for the rotation
		drawPosition := rl.Vector2{
			X: s.Position.X + float32(s.Animation.Frames[s.Animation.CurrentFrame].Width)*s.Scale,
			Y: s.Position.Y + float32(s.Animation.Frames[s.Animation.CurrentFrame].Height)*s.Scale,
		}
		rl.DrawTextureEx(s.Animation.Frames[s.Animation.CurrentFrame], drawPosition, rotation, s.Scale, rl.White)
	} else {
		rl.DrawTextureEx(s.Animation.Frames[s.Animation.CurrentFrame], s.Position, rotation, s.Scale, rl.White)
	}
}

func (s *Sword) GetSwordRect() rl.Rectangle {
	width := float32(s.Animation.Frames[s.Animation.CurrentFrame].Width) * 0.5 * s.Scale
	height := float32(s.Animation.Frames[s.Animation.CurrentFrame].Height) * s.Scale

	if s.Direction == "right" {
		return rl.NewRectangle(s.Position.X+width, s.Position.Y, width, height)
	}

	return rl.NewRectangle(s.Position.X, s.Position.Y, width, height)
}

func (s *Sword) ResetAttack() {
	s.Visible = true
	s.Animation.CurrentFrame = 0
}
