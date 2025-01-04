package screens

import (
	"crydes/audio"
	"crydes/helpers"
	"crydes/player"
	"crydes/world"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TitleScreen struct {
	buttons      []*Button
	soundManager *audio.SoundManager
	nextScreen   ScreenType
	muteButton   *Button

	// Demo scene components
	demoWorld        *world.World
	demoPlayer       *player.Player
	demoCamera       rl.Camera2D
	moveDirection    rl.Vector2
	moveTimer        float32
	pathfinder       *world.Pathfinder
	currentPath      []rl.Vector2
	pathIndex        int
	demoCollectibles []rl.Vector2
	collectibleTimer float32
	attackTimer      float32
	attackInterval   float32
}

func NewTitleScreen(soundManager *audio.SoundManager) *TitleScreen {
	ts := &TitleScreen{
		soundManager:     soundManager,
		nextScreen:       TITLE,
		demoCollectibles: make([]rl.Vector2, 0),
		attackTimer:      0,
		attackInterval:   2.0,
	}

	// Initialize demo world
	ts.demoWorld = world.NewWorld()

	// Initialize pathfinder
	ts.pathfinder = world.NewPathfinder(ts.demoWorld.Map)

	// Initialize demo player at spawn position
	x, y := ts.demoWorld.PlayerSpawn()
	ts.demoPlayer = player.NewPlayer(x, y, ts.demoWorld.Map, soundManager, nil)

	// Set up camera for demo scene
	ts.demoCamera = rl.Camera2D{
		Offset:   rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 3},
		Target:   rl.Vector2{X: x, Y: y},
		Rotation: 0.0,
		Zoom:     3.0,
	}

	// Initialize random movement direction
	ts.moveDirection = rl.Vector2{
		X: float32(rand.Float64()*2 - 1),
		Y: float32(rand.Float64()*2 - 1),
	}
	ts.moveTimer = 0

	ts.Init()
	return ts
}

func (ts *TitleScreen) Init() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	buttonWidth := float32(200)
	buttonHeight := float32(50)
	startX := (screenWidth - buttonWidth) / 2
	startY := screenHeight/2 - buttonHeight

	// Create mute button in top-right corner
	muteButtonSize := float32(40)
	ts.muteButton = NewButton(
		screenWidth-muteButtonSize-10,
		10,
		muteButtonSize,
		muteButtonSize,
		"🔊",
		func() {
			if ts.soundManager.GetCurrentMusic() != "" {
				ts.soundManager.SetVolume(audio.MUSIC, 0)
				ts.muteButton.Text = "🔇"
			} else {
				ts.soundManager.SetVolume(audio.MUSIC, audio.MUSIC_BASE)
				ts.muteButton.Text = "🔊"
			}
			ts.soundManager.RequestSound("menu_select", 1.0, 1.0)
		},
	)

	ts.buttons = []*Button{
		NewButton(startX, startY, buttonWidth, buttonHeight, "Play", func() {
			ts.soundManager.RequestSound("menu_select", 1.0, 1.0)
			ts.nextScreen = GAME
		}),
		NewButton(startX, startY+buttonHeight+20, buttonWidth, buttonHeight, "Quit", func() {
			ts.soundManager.RequestSound("menu_select", 1.0, 1.0)
			rl.CloseWindow()
		}),
	}

	// Start title screen music
	ts.soundManager.RequestMusic("title_theme", true)
}

func (ts *TitleScreen) Update(deltaTime float32) bool {
	// Update demo scene
	ts.updateDemoScene(deltaTime)

	// Update buttons including mute button
	ts.muteButton.Update()
	for _, button := range ts.buttons {
		button.Update()
	}

	if ts.nextScreen == GAME {
		ts.soundManager.RequestMusic("dungeon_theme", true)
	}

	return ts.nextScreen == GAME
}

