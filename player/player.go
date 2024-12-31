package player

import (
	"crydes/audio"
	effects "crydes/effects/particle"
	helpers "crydes/helpers"
	wrld "crydes/world"

	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Effect struct {
	Type      string
	Value     float32
	Duration  time.Duration
	ExpiresAt time.Time
}

type Player struct {
	Position rl.Vector2
	Speed    float32
	Health   int
	Scale    float32

	CurrentAnim *helpers.Animation
	Animations  map[string]*helpers.Animation

	Map   *wrld.Map
	Sword *Sword

	DamageChan     chan bool
	IsTakingDamage bool
	AttackChan     chan rl.Rectangle

	LastDirection  string
	State          string // Add a state field to track the current state
	HeartTexture   rl.Texture2D
	heartParticles *effects.ParticleSystem
	lastHealth     int

	audio         *audio.SoundManager
	effectsChan   <-chan wrld.ItemEffectEvent
	ActiveEffects map[string]*Effect
}

func NewPlayer(x, y float32, mp *wrld.Map, sm *audio.SoundManager, effectsChan <-chan wrld.ItemEffectEvent) *Player {
	idleRight := helpers.LoadAnimation("IDLE_R",
		"assets/player/1.png",
		"assets/player/2.png",
		"assets/player/3.png",
	)
	moveRight := helpers.LoadAnimation("MOV_R",
		"assets/player/15.png",
		"assets/player/16.png",
		"assets/player/17.png",
		"assets/player/18.png",
	)
	idleLeft := helpers.LoadAnimation("IDLE_L",
		"assets/player/8.png",
		"assets/player/9.png",
		"assets/player/10.png",
	)
	moveLeft := helpers.LoadAnimation("MOV_L",
		"assets/player/22.png",
		"assets/player/23.png",
		"assets/player/24.png",
		"assets/player/25.png",
	)
	damageLeft := helpers.LoadAnimation("DAMAGE_R",
		"assets/player/29.png",
		"assets/player/30.png",
		"assets/player/31.png",
		"assets/player/32.png",
		"assets/player/33.png",
	)
	damageRight := helpers.LoadAnimation("DAMAGE_L",
		"assets/player/36.png",
		"assets/player/37.png",
		"assets/player/38.png",
		"assets/player/39.png",
		"assets/player/40.png",
	)
	die := helpers.LoadAnimation("DIE",
		"assets/player/57.png",
		"assets/player/58.png",
		"assets/player/59.png",
		"assets/player/60.png",
		"assets/player/61.png",
		"assets/player/62.png",
		"assets/player/63.png",
	)

	heartTexture := rl.LoadTexture("assets/ui/heart.png")

	p := &Player{
		Position: rl.NewVector2(x, y),
		Speed:    200.0,
		Animations: map[string]*helpers.Animation{
			"idle_right":   idleRight,
			"move_right":   moveRight,
			"idle_left":    idleLeft,
			"move_left":    moveLeft,
			"damage_right": damageLeft,
			"damage_left":  damageRight,
			"die":          die,
		},
		CurrentAnim:   idleRight,
		LastDirection: "right",
		Map:           mp,
		Sword: NewSword(
			rl.NewVector2(-8, -4),
			"right",
		),
		DamageChan:     make(chan bool, 10),
		AttackChan:     make(chan rl.Rectangle, 10),
		Health:         5,
		Scale:          0.5,
		HeartTexture:   heartTexture,
		heartParticles: effects.NewParticleSystem(),
		lastHealth:     5,
		audio:          sm,
		effectsChan:    effectsChan,
		ActiveEffects:  make(map[string]*Effect),
	}

	go p.listenForEffects()
	go p.listenForDamage()

	return p
}

const MOV_SPEED = 0.006

func (p *Player) Update(refreshRate float32) {
	// Update effects at the start of each frame
	p.updateEffects()

	if p.CheckHealth(); p.State == "dying" {
		p.CurrentAnim = p.Animations["die"]
		p.UpdateAnimation(refreshRate)
		return
	}

	switch p.State {
	case "taking_damage":
		// Let the damage animation play out; no other actions allowed.
		p.HandlePlayerMovement()

		p.CurrentAnim = p.Animations["damage_"+p.LastDirection]

		if rl.IsKeyPressed(rl.KeySpace) {
			p.Attack()
		}

		// change rotation of the player
		// if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// 	mousePos := rl.GetMousePosition()
		// 	helpers.DEBUG("Mouse Click", mousePos)
		// 	helpers.DEBUG("Player Position", p.Position)
		// 	p.HandleMouseClick(mousePos)
		// }

	case "dying":
		// Play the death animation, no other actions should be allowed.

	default:
		// Allow player to move and attack if not taking damage or dying.
		moving := p.HandlePlayerMovement()

		if !moving {
			p.SetIdleAnimation()
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			p.Attack()
		}

		// change rotation of the player
		// if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// 	mousePos := rl.GetMousePosition()
		// 	helpers.DEBUG("Mouse Click", mousePos)
		// 	helpers.DEBUG("Player Position", p.Position)
		// 	p.HandleMouseClick(mousePos)
		// }

		//! TEMPORATRRARAR

		if rl.IsKeyPressed(rl.KeyE) {
			p.Die()
		}

		//! TEMPORATRRARAR

	}

	// If the sword is visible, update its position
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
	// fmt.Println(p.Speed)

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

	if p.State == "dying" && p.CurrentAnim.CurrentFrame == len(p.CurrentAnim.Frames)-1 {
		return
	}

	if p.CurrentAnim.Timer >= p.CurrentAnim.FrameTime {
		p.CurrentAnim.CurrentFrame = (p.CurrentAnim.CurrentFrame + 1) % len(p.CurrentAnim.Frames)
		p.CurrentAnim.Timer = 0

		if p.State == "taking_damage" && p.CurrentAnim.CurrentFrame == len(p.CurrentAnim.Frames)-1 {
			p.IsTakingDamage = false
			p.State = "idle" // Reset the state to idle
		}

	}
}

func (p *Player) GetPosition() rl.Vector2 {
	return p.Position
}

func (p *Player) ConvertToMapPosition() (int, int) {
	return int(p.Position.X / helpers.TILE_SIZE), int(p.Position.Y / helpers.TILE_SIZE)
}

func (p *Player) Render() {
	// Draw the current animation frame.
	rl.DrawTextureEx(p.CurrentAnim.Frames[p.CurrentAnim.CurrentFrame], p.Position, 0, p.Scale, rl.White)

	// Render the sword if visible.
	p.Sword.Render()
}

// TakeDamage method to trigger the damage effect
func (p *Player) TakeDamage() {
	// If already taking damage or dying, ignore further damage.
	if p.State == "taking_damage" || p.State == "dying" {
		return
	}

	p.audio.RequestSound("damage", 1.0, 1.0)
	// Change the player's state to taking damage.
	p.State = "taking_damage"
	p.DamageChan <- true
}

// listenForDamage listens for damage events and handles them
func (p *Player) listenForDamage() {
	for {
		select {
		case <-p.DamageChan:
			// Set the damage animation and disable other actions.

			p.Health--
			helpers.DEBUG("Player Health", p.Health)

			p.IsTakingDamage = true

			// helpers.DEBUG("Player Taking Damage DIRECTION : ", p.LastDirection)
			p.CurrentAnim = p.Animations["damage_"+p.LastDirection]

			// Wait for the duration of the damage animation to complete.
			// time.Sleep(helpers.DAMAGE_DURATION)

			// Reset to idle state if not dead.
			// if p.State != "dying" {
			// 	p.IsTakingDamage = false
			// 	p.State = "idle" // Reset the state to idle
			// }
		}
	}
}

func (p *Player) CheckHealth() {
	if p.Health <= 0 {
		p.Die()
	}
}

func (p *Player) Die() {
	// Set the state to dying and play the death animation.
	p.State = "dying"
	p.audio.RequestSound("death", 1.0, 1.0)
}

func (p *Player) Attack() {
	p.Sword.Visible = true
	// p.Sword.ResetAttack()
	area := p.Sword.GetSwordRect()

	helpers.DEBUG("Player Attack", area)

	p.AttackChan <- area
}

func (p *Player) GameHasEnded() bool {
	return p.State == "dying" && p.CurrentAnim.CurrentFrame == len(p.CurrentAnim.Frames)-1
}

func (p *Player) HandleMouseClick(mousePos rl.Vector2) {

	mouseX, _ := mousePos.X, mousePos.Y

	if mouseX > p.Position.X {
		p.LastDirection = "right"
	} else {
		p.LastDirection = "left"
	}

	p.Attack()
}

func (p *Player) GetPlayerRoom() int {
	return p.Map.CurrentRoomIndex(p.Position)
}

func (p *Player) GetPlayerCenterPoint() rl.Vector2 {
	return rl.NewVector2(
		p.Position.X+float32(p.CurrentAnim.Frames[0].Width/2)*p.Scale,
		p.Position.Y+float32(p.CurrentAnim.Frames[0].Width/2)*p.Scale,
	)
}

func (p *Player) RenderHearts() {
	// Common calculations
	heartScale := float32(5.0)
	heartSize := float32(8) * heartScale
	padding := float32(5)
	startX := float32(20)
	startY := float32(rl.GetScreenHeight() - int(heartSize) - 20)

	// Draw blurry background - make it taller to accommodate effects
	totalWidth := (heartSize+padding)*float32(5) + padding
	effectHeight := float32(30) // Height for effect indicators
	bgRect := rl.Rectangle{
		X:      startX - padding,
		Y:      startY - padding - effectHeight,
		Width:  totalWidth,
		Height: heartSize + padding*2 + effectHeight,
	}
	rl.DrawRectangle(
		int32(bgRect.X),
		int32(bgRect.Y),
		int32(bgRect.Width),
		int32(bgRect.Height),
		rl.NewColor(0, 0, 0, 100),
	)

	// Draw border
	rl.DrawRectangleLinesEx(bgRect, 2, rl.ColorAlpha(rl.White, 0.3))

	// Render active effects
	effectY := startY - effectHeight + padding
	effectX := startX
	for effectType, effect := range p.ActiveEffects {
		remaining := effect.ExpiresAt.Sub(time.Now())
		if remaining > 0 {
			// Draw effect indicator
			effectColor := rl.White
			effectText := effectType
			switch effectType {
			case "speed":
				effectColor = rl.Green
				effectText = "Speed Boost"
			case "poison":
				effectColor = rl.Purple
				effectText = "Poisoned"
			}

			// Draw effect name
			rl.DrawText(effectText, int32(effectX), int32(effectY), 15, effectColor)

			// Draw timer bar
			barWidth := float32(100)
			barHeight := float32(5)
			progress := float32(remaining) / float32(effect.Duration)
			rl.DrawRectangle(
				int32(effectX),
				int32(effectY)+22,
				int32(barWidth*progress),
				int32(barHeight),
				effectColor,
			)

			effectX += barWidth + padding
		}
	}

	// Check if health has decreased
	if p.Health < p.lastHealth {
		// Emit particles at the position of each lost heart
		for i := p.Health; i < p.lastHealth; i++ {
			position := rl.Vector2{
				X: startX + (heartSize+padding)*float32(i) + heartSize/2,
				Y: startY + heartSize/2,
			}
			p.heartParticles.EmitParticles(position, 15, rl.Red, "hit")
		}
	}
	p.lastHealth = p.Health

	// Update and draw particles
	p.heartParticles.Update(rl.GetFrameTime())
	p.heartParticles.Draw()

	// Draw hearts
	for i := 0; i < p.Health; i++ {
		position := rl.Vector2{
			X: startX + (heartSize+padding)*float32(i),
			Y: startY,
		}
		rl.DrawTextureEx(p.HeartTexture, position, 0, heartScale, rl.White)
	}
}

func (p *Player) updateEffects() {
	now := time.Now()
	for effectType, effect := range p.ActiveEffects {
		if now.After(effect.ExpiresAt) {
			println("death")
			// Remove expired effect and restore original values
			switch effectType {
			case "speed":
				println("death 1")
				p.Speed /= effect.Value // Restore original speed
			case "poison":
				// Stop taking poison damage
			}
			println("death 2")
			delete(p.ActiveEffects, effectType)
		}
	}
}

func (p *Player) applyEffect(effectType string, value float32, duration time.Duration) {
	// Remove existing effect of the same type if it exists
	if existing, exists := p.ActiveEffects[effectType]; exists {
		// Remove the existing effect first
		println("herhehrehrkehrkher")
		switch effectType {
		case "speed":
			p.Speed /= existing.Value // Restore original speed before applying new one
		}
	}

	// Apply the new effect
	p.ActiveEffects[effectType] = &Effect{
		Type:      effectType,
		Value:     value,
		Duration:  duration,
		ExpiresAt: time.Now().Add(duration),
	}

	// Apply immediate effect
	switch effectType {
	case "speed":
		p.Speed *= value
	case "poison":
		go p.handlePoisonEffect(value, duration)
	}
}

func (p *Player) handlePoisonEffect(damage float32, duration time.Duration) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)

	for {
		select {
		case <-ticker.C:
			if time.Now().After(endTime) {
				return
			}
			if p.State != "dying" {
				p.Health -= int(damage)
				if p.Health <= 0 {
					p.Die()
					return
				}
			}
		}
	}
}

func (p *Player) listenForEffects() {
	for effect := range p.effectsChan {
		switch effect.Effect.Type {
		case "heal":
			p.audio.RequestSound("heal", 1.0, 1.0)
			p.Health = helpers.Min(p.Health+int(effect.Effect.Value), 5)
		case "speed":
			println("HELL YEAH")
			p.applyEffect("speed", effect.Effect.Value, effect.Effect.Duration)
		case "poison":
			p.applyEffect("poison", effect.Effect.Value, effect.Effect.Duration)
		case "key":
			// Implement key collection logic
		case "coin":
			// Implement coin collection logic
		}
	}
}
