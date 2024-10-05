package helpers

import rl "github.com/gen2brain/raylib-go/raylib"

type Animation struct {
	Frames       []rl.Texture2D
	CurrentFrame int
	FrameTime    float32
	Timer        float32
	ID           string
}
