package core

import (
	"crydes/audio"
	"crydes/core/screens"
	"crydes/effects"
	"crydes/enemies"
	"crydes/helpers"
	"crydes/player"
	"crydes/world"
	"time"

	"fmt"

	"crydes/core/minimap"

	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ! FOR DEVELOPMENT
const (
	RENDER_LIGHTING = 1 << iota
)

//! END DEVELOPMENT

type Game struct {
	player    *player.Player
	world     *world.World
	lightning *effects.RetroLightingEffect

	soundManager *audio.SoundManager // Reference to the sound manager

	enemies *enemies.EnemiesManager

	camera        rl.Camera2D
	width, height int

	// ! FOR DEVELOPMENT
	flags int
	//! END DEVELOPMENT

	pauseScreen *screens.PauseScreen
	titleScreen *screens.TitleScreen
	outroScreen *screens.OutroScreen
	isPaused    bool
	showTitle   bool
	showOutro   bool

	minimap             *minimap.Minimap
	collectiblesManager *world.CollectibleManager

	// Dungeon shifting
	shiftTimer       float32
	isShifting       bool
	shiftDelay       float32
	shiftText        string
	shiftTextTimer   float32
	fadeAlpha        float32
	shiftSoundPlayed bool
}

// NewGame initializes a new game instance
func NewGame(soundManager *audio.SoundManager, width, height int) *Game {

	w := world.NewWorld()
	collectibleManager := world.NewCollectibleManager()

	x, y := w.PlayerSpawn()
	p := player.NewPlayer(x, y, w.Map, soundManager, collectibleManager.GetEffectsChan())
	collectibleManager.SetPlayerPosition(&p.Position)

	em := enemies.NewEnemiesManager(x, y, w.Map, p.AttackChan, w.Map.GetRoomsRects(), soundManager)
	em.SpawnEnemies()

	rle := effects.NewRetroLightingEffect(
		int32(helpers.MAP_WIDTH*helpers.TILE_SIZE), int32(helpers.MAP_HEIGHT*helpers.TILE_SIZE), 50, 2, p,
	)

	collectibleManager.ScatterCollectibles(w.Map.GetRoomsRects(), w.Map)
	collectibleManager.AddItem(1, world.HealthPotion, x+20, y+20)
	collectibleManager.AddItem(2, world.SpeedPotion, x+30, y+30)
	collectibleManager.AddItem(3, world.Poison, x-30, y+30)

	rle.SetUpPropsLightning(w.PropsManager.GetProps())

	mm := minimap.NewMinimap(w.Map)

	return &Game{
		player:              p,
		world:               w,
		soundManager:        soundManager,
		enemies:             em,
		width:               width,
		height:              height,
		lightning:           rle,
		flags:               0,
		pauseScreen:         screens.NewPauseScreen(soundManager),
		titleScreen:         screens.NewTitleScreen(soundManager),
		outroScreen:         screens.NewOutroScreen(soundManager),
		isPaused:            false,
		showTitle:           true,
		showOutro:           false,
		minimap:             mm,
		collectiblesManager: collectibleManager,
		shiftTimer:          0,
		isShifting:          false,
		// shiftDelay:          float32(4 + rand.Intn(41)), // Random value between 40 and 80 seconds
		shiftDelay:     float32(4 + rand.Intn(3)), // Random value between 40 and 80 seconds
		shiftText:      "The dungeon shifts beneath your feet...",
		shiftTextTimer: 0,
		fadeAlpha:      0,
	}

}

func (g *Game) Run() {
	rl.SetTargetFPS(60)
	previousTime := rl.GetTime()

	// Initialize the camera with correct offset for centering
	g.camera = rl.Camera2D{
		Offset:   rl.Vector2{X: float32(g.width) / 2, Y: float32(g.height) / 2},
		Target:   rl.Vector2{X: float32(g.player.Position.X), Y: float32(g.player.Position.Y)},
		Rotation: 0.0,
		Zoom:     4.5,
	}

	for !rl.WindowShouldClose() {
		deltaTime := float32(rl.GetTime() - previousTime)
		previousTime = rl.GetTime()

		// Update logic
		if !g.isPaused && !g.showTitle && !g.showOutro {
			g.Update(deltaTime)
			g.checkGameEnd() // Check for game end conditions
		} else if g.showTitle {
			if g.titleScreen.Update(deltaTime) {
				g.showTitle = false
			}
		} else if g.isPaused {
			if g.pauseScreen.Update(deltaTime) {
				g.isPaused = false
			}
		} else if g.showOutro {
			if g.outroScreen.Update(deltaTime) {
				g.showOutro = false
				g.showTitle = true // Return to title screen
			}
		}

		helpers.LogOnce(1, "HELLOOOO")

		// Update camera target to follow the player
		g.camera.Target = rl.Vector2{X: float32(g.player.Position.X), Y: float32(g.player.Position.Y)}

		// Ensure the camera offset stays centered even if window size changes
		g.camera.Offset = rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2}

		rl.BeginDrawing()
		rl.ClearBackground(helpers.VOID_COLOR) // rgb(88, 68, 34)

		// Render game world if not in title screen
		if !g.showTitle {
			rl.BeginMode2D(g.camera)
			g.Render()
			rl.EndMode2D()
		}

		// Render overlay screens
		if g.showTitle {
			g.titleScreen.Render()
		} else if g.isPaused {
			g.pauseScreen.Render()
		} else if g.showOutro {
			g.outroScreen.Render()

			rl.EndDrawing()
			continue
		}

		rl.DrawText("Cryptic Descent", 10, 10, 20, rl.Gray)
		fpsText := fmt.Sprintf("FPS: %d", rl.GetFPS())
		rl.DrawText(fpsText, 10, 35, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("CAM ZOOM : %.3f", g.camera.Zoom), 10, 60, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("DECAY FACTOR : %.3f", helpers.DECAY_FACTOR), 10, 85, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("LIGHT RADIUS : %.3f", helpers.LIGHT_RADIUS), 10, 110, 20, rl.Gray)

		// Debug information
		rl.DrawText(fmt.Sprintf("Shift Timer: %.1f / %.1f", g.shiftTimer, g.shiftDelay), 10, 135, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("Is Shifting: %v", g.isShifting), 10, 160, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("Fade Alpha: %.2f", g.fadeAlpha), 10, 185, 20, rl.Gray)
		// if g.isShifting {
		// rl.DrawText(fmt.Sprintf("Text Progress: %d/%d", g.shiftTextIndex, len(g.shiftText)), 10, 210, 20, rl.Gray)
		// }

		rl.EndDrawing()
	}
}

