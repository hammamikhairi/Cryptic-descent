package effects

import (
	"crydes/player"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RetroLightingEffect struct {
	lightMask   rl.RenderTexture2D
	lightRadius float32
	pixelSize   int32
	player      *player.Player
	smoothness  float32
}

func NewRetroLightingEffect(screenWidth, screenHeight int32, lightRadius float32, pixelSize int32, p *player.Player) *RetroLightingEffect {
	return &RetroLightingEffect{
		lightMask:   rl.LoadRenderTexture(screenWidth, screenHeight),
		lightRadius: lightRadius,
		pixelSize:   pixelSize,
		player:      p,
		smoothness:  0.8, // Adjust this value between 0 and 1 for desired smoothness
	}
}

func (rle *RetroLightingEffect) Update() {
	rl.BeginTextureMode(rle.lightMask)
	rl.ClearBackground(rl.Black)

	playerCenter := rle.player.GetPlayerCenterPoint()

	for x := int32(0); x < int32(rle.lightMask.Texture.Width); x += rle.pixelSize {
		for y := int32(0); y < int32(rle.lightMask.Texture.Height); y += rle.pixelSize {
			dx := float32(x) - playerCenter.X
			dy := float32(y) - playerCenter.Y
			distance := rl.Vector2Length(rl.Vector2{X: dx, Y: dy})

			if distance < rle.lightRadius {
				intensity := 1.0 - (distance / rle.lightRadius)
				intensity = float32(math.Pow(float64(intensity), float64(rle.smoothness))) // Smooth transition
				quantizedIntensity := float32(int32(intensity*8)) / 8                      // Increase quantization levels
				finalIntensity := intensity*rle.smoothness + quantizedIntensity*(1-rle.smoothness)
				color := rl.ColorAlpha(rl.White, finalIntensity)
				rl.DrawRectangle(x, y, rle.pixelSize, rle.pixelSize, color)
			}
		}
	}

	rl.EndTextureMode()
}

func (rle *RetroLightingEffect) Render() {
	rl.BeginBlendMode(rl.BlendMultiplied)
	rl.DrawTextureRec(rle.lightMask.Texture,
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(rle.lightMask.Texture.Width),
			Height: float32(-rle.lightMask.Texture.Height),
		},
		rl.Vector2{X: 0, Y: 0},
		rl.White)
	rl.EndBlendMode()
}

func (rle *RetroLightingEffect) Unload() {
	rl.UnloadRenderTexture(rle.lightMask)
}