func (ts *TitleScreen) updateDemoScene(deltaTime float32) {
	// Check if we need a new path
	if len(ts.currentPath) == 0 || ts.pathIndex >= len(ts.currentPath) {
		// Get current position in tile coordinates
		startX := int(ts.demoPlayer.Position.X) / helpers.TILE_SIZE
		startY := int(ts.demoPlayer.Position.Y) / helpers.TILE_SIZE

		// Generate random target position
		targetX, targetY := startX, startY
		for attempts := 0; attempts < 100; attempts++ {
			targetX = rand.Intn(helpers.MAP_WIDTH)
			targetY = rand.Intn(helpers.MAP_HEIGHT)
			if ts.demoWorld.Map.IsWalkable(targetX, targetY) {
				break
			}
		}

		// Find path to target
		ts.pathfinder.Update(startX, startY, targetX, targetY)
		ts.currentPath = ts.pathfinder.CreateSmoothPath()
		ts.pathIndex = 0
	}

	// Move towards next path point
	if ts.pathIndex < len(ts.currentPath) {
		target := ts.currentPath[ts.pathIndex]
		direction := rl.Vector2Subtract(target, ts.demoPlayer.Position)
		distance := rl.Vector2Length(direction)

		if distance < 2.0 {
			ts.pathIndex++
		} else {
			moveSpeed := float32(100.0)
			direction = rl.Vector2Normalize(direction)
			ts.demoPlayer.Position.X += direction.X * moveSpeed * deltaTime
			ts.demoPlayer.Position.Y += direction.Y * moveSpeed * deltaTime
			// Determine the appropriate animation
			var newAnim *helpers.Animation
			if direction.X > 0 {
				newAnim = ts.demoPlayer.Animations["move_right"]
				ts.demoPlayer.LastDirection = "right"
			} else if direction.X < 0 {
				newAnim = ts.demoPlayer.Animations["move_left"]
				ts.demoPlayer.LastDirection = "left"
			} else if direction.Y < 0 {
				newAnim = ts.demoPlayer.Animations["move_up"]
			} else if direction.Y > 0 {
				newAnim = ts.demoPlayer.Animations["move_down"]
			}

			// Only update the animation if it's different from the current one
			if ts.demoPlayer.CurrentAnim.ID != newAnim.ID {
				ts.demoPlayer.CurrentAnim = newAnim
			}
		}
	}

	// Update attack timer
	ts.attackTimer += deltaTime
	if ts.attackTimer >= ts.attackInterval {
		ts.demoPlayer.Attack()
		ts.attackTimer = 0
	}

	// Update player animation
	ts.demoPlayer.Update(deltaTime)

	// Update camera to follow player
	ts.demoCamera.Target = ts.demoPlayer.Position

}

func (ts *TitleScreen) Type() ScreenType {
	return TITLE
}

func (ts *TitleScreen) Render() {
	// Draw demo scene
	rl.BeginMode2D(ts.demoCamera)
	ts.demoWorld.Render()

	// Draw collectibles
	// for _, pos := range ts.demoCollectibles {
	// 	rl.DrawCircleV(pos, 8, rl.Yellow)
	// 	rl.DrawCircleV(pos, 4, rl.Gold)
	// }

	ts.demoPlayer.Render()
	rl.EndMode2D()

	// Draw semi-transparent overlay
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()),
		rl.ColorAlpha(rl.Black, 0.8))

	// Draw title with gradient
	titleText := "Cryptic Descent"
	fontSize := int32(60)
	textWidth := rl.MeasureText(titleText, fontSize)
	centerX := float32(rl.GetScreenWidth()-int(textWidth)) / 2
	centerY := float32(100)

	// Draw gradient title (multiple layers for bold effect)
	for offset := -2; offset <= 2; offset++ {
		for yOffset := -2; yOffset <= 2; yOffset++ {
			// Calculate gradient color
			gradientFactor := float32(offset+2) / 4.0
			color := rl.Color{
				R: uint8(192 + gradientFactor*63), // 192-255 range for silver
				G: uint8(192 + gradientFactor*63),
				B: uint8(192 + gradientFactor*63),
				A: 255,
			}

			rl.DrawText(titleText,
				int32(centerX)+int32(offset),
				int32(centerY)+int32(yOffset),
				fontSize,
				color)
		}
	}

	// Draw mute button
	ts.muteButton.Render()

	// Draw credits at the bottom
	creditsText := "Created by Khairi Hammami"
	copyrightText := "© 2025 All Rights Reserved"
	creditsFontSize := int32(20)

	creditsWidth := rl.MeasureText(creditsText, creditsFontSize)
	copyrightWidth := rl.MeasureText(copyrightText, creditsFontSize)

	screenHeight := float32(rl.GetScreenHeight())

	rl.DrawText(creditsText,
		int32(float32(rl.GetScreenWidth()-int(creditsWidth))/2),
		int32(screenHeight-80),
		creditsFontSize,
		rl.Gray)

	rl.DrawText(copyrightText,
		int32(float32(rl.GetScreenWidth()-int(copyrightWidth))/2),
		int32(screenHeight-40),
		creditsFontSize,
		rl.Gray)

	// Draw buttons
	for _, button := range ts.buttons {
		button.Render()
	}
}

func (ts *TitleScreen) Unload() {
	// Cleanup if needed
}