var lastLightSwitch time.Time
var lastLightningSwitch time.Time

func (g *Game) GetLastRoomPos() (int, int) {
	lastRoom := (*g.world.Map.GetRooms())[len(*g.world.Map.GetRooms())-1]
	return int(lastRoom.X + lastRoom.Height/2), int(lastRoom.Y + lastRoom.Width/2)
}

func (g *Game) Update(deltaTime float32) {

	// g.world.Update(deltaTime)

	scrollY := rl.GetMouseWheelMove()

	if scrollY != 0 {
		g.camera.Zoom += float32(scrollY) * 0.2
	}

	if g.player.GameHasEnded() {
		// Game over
		return
	}

	// ! FOR DEVELOPMENT
	if rl.IsKeyDown(rl.KeyR) {
		x, y := g.world.SwitchMap()
		g.player.Position = rl.NewVector2(x, y)

		// Reset enemies
		g.enemies.Rooms = g.world.Map.GetRoomsRects()
		g.enemies.ResetEnemies()

		// Reset collectibles
		g.collectiblesManager.ScatterCollectibles(g.world.Map.GetRoomsRects(), g.world.Map)

		// Reset lighting
		// g.lightning.SetUpPropsLightning(g.world.PropsManager.GetProps())

		// Reset minimap
		g.minimap.SetDirty()
	}

	if rl.IsKeyDown(rl.KeyK) {
		helpers.DECAY_FACTOR += 0.25
	}

	if rl.IsKeyDown(rl.KeyJ) {
		helpers.DECAY_FACTOR -= 0.25
	}

	if rl.IsKeyDown(rl.KeyI) {
		helpers.LIGHT_RADIUS += 0.5
	}

	if rl.IsKeyDown(rl.KeyU) {
		helpers.LIGHT_RADIUS -= 0.5
	}

	if rl.IsKeyDown(rl.KeyL) {
		if time.Since(lastLightSwitch) > time.Second {
			g.flags ^= RENDER_LIGHTING
			lastLightSwitch = time.Now()
		}
	}

	if g.flags&RENDER_LIGHTING != 0 {
		g.lightning.Update()
	}

	if rl.IsKeyDown(rl.KeyP) {
		if time.Since(lastLightningSwitch) > time.Second {
			g.lightning.NextLightningMode()
			lastLightningSwitch = time.Now()
		}
	}

	//! END DEVELOPMENT

	// Handle pause toggle
	if rl.IsKeyPressed(rl.KeyY) {
		g.isPaused = !g.isPaused
	}

	g.world.PropsManager.Update(deltaTime)
	g.player.Update(deltaTime)
	g.enemies.Update(deltaTime, g.player)
	g.collectiblesManager.Update(deltaTime)

	// PATH FINDING
	// helpers.DEBUG("PLAYER POS", g.player.Position)
	// g.world.Pathfinder.Update(
	// 	int(g.player.Position.X/helpers.TILE_SIZE),
	// 	int(g.player.Position.Y/helpers.TILE_SIZE),
	// 	int(g.world.Map.GetRoomsRects()[len(g.world.Map.GetRoomsRects())-1].X),
	// 	int(g.world.Map.GetRoomsRects()[len(g.world.Map.GetRoomsRects())-1].Y),
	// )

	// g.teleportTimer.Update(deltaTime)
	// g.transition.Update()

	// if g.teleportTimer.lastTickTime == 0 {
	// 	g.teleportTimer.TriggerTick(g.soundManager)
	// }

	// if g.teleportTimer.TeleportTriggered() {
	// 	g.soundManager.PlayBiwaSound()
	// 	g.transition.Start()

	// 	if !g.transition.active {
	// 		g.TeleportPlayer()
	// 		g.teleportTimer.Reset(15, 30) // Reset for next teleport
	// 	}
	// }

	// Toggle map view with T key
	if rl.IsKeyPressed(rl.KeyT) {
		g.minimap.ToggleView()
	}

	if destX, destY := g.GetLastRoomPos(); destX != -1 {
		g.minimap.SetDestination(destX, destY)
	}

	g.minimap.Update(g.player.GetPosition())
	println(g.lightning.Count(), len(*g.world.PropsManager.GetProps()))

	// Update dungeon shift timer
	if !g.isShifting {
		g.shiftTimer += deltaTime
		if g.shiftTimer >= g.shiftDelay-5 { // Start effect 5 seconds before shift
			g.lightning.SetMode("heartbeat") // Set to HandleGlitchLighting (or any other mode you prefer)
		}
		if g.shiftTimer >= g.shiftDelay {
			g.isShifting = true
			g.fadeAlpha = 0
		}
	} else {

		// Handle the shift transition
		const fadeSpeed = 1.0
		const textDuration = 2.0 // Show text for 2 seconds

		if g.shiftTextTimer < textDuration {
			if !g.shiftSoundPlayed {
				g.soundManager.RequestSound("biwa", 1.0, 1.0) // Adjust sound name as needed
				g.shiftSoundPlayed = true
			}
			// First phase: fade to black
			g.fadeAlpha += fadeSpeed * deltaTime
			if g.fadeAlpha >= 1.0 {
				g.fadeAlpha = 1.0
				g.shiftTextTimer += deltaTime
			}
		} else if g.shiftTextTimer >= textDuration && g.shiftTextTimer < textDuration+0.1 {
			// Perform the actual shift exactly once

			x, y := g.world.SwitchMap()
			g.player.Position = rl.NewVector2(x, y)
			g.enemies.Rooms = g.world.Map.GetRoomsRects()
			g.enemies.ResetEnemies()
			g.collectiblesManager.ScatterCollectibles(g.world.Map.GetRoomsRects(), g.world.Map)
			g.lightning.SetUpPropsLightning(g.world.PropsManager.GetProps())
			g.lightning.SetMode("static") // Reset to default lighting mode
			g.minimap.SetDirty()
			g.shiftTextTimer = textDuration + 0.1
			g.shiftSoundPlayed = false
		} else {
			// Final phase: fade back in
			g.fadeAlpha -= fadeSpeed * deltaTime
			if g.fadeAlpha <= 0 {
				// Reset for next shift
				g.isShifting = false
				g.shiftTimer = 0
				g.shiftTextTimer = 0
				g.fadeAlpha = 0
				g.shiftDelay = float32(40 + rand.Intn(41)) // Random value between 40 and 80 seconds
				// g.shiftDelay = float32(1 + rand.Intn(3)) // Random value between 40 and 80 seconds
			}
		}
	}
}

