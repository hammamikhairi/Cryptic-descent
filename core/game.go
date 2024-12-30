package core

import (
	"crydes/audio"
	"crydes/effects"
	"crydes/enemies"
	"crydes/helpers"
	"crydes/player"
	"crydes/world"
	"time"

	"fmt"

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
	transition   *effects.Transition // Reference to the transition effect

	enemies *enemies.EnemiesManager

	camera        rl.Camera2D
	width, height int

	// ! FOR DEVELOPMENT
	flags int
	//! END DEVELOPMENT

}

// NewGame initializes a new game instance
func NewGame(soundManager *audio.SoundManager, transition *effects.Transition, width, height int) *Game {

	w := world.NewWorld()

	x, y := w.PlayerSpawn()
	// x, y := w.Map.GetRandomRoomCenterBySize(1)
	p := player.NewPlayer(x, y, w.Map)

	em := enemies.NewEnemiesManager(x, y, w.Map, p.AttackChan, w.Map.GetRoomsRects())
	em.SpawnEnemies()

	rle := effects.NewRetroLightingEffect(
		int32(helpers.MAP_WIDTH*helpers.TILE_SIZE), int32(helpers.MAP_HEIGHT*helpers.TILE_SIZE), 50, 2, p,
	)

	rle.SetUpPropsLightning(w.PropsManager.GetProps())

	return &Game{
		player:       p,
		world:        w,
		soundManager: soundManager,
		transition:   transition,
		enemies:      em,
		width:        width,
		height:       height,
		lightning:    rle,
		flags:        0,
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

		g.Update(deltaTime)

		helpers.LogOnce(1, "HELLOOOO")

		// Update camera target to follow the player
		g.camera.Target = rl.Vector2{X: float32(g.player.Position.X), Y: float32(g.player.Position.Y)}

		// Ensure the camera offset stays centered even if window size changes
		g.camera.Offset = rl.Vector2{X: float32(rl.GetScreenWidth()) / 2, Y: float32(rl.GetScreenHeight()) / 2}

		rl.BeginDrawing()
		rl.ClearBackground(helpers.VOID_COLOR) // rgb(88, 68, 34)

		rl.BeginMode2D(g.camera)
		g.Render()
		rl.EndMode2D()
		rl.DrawText("Cryptic Descent", 10, 10, 20, rl.Gray)
		fpsText := fmt.Sprintf("FPS: %d", rl.GetFPS())
		rl.DrawText(fpsText, 10, 35, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("CAM ZOOM : %.3f", g.camera.Zoom), 10, 60, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("DECAY FACTOR : %.3f", helpers.DECAY_FACTOR), 10, 85, 20, rl.Gray)
		rl.DrawText(fmt.Sprintf("LIGHT RADIUS : %.3f", helpers.LIGHT_RADIUS), 10, 110, 20, rl.Gray)

		rl.EndDrawing()
	}
}

var lastLightSwitch time.Time
var lastLightningSwitch time.Time

func (g *Game) Update(deltaTime float32) {

	g.world.Update(deltaTime)

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
		g.enemies.Rooms = g.world.Map.GetRoomsRects()
		g.enemies.ResetEnemies()
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

	g.player.Update(deltaTime)
	g.enemies.Update(deltaTime, g.player)

	// PATH FINDING
	helpers.DEBUG("PLAYER POS", g.player.Position)
	g.world.Pathfinder.Update(
		int(g.player.Position.X/helpers.TILE_SIZE),
		int(g.player.Position.Y/helpers.TILE_SIZE),
		int(g.world.Map.GetRoomsRects()[len(g.world.Map.GetRoomsRects())-1].X),
		int(g.world.Map.GetRoomsRects()[len(g.world.Map.GetRoomsRects())-1].Y),
	)

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
	// g.world.Pathfinder.Render()

	g.world.Pathfinder.Render(
		g.player.GetPlayerRoom(),
	)

	if g.flags&RENDER_LIGHTING != 0 {
		g.lightning.Render()
	}
}
