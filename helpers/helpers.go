package helpers

import rl "github.com/gen2brain/raylib-go/raylib"

func LoadAnimation(id string, filePaths ...string) *Animation {

	var textures []rl.Texture2D = []rl.Texture2D{}
	for _, path := range filePaths {
		texture := rl.LoadTexture(path)

		if texture.ID == 0 {
			rl.TraceLog(rl.LogError, "Failed to load texture: %s", path)
			panic("Failed to load texture")
		}

		textures = append(textures, texture)
	}

	return &Animation{
		Frames:       textures, // You can load multiple frames here if needed
		FrameTime:    0.1,      // Adjust frame time for animation speed
		Timer:        0,
		CurrentFrame: 0,
		ID:           id,
	}
}

func GetDistance(a, b rl.Vector2) float32 {
	return rl.Vector2Distance(a, b)
}
