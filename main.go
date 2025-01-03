package main

import (
	"crydes/audio"
	"crydes/core"
	"crydes/helpers"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Set up monitor info for fullscreen
	var screenWidth, screenHeight int32
	if helpers.FULLSCREEN {
		monitor := rl.GetCurrentMonitor()
		screenWidth = int32(rl.GetMonitorWidth(monitor))
		screenHeight = int32(rl.GetMonitorHeight(monitor))
		helpers.SCREEN_WIDTH = int32(rl.GetMonitorWidth(monitor))
		helpers.SCREEN_HEIGHT = int32(rl.GetMonitorHeight(monitor))
	} else {
		screenWidth = helpers.SCREEN_WIDTH
		screenHeight = helpers.SCREEN_HEIGHT
	}

	// Initialize window with proper flags
	if helpers.FULLSCREEN {
		rl.SetConfigFlags(rl.FlagFullscreenMode)
	}
	rl.InitWindow(screenWidth, screenHeight, "Cryptic Descent")
	defer rl.CloseWindow()

	// Toggle fullscreen with Alt+Enter
	rl.SetExitKey(0) // Disable exit on ESC

	soundManager := audio.NewSoundManager()
	defer soundManager.Unload()

	game := core.NewGame(soundManager, int(screenWidth), int(screenHeight))

	for !rl.WindowShouldClose() {
		// Handle fullscreen toggle
		if rl.IsKeyPressed(rl.KeyEnter) && (rl.IsKeyDown(rl.KeyLeftAlt) || rl.IsKeyDown(rl.KeyRightAlt)) {
			if rl.IsWindowFullscreen() {
				rl.ToggleFullscreen()
				rl.SetWindowSize(int(helpers.SCREEN_WIDTH), int(helpers.SCREEN_HEIGHT))
			} else {
				monitor := rl.GetCurrentMonitor()
				rl.SetWindowSize(rl.GetMonitorWidth(monitor), rl.GetMonitorHeight(monitor))
				rl.ToggleFullscreen()
			}
		}

		game.Run()
	}
}
