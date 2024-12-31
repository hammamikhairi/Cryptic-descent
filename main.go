package main

import (
	"crydes/audio"
	"crydes/core"
	"crydes/effects"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	height = 800
	width  = 800
)

func main() {
	rl.InitWindow(height, width, "Cryptic Descent")
	defer rl.CloseWindow()

	soundManager := audio.NewSoundManager()
	defer soundManager.Unload() // Ensure sounds are unloaded properly

	transition := effects.NewTransition(0.01) // Create a transition effect

	game := core.NewGame(soundManager, transition, height, width) // Pass sound manager and transition to the game
	game.Run()                                                    // Start the game loop
}
