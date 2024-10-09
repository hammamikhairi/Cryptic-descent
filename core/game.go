package core

import (
	"crydes/audio"
	"crydes/effects"
	"crydes/enemies"
	"crydes/helpers"
	"crydes/player"
	"crydes/world"

	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	player    *player.Player
	world     *world.World
	lightning *effects.RetroLightingEffect

	soundManager *audio.SoundManager // Reference to the sound manager
	transition   *effects.Transition // Reference to the transition effect

	enemies *enemies.EnemiesManager

	camera        rl.Camera2D
	width, height int
}

// NewGame initializes a new game instance
func NewGame(soundManager *audio.SoundManager, transition *effects.Transition, width, height int) *Game {

	w := world.NewWorld()

	x, y := w.PlayerSpawn()
	p := player.NewPlayer(x, y, w.Map)

	w.NewProp(1, x+10, y+10, 0.7, rl.NewVector2(16, 16),
		helpers.LoadAnimation(
			"assets/fireplace/1.png",
			"assets/fireplace/2.png",
			"assets/fireplace/3.png",
			"assets/fireplace/4.png",
		),
		true)

	em := enemies.NewEnemiesManager(x, y, w.Map, p.AttackChan, w.Map.GetRooms())
	em.SpawnEnemies()

	return &Game{
		player:       p,
		world:        w,
		soundManager: soundManager,
		transition:   transition,
		enemies:      em,
		width:        width,
		height:       height,
		lightning: effects.NewRetroLightingEffect(
			int32(helpers.MAP_WIDTH*helpers.TILE_SIZE), int32(helpers.MAP_HEIGHT*helpers.TILE_SIZE), 50, 2, p,
		),
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
		Zoom:     5.0,
	}

	for !rl.WindowShouldClose() {
		deltaTime := float32(rl.GetTime() - previousTime)
		previousTime = rl.GetTime()

		g.Update(deltaTime)

		helpers.LogOnce(1, "HELLOOOO")

		// Update camera target to follow the player
		g.camera.Target = rl.Vector2{X: float32(g.player.Position.X), Y: float32(g.player.Position.Y)}

		// Ensure the camera offset stays centered even if window size changes
		g.camera.Offset = rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2}

		rl.BeginDrawing()
		rl.ClearBackground(color.RGBA{88, 68, 34, 255}) // rgb(88, 68, 34)

		rl.BeginMode2D(g.camera)
		g.Render()
		rl.EndMode2D()
		rl.DrawText("Cryptic Descent", 10, 10, 20, rl.Black)
		fpsText := fmt.Sprintf("FPS: %d", rl.GetFPS())
		rl.DrawText(fpsText, 10, 35, 20, rl.Black)

		rl.EndDrawing()
	}
}

func (g *Game) Update(deltaTime float32) {

	g.world.Update(deltaTime)

	scrollY := rl.GetMouseWheelMove()

	if scrollY != 0 {
		g.camera.Zoom += float32(scrollY) * 0.1
	}

	if g.player.GameHasEnded() {
		// Game over
		return
	}

	g.player.Update(deltaTime)

	if rl.IsKeyDown(rl.KeyR) {
		x, y := g.world.SwitchMap()
		g.player.Position = rl.NewVector2(x, y)
		g.enemies.Rooms = g.world.Map.GetRooms()
		g.enemies.ResetEnemies()
	}

	g.lightning.Update()
	g.enemies.Update(deltaTime, g.player)

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
}

func (g *Game) Render() {
	g.world.Render()
	g.enemies.Render()
	g.player.Render() // Render player here
	g.transition.Render()
	g.lightning.Render()
}
