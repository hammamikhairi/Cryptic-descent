package enemies

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	helpers "crydes/helpers"
	"crydes/player"
)

type EnemiesManager struct {
	Enemies    []*Enemy
	Animations map[string]*map[string]*helpers.Animation
}

type Enemy struct {
	ID       int
	Position rl.Vector2
	Size     rl.Vector2
	Scale    float32
	Speed    float32
	Health   int

	CurrentAnim   *helpers.Animation
	Animations    *map[string]*helpers.Animation
	LastDirection string

	DamageChan     chan bool
	IsTakingDamage bool
}

func NewEnemiesManager(pX, pY float32) *EnemiesManager {
	manager := &EnemiesManager{
		Enemies:    []*Enemy{},
		Animations: map[string]*map[string]*helpers.Animation{},
	}

	manager.LoadAnimations()

	manager.Enemies = append(manager.Enemies, NewEnemy(1, pX, pY, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["spider"]))

	return manager
}

func (em *EnemiesManager) LoadAnimations() {
	SPIRDER_idleRight := helpers.LoadAnimation("IDLE_R",
		"assets/spider/1.png",
		"assets/spider/2.png",
	)
	SPIRDER_moveRight := helpers.LoadAnimation("MOV_R",
		"assets/spider/9.png",
		"assets/spider/10.png",
		"assets/spider/11.png",
		"assets/spider/12.png",
	)
	SPIRDER_idleLeft := helpers.LoadAnimation("IDLE_L",
		"assets/spider/5.png",
		"assets/spider/6.png",
	)
	SPIRDER_moveLeft := helpers.LoadAnimation("MOV_L",
		"assets/spider/13.png",
		"assets/spider/14.png",
		"assets/spider/15.png",
		"assets/spider/16.png",
	)

	em.Animations["spider"] = &map[string]*helpers.Animation{
		"idle_right": SPIRDER_idleRight,
		"move_right": SPIRDER_moveRight,
		"idle_left":  SPIRDER_idleLeft,
		"move_left":  SPIRDER_moveLeft,
	}
}

func (em *EnemiesManager) Update(refreshRate float32, p *player.Player) {
	for _, e := range em.Enemies {
		e.Update(refreshRate, p)
	}
}

func (em *EnemiesManager) Render() {
	for _, e := range em.Enemies {
		e.Render()
	}
}

func (e *Enemy) GetBounds() rl.Rectangle {
	return rl.NewRectangle(e.Position.X, e.Position.Y, e.Size.X*e.Scale, e.Size.Y*e.Scale)
}

func NewEnemy(
	id int,
	x, y float32,
	scale float32,
	size rl.Vector2,
	speed float32,
	animations *map[string]*helpers.Animation,
) *Enemy {

	e := &Enemy{
		Position:      rl.NewVector2(x, y),
		Speed:         speed,
		Animations:    animations,
		Scale:         scale,
		Size:          size,
		CurrentAnim:   (*animations)["idle_right"],
		LastDirection: "right",
		DamageChan:    make(chan bool, 10),
	}

	return e
}

const MOV_SPEED = 0.001 // Movement speed
const EPSILON = 0.001
const BOUNCE_BACK_DISTANCE = 6         // Distance to bounce back
const DIRECTION_CHANGE_THRESHOLD = 5.0 // Adjust this value as needed

func (e *Enemy) Update(refreshRate float32, p *player.Player) {
	// Calculate the distance between the enemy and the player
	distance := helpers.GetDistance(e.Position, p.Position)

	// Debug log for positions
	// fmt.Println("Enemy Position:", e.Position, "Player Position:", p.Position)

	// Check if the player is within range
	if distance < 100 {

		// Calculate the difference in position
		deltaX := p.Position.X - e.Position.X
		deltaY := p.Position.Y - e.Position.Y

		// Check if the enemy is outside the deadzone
		// Normalize the movement direction
		moveX := float32(0.0)
		moveY := float32(0.0)

		// Move towards the player
		if abs(deltaX) > DIRECTION_CHANGE_THRESHOLD {
			if deltaX > 0 {
				moveX = MOV_SPEED * e.Speed
				e.CurrentAnim = (*e.Animations)["move_right"]
				e.LastDirection = "right"
			} else {
				moveX = -MOV_SPEED * e.Speed
				e.CurrentAnim = (*e.Animations)["move_left"]
				e.LastDirection = "left"
			}
		} else {
			// When close to aligned on X-axis, maintain the last direction
			if e.LastDirection == "right" {
				moveX = MOV_SPEED * e.Speed
				e.CurrentAnim = (*e.Animations)["move_right"]
			} else {
				moveX = -MOV_SPEED * e.Speed
				e.CurrentAnim = (*e.Animations)["move_left"]
			}
		}

		if deltaY > 0 {
			moveY = MOV_SPEED * e.Speed
		} else if deltaY < 0 {
			moveY = -MOV_SPEED * e.Speed
		}

		// Update position based on calculated movements
		e.Position.X += moveX
		e.Position.Y += moveY

		// Check for collision with the player
		if distance < 7 { // Adjust this threshold based on the size of the player and enemy
			p.TakeDamage() // Invoke the damage method on the player

			// Bounce the enemy back slightly
			if e.LastDirection == "right" {
				e.Position.X -= BOUNCE_BACK_DISTANCE // Bounce back to the left
			} else {
				e.Position.X += BOUNCE_BACK_DISTANCE // Bounce back to the right
			}
		}
	} else {
		// Enemy is idle when out of range
		if e.LastDirection == "right" {
			e.CurrentAnim = (*e.Animations)["idle_right"]
		} else {
			e.CurrentAnim = (*e.Animations)["idle_left"]
		}
	}

	// Update the current animation
	e.UpdateAnimation(refreshRate)
}

// Helper function to get the absolute value
func abs(value float32) float32 {
	if value < 0 {
		return -value
	}
	return value
}
func (e *Enemy) UpdateAnimation(refreshRate float32) {
	e.CurrentAnim.Timer += refreshRate
	if e.CurrentAnim.Timer >= e.CurrentAnim.FrameTime {
		e.CurrentAnim.CurrentFrame = (e.CurrentAnim.CurrentFrame + 1) % len(e.CurrentAnim.Frames)
		e.CurrentAnim.Timer = 0
	}
}

func (e *Enemy) Render() {
	var drawColor rl.Color
	if e.IsTakingDamage {
		drawColor = helpers.DAMAGE_COLOR
	} else {
		drawColor = rl.White // Default color
	}

	// Render the enemy's current animation.
	rl.DrawTextureEx(e.CurrentAnim.Frames[e.CurrentAnim.CurrentFrame], e.Position, 0, e.Scale, drawColor)

	// Draw a rectangle outline around the enemy after rendering the texture
	bounds := e.GetBounds()
	rl.DrawRectangleLines(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), rl.Red)
}
