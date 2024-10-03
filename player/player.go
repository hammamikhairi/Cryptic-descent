package player

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Position    rl.Vector2
	Speed       float32
	CurrentAnim *Animation
	Animations  map[string]*Animation

	LastDirection string
}

type Animation struct {
	Frames       []rl.Texture2D
	CurrentFrame int
	FrameTime    float32
	Timer        float32
	ID           string
}

func NewPlayer(x, y float32) *Player {
	idleRight := loadAnimation("IDLE_R",
		"assets/player/1.png",
		"assets/player/2.png",
		"assets/player/3.png",
	)
	moveRight := loadAnimation("MOV_R",
		"assets/player/15.png",
		"assets/player/16.png",
		"assets/player/17.png",
		"assets/player/18.png",
	)
	idleLeft := loadAnimation("IDLE_L",
		"assets/player/8.png",
		"assets/player/9.png",
		"assets/player/10.png",
	)
	moveLeft := loadAnimation("MOV_L",
		"assets/player/22.png",
		"assets/player/23.png",
		"assets/player/24.png",
		"assets/player/25.png",
	)

	return &Player{
		Position: rl.NewVector2(x, y),
		Speed:    200.0,
		Animations: map[string]*Animation{
			"idle_right": idleRight,
			"move_right": moveRight,
			"idle_left":  idleLeft,
			"move_left":  moveLeft,
		},
		CurrentAnim: idleRight,
	}
}

func loadAnimation(id string, filePaths ...string) *Animation {

	var textures []rl.Texture2D = []rl.Texture2D{}
	for _, path := range filePaths {
		texture := rl.LoadTexture(path)

		if texture.ID == 0 {
			rl.TraceLog(rl.LogError, "Failed to load texture: %s", path)
			panic("Failed to load texture")
		}

		textures = append(textures, texture)
	}

	return &Animation{
		Frames:       textures, // You can load multiple frames here if needed
		FrameTime:    0.1,      // Adjust frame time for animation speed
		Timer:        0,
		CurrentFrame: 0,
		ID:           id,
	}
}

const MOV_SPEED = 0.006

func (p *Player) Update(refreshRate float32) {
	// Movement input handling
	moving := false

	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		p.Position.X += p.Speed * MOV_SPEED
		p.CurrentAnim = p.Animations["move_right"]
		p.LastDirection = "right"
		moving = true
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		p.Position.X -= p.Speed * MOV_SPEED
		p.CurrentAnim = p.Animations["move_left"]
		p.LastDirection = "left"
		moving = true
	}
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		p.Position.Y -= p.Speed * MOV_SPEED
		moving = true
		// Set animation based on last horizontal direction
		if p.LastDirection == "left" {
			p.CurrentAnim = p.Animations["move_left"]
		} else {
			p.CurrentAnim = p.Animations["move_right"]
		}
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		p.Position.Y += p.Speed * MOV_SPEED
		moving = true
		// Set animation based on last horizontal direction
		if p.LastDirection == "left" {
			p.CurrentAnim = p.Animations["move_left"]
		} else {
			p.CurrentAnim = p.Animations["move_right"]
		}
	}

	if !moving {
		// If no keys are pressed, set the idle animation
		if p.LastDirection == "left" {
			p.CurrentAnim = p.Animations["idle_left"]
		} else {
			p.CurrentAnim = p.Animations["idle_right"]
		}
	}

	// Update animation
	p.CurrentAnim.Timer += refreshRate
	if p.CurrentAnim.Timer >= p.CurrentAnim.FrameTime {
		p.CurrentAnim.CurrentFrame = (p.CurrentAnim.CurrentFrame + 1) % len(p.CurrentAnim.Frames)
		p.CurrentAnim.Timer = 0
	}
}

func (p *Player) Render() {

	rl.DrawTextureEx(p.CurrentAnim.Frames[p.CurrentAnim.CurrentFrame], p.Position, 0, 0.5, rl.White)
}
