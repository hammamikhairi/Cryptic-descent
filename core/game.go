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

	pauseScreen *screens.PauseScreen
	titleScreen *screens.TitleScreen
	isPaused    bool
	showTitle   bool

	minimap             *minimap.Minimap
	collectiblesManager *world.CollectibleManager
}

// NewGame initializes a new game instance
func NewGame(soundManager *audio.SoundManager, transition *effects.Transition, width, height int) *Game {

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
		transition:          transition,
		enemies:             em,
		width:               width,
		height:              height,
		lightning:           rle,
		flags:               0,
		pauseScreen:         screens.NewPauseScreen(soundManager),
		titleScreen:         screens.NewTitleScreen(soundManager),
		isPaused:            false,
		showTitle:           true,
		minimap:             mm,
		collectiblesManager: collectibleManager,
	}

}

func (g *Game) Run() {
	rl.SetTargetFPS(45)
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
		if !g.isPaused && !g.showTitle {
			g.Update(deltaTime)
		} else if g.showTitle {
			if g.titleScreen.Update(deltaTime) {
				g.showTitle = false
			}
		} else if g.isPaused {
			if g.pauseScreen.Update(deltaTime) {
				g.isPaused = false
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
		}

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

func (g *Game) GetLastRoomPos() (int, int) {
	rm := g.world.Map.GetRoomsRects()[len(g.world.Map.GetRoomsRects())-1]
	return int(rm.X + rm.Height/2), int(rm.Y + rm.Width/2)
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
		g.lightning.SetUpPropsLightning(g.world.PropsManager.GetProps())

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
}

func (g *Game) Render() {
	rl.BeginMode2D(g.camera)
	g.world.Render()
	g.enemies.Render()
	g.player.Render()
	g.transition.Render()
	// g.world.Pathfinde<r.Render()
	// //
	// g.world.Pathfinder.Render(
	// 	g.player.GetPlayerRoom(),
	// )

	if g.flags&RENDER_LIGHTING != 0 {
		g.lightning.Render()
	}
	g.collectiblesManager.Render()
	rl.EndMode2D()

	// Render minimap after EndMode2D so it stays fixed on screen
	g.minimap.Render(g.player.Position)
	g.player.RenderHearts()
}
