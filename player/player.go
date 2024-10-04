package player

import (
	"time"

	helpers "crydes/helpers"
	wrld "crydes/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	Position rl.Vector2
	Speed    float32

	CurrentAnim *Animation
	Animations  map[string]*Animation

	Map   *wrld.Map
	Sword *Sword

	DamageChan     chan bool
	IsTakingDamage bool

	LastDirection string
}

type Animation struct {
	Frames       []rl.Texture2D
	CurrentFrame int
	FrameTime    float32
	Timer        float32
	ID           string
}

func NewPlayer(x, y float32, mp *wrld.Map) *Player {
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

	p := &Player{
		Position: rl.NewVector2(x, y),
		Speed:    200.0,
		Animations: map[string]*Animation{
			"idle_right": idleRight,
			"move_right": moveRight,
			"idle_left":  idleLeft,
			"move_left":  moveLeft,
		},
		CurrentAnim:   idleRight,
		LastDirection: "right",
		Map:           mp,
		Sword: NewSword(
			rl.NewVector2(-8, -4),
			"right",
		),
		DamageChan: make(chan bool, 10),
	}

	go p.listenForDamage()

	return p
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
	moving := p.HandlePlayerMovement()

	// Set the idle animation if the player is not moving
	if !moving {
		p.SetIdleAnimation()
	}

	//! TEMOPOARYYYYYYYYYYYYYYYYY
	if rl.IsKeyPressed(rl.KeyT) {
		p.TakeDamage()
	}
	//! TEMOPOARYYYYYYYYYYYYYYYYY

	// Check for sword attack input.
	if rl.IsKeyPressed(rl.KeyF) {
		p.Sword.Visible = true
	}

	// Check for sword attack input.
	if rl.IsKeyPressed(rl.KeyF) {
		p.Sword.Visible = true
	}

	if p.Sword.Visible {
		p.Sword.Update(refreshRate, p.GetPosition(), p.LastDirection)
	}

	// Update the animation frames
	p.UpdateAnimation(refreshRate)
}

func (p *Player) HandlePlayerMovement() bool {
	// Determine target positions based on current position and speed
	targetX, targetY := p.Position.X, p.Position.Y
	moving := false

	// Horizontal movement
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		targetX += p.Speed * MOV_SPEED
		if p.IsTargetPositionWalkable(targetX, p.Position.Y) {
			p.Position.X = targetX
			p.CurrentAnim = p.Animations["move_right"]
			p.LastDirection = "right"
			moving = true
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		targetX -= p.Speed * MOV_SPEED
		if p.IsTargetPositionWalkable(targetX, p.Position.Y) {
			p.Position.X = targetX
			p.CurrentAnim = p.Animations["move_left"]
			p.LastDirection = "left"
			moving = true
		}
	}

	// Vertical movement
	if rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeyW) {
		targetY -= p.Speed * MOV_SPEED
		if p.IsTargetPositionWalkable(p.Position.X, targetY) {
			p.Position.Y = targetY
			moving = true
			p.CurrentAnim = p.Animations["move_"+p.LastDirection] // Use the last horizontal direction
		}
	}
	if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
		targetY += p.Speed * MOV_SPEED
		if p.IsTargetPositionWalkable(p.Position.X, targetY) {
			p.Position.Y = targetY
			moving = true
			p.CurrentAnim = p.Animations["move_"+p.LastDirection] // Use the last horizontal direction
		}
	}

	return moving
}

// IsTargetPositionWalkable checks if the target position is within a walkable block.
func (p *Player) IsTargetPositionWalkable(targetX, targetY float32) bool {
	// Convert the target world position to map coordinates
	mapX, mapY := int(targetX/helpers.TILE_SIZE), int(targetY/helpers.TILE_SIZE)
	// println(mapX, mapY)

	// Check if the map position is within the bounds and walkable
	return p.Map.IsWalkable(mapX, mapY)
}

// SetMovementAnimation sets the animation based on the direction.
func (p *Player) SetMovementAnimation(direction string) {
	switch direction {
	case "right":
		p.CurrentAnim = p.Animations["move_right"]
	case "left":
		p.CurrentAnim = p.Animations["move_left"]
	}
	p.LastDirection = direction
}

// SetIdleAnimation sets the idle animation based on the last direction.
func (p *Player) SetIdleAnimation() {
	if p.LastDirection == "left" {
		p.CurrentAnim = p.Animations["idle_left"]
	} else {
		p.CurrentAnim = p.Animations["idle_right"]
	}
}

// UpdateAnimation updates the current animation frame based on the refresh rate.
func (p *Player) UpdateAnimation(refreshRate float32) {
	p.CurrentAnim.Timer += refreshRate
	if p.CurrentAnim.Timer >= p.CurrentAnim.FrameTime {
		p.CurrentAnim.CurrentFrame = (p.CurrentAnim.CurrentFrame + 1) % len(p.CurrentAnim.Frames)
		p.CurrentAnim.Timer = 0
	}
}

func (p *Player) GetPosition() rl.Vector2 {
	return p.Position
}

func (p *Player) ConvertToMapPosition() (int, int) {
	return int(p.Position.X / helpers.TILE_SIZE), int(p.Position.Y / helpers.TILE_SIZE)
}

func (p *Player) Render() {
	var drawColor rl.Color
	if p.IsTakingDamage {
		drawColor = helpers.DAMAGE_COLOR
	} else {
		drawColor = rl.White // Default color
	}

	// Render the player's current animation.
	rl.DrawTextureEx(p.CurrentAnim.Frames[p.CurrentAnim.CurrentFrame], p.Position, 0, 0.5, drawColor)
	p.Sword.Render()
}

// TakeDamage method to trigger the damage effect
func (p *Player) TakeDamage() {
	p.DamageChan <- true // Send a damage event to the channel
}

// listenForDamage listens for damage events and handles them
func (p *Player) listenForDamage() {
	for {
		select {
		case <-p.DamageChan:
			p.IsTakingDamage = true
			// Wait for the duration of the damage effect
			time.Sleep(helpers.DAMAGE_DURATION)
			p.IsTakingDamage = false
		}
	}
}