func (g *Game) Render() {
	rl.BeginMode2D(g.camera)
	g.world.Render()
	g.enemies.Render()
	g.player.Render()
	// g.transition.Render()
	// g.world.Pathfinde<r.Render()
	// //
	g.world.Pathfinder.Render(
		g.player.GetPlayerRoom(),
	)

	if g.flags&RENDER_LIGHTING != 0 {
		g.lightning.Render()
	}
	g.collectiblesManager.Render()
	rl.EndMode2D()

	// Render minimap after EndMode2D so it stays fixed on screen
	g.minimap.Render(g.player.Position)
	g.player.RenderHearts()

	// Render shift transition effects
	if g.isShifting {
		// Draw darkening overlay
		rl.DrawRectangle(0, 0, int32(g.width), int32(g.height),
			rl.ColorAlpha(rl.Black, g.fadeAlpha))

		// Draw text if faded to black
		if g.fadeAlpha >= 1.0 {
			fontSize := int32(30)
			textWidth := rl.MeasureText(g.shiftText, fontSize)
			textX := int32(g.width/2) - textWidth/2
			textY := int32(g.height / 2)

			rl.DrawText(g.shiftText, textX, textY, fontSize, rl.White)
		}
	}

	// Debug information
	rl.DrawText(fmt.Sprintf("Shift Timer: %.1f / %.1f", g.shiftTimer, g.shiftDelay), 10, 135, 20, rl.Gray)
	rl.DrawText(fmt.Sprintf("Is Shifting: %v", g.isShifting), 10, 160, 20, rl.Gray)
	rl.DrawText(fmt.Sprintf("Fade Alpha: %.2f", g.fadeAlpha), 10, 185, 20, rl.Gray)
	// if g.isShifting {
	// 	rl.DrawText(fmt.Sprintf("Text Progress: %d/%d", g.shiftTextIndex, len(g.shiftText)), 10, 210, 20, rl.Gray)
	// }
}

func (g *Game) checkGameEnd() bool {
	// For now, just check player's game end condition
	if g.player.GameHasEnded() {
		g.showOutro = true
		return true
	}
	return false
}
