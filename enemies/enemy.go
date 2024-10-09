package enemies

import (
	"crydes/helpers"
	"crydes/player"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

	DamageChan     chan rl.Rectangle
	IsTakingDamage bool
	isDead         bool

	CurrentRoom int
}

// Constructor for a new Enemy instance
func NewEnemy(
	id int,
	x, y float32,
	scale float32,
	size rl.Vector2,
	speed float32,
	animations *map[string]*helpers.Animation,
	health int,
	CurrentRoom int,
) *Enemy {
	e := &Enemy{
		Position:      rl.NewVector2(x, y),
		Speed:         speed,
		Animations:    animations,
		Scale:         scale,
		Size:          size,
		CurrentAnim:   (*animations)["idle_right"],
		LastDirection: "right",
		DamageChan:    make(chan rl.Rectangle, 10),
		Health:        health,
		CurrentRoom:   CurrentRoom,
	}

	go e.ListenForDamage()
	return e
}

// Updates the current animation of the enemy based on a refresh rate.
func (e *Enemy) UpdateAnimation(refreshRate float32) {
	e.CurrentAnim.Timer += refreshRate
	if e.CurrentAnim.Timer >= e.CurrentAnim.FrameTime {
		e.CurrentAnim.CurrentFrame = (e.CurrentAnim.CurrentFrame + 1) % len(e.CurrentAnim.Frames)
		e.CurrentAnim.Timer = 0

		if e.ShouldDie() && e.CurrentAnim.CurrentFrame == len(e.CurrentAnim.Frames)-1 {
			e.isDead = true
		}
	}
}

// Handles rendering of the enemy, considering its current state.
func (e *Enemy) Render() {
	if e.isDead {
		return
	}

	var drawColor rl.Color
	if e.IsTakingDamage {
		drawColor = helpers.DAMAGE_COLOR
	} else {
		drawColor = rl.White
	}

	// Draw the enemy's current animation frame.
	rl.DrawTextureEx(e.CurrentAnim.Frames[e.CurrentAnim.CurrentFrame], e.Position, 0, e.Scale, drawColor)
}

// Updates the enemy's state based on its interactions with the player.
func (e *Enemy) Update(refreshRate float32, p *player.Player) {
	if e.isDead {
		return
	}

	// Handle enemy death state and animation
	if e.ShouldDie() && !e.isDead {
		e.TriggerDeath()
		e.UpdateAnimation(refreshRate)
		return
	}

	// Handle movement and interaction with the player
	e.MoveTowardsPlayer(refreshRate, p)
	e.UpdateAnimation(refreshRate)
}

// Moves the enemy towards the player, adjusting its position and animation accordingly.
func (e *Enemy) MoveTowardsPlayer(refreshRate float32, p *player.Player) {
	// Calculate distance to the player and adjust position
	distance := helpers.GetDistance(e.Position, p.Position)

	if distance < helpers.ENEMIES_PLAYER_RANGE && (p.GetPlayerRoom() == e.CurrentRoom || e.CurrentRoom == -1) {
		deltaX := p.Position.X - e.Position.X
		deltaY := p.Position.Y - e.Position.Y

		if e.CurrentRoom != -1 {
			e.CurrentRoom = -1
		}

		moveX, moveY := e.CalculateMovement(deltaX, deltaY)

		// Update position
		e.Position.X += moveX
		e.Position.Y += moveY

		// Check for collision with the player and bounce back if necessary
		if distance < 7 {
			p.TakeDamage()
			e.BounceBack(p.Position.X, p.Position.Y)
		}
	} else {
		// Enemy is idle when out of range
		e.SetIdleAnimation()
	}
}

// Calculates movement towards the player.
func (e *Enemy) CalculateMovement(deltaX, deltaY float32) (float32, float32) {
	var moveX, moveY float32

	if helpers.ABS(deltaX) > helpers.ENEMIES_DIRECTION_CHANGE_THRESHOLD {
		if deltaX > 0 {
			moveX = helpers.ENEMIES_MOV_SPEED * e.Speed
			e.CurrentAnim = (*e.Animations)["move_right"]
			e.LastDirection = "right"
		} else {
			moveX = -helpers.ENEMIES_MOV_SPEED * e.Speed
			e.CurrentAnim = (*e.Animations)["move_left"]
			e.LastDirection = "left"
		}
	} else {
		if e.LastDirection == "right" {
			moveX = helpers.ENEMIES_MOV_SPEED * e.Speed
			e.CurrentAnim = (*e.Animations)["move_right"]
		} else {
			moveX = -helpers.ENEMIES_MOV_SPEED * e.Speed
			e.CurrentAnim = (*e.Animations)["move_left"]
		}
	}

	if deltaY > 0 {
		moveY = helpers.ENEMIES_MOV_SPEED * e.Speed
	} else if deltaY < 0 {
		moveY = -helpers.ENEMIES_MOV_SPEED * e.Speed
	}

	return moveX, moveY
}

// Sets the enemy to idle animation based on its last direction.
func (e *Enemy) SetIdleAnimation() {
	if e.LastDirection == "right" {
		e.CurrentAnim = (*e.Animations)["idle_right"]
	} else {
		e.CurrentAnim = (*e.Animations)["idle_left"]
	}
}

// Bounces the enemy back slightly after colliding with the player.
func (e *Enemy) BounceBack(x, y float32) {
	// get the angle between the enemy and the player
	angle := math.Atan2(float64(e.Position.Y-y), float64(e.Position.X-x))

	// calculate the new position
	newX := e.Position.X + float32(math.Cos(angle))*helpers.ENEMIES_BOUNCE_BACK_DISTANCE
	newY := e.Position.Y + float32(math.Sin(angle))*helpers.ENEMIES_BOUNCE_BACK_DISTANCE

	e.Position = rl.NewVector2(newX, newY)

}

// Triggers the death animation for the enemy.
func (e *Enemy) TriggerDeath() {
	if e.LastDirection != "right" {
		e.CurrentAnim = (*e.Animations)["death_right"]
	} else {
		e.CurrentAnim = (*e.Animations)["death_left"]
	}
}

// Returns whether the enemy should enter the death state.
func (e *Enemy) ShouldDie() bool {
	return e.Health <= 0
}

// Returns whether the enemy is considered dead.
func (e *Enemy) IsDead() bool {
	return e.isDead
}

// Retrieves the bounding box of the enemy.
func (e *Enemy) GetBounds() rl.Rectangle {
	return rl.NewRectangle(e.Position.X, e.Position.Y, e.Size.X*e.Scale, e.Size.Y*e.Scale)
}

// Handles damage taken by the enemy.
func (e *Enemy) TakeDamage(area rl.Rectangle) {
	if e.isDead || !helpers.CheckCollisionRecs(area, e.GetBounds()) {
		return
	}

	e.Health--
	e.IsTakingDamage = true

	// center of area
	centerX := area.X + area.Width/2
	centerY := area.Y + area.Height/2

	e.BounceBack(centerX, centerY)

	// Trigger death logic if health falls below zero
	if e.ShouldDie() {
		e.IsTakingDamage = false
		e.TriggerDeath()
		return
	}

	// Reset the damage state after a short duration
	go func() {
		<-time.After(helpers.DAMAGE_DURATION)
		e.IsTakingDamage = false
	}()
}

// Listens for damage signals on the DamageChan channel.
func (e *Enemy) ListenForDamage() {
	for area := range e.DamageChan {
		e.TakeDamage(area)
	}
}

// Destroys the enemy, releasing resources and closing channels.
func (e *Enemy) Destroy() {
	close(e.DamageChan)
}
