package main

import (
	"crydes/audio"
	"crydes/core"
	"crydes/effects"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(800, 800, "Gale Game with Raylib and Go")
	defer rl.CloseWindow()

	soundManager := audio.NewSoundManager()
	defer soundManager.Unload() // Ensure sounds are unloaded properly

	transition := effects.NewTransition(0.01) // Create a transition effect

	game := core.NewGame(soundManager, transition) // Pass sound manager and transition to the game
	game.Run(1000, 1000)                           // Start the game loop
}
