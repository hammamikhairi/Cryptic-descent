package helpers

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

var LOGS map[int]bool = make(map[int]bool)

func LogOnce(id int, msg string) {
	if val, _ := LOGS[id]; !val {
		fmt.Printf("[LOG %d] %s\n", id, msg)
		LOGS[id] = true
	}
}

func ABS(value float32) float32 {
	if value < 0 {
		return -value
	}
	return value
}

func CheckCollisionRecs(r1, r2 rl.Rectangle) bool {
	return r1.X < r2.X+r2.Width && r1.X+r1.Width > r2.X && r1.Y < r2.Y+r2.Height && r1.Y+r1.Height > r2.Y
}

func DEBUG(tag string, msg any) {
	fmt.Printf("[DEBUG] %s : %+v\n", tag, msg)
}

func (r1 *Rectangle) Intersects(r2 Rectangle) bool {
	return r1.X < r2.X+r2.Width && r1.X+r1.Width > r2.X && r1.Y < r2.Y+r2.Height && r1.Y+r1.Height > r2.Y
}

func (r1 *Rectangle) GetRandomPosInRect() rl.Vector2 {
	enemyPos := rl.NewVector2(TILE_SIZE*(float32(r1.X)+float32(rl.GetRandomValue(0, int32(r1.Width)))), TILE_SIZE*(float32(r1.Y)+float32(rl.GetRandomValue(0, int32(r1.Height)))))
	return enemyPos
}

func (r1 *Rectangle) ContainsPoint(p rl.Vector2) bool {
	// DEBUG("CHECKING POINT", p)
	// DEBUG("RECTANGLE", r1)
	return p.X >= float32(r1.X)*TILE_SIZE && p.X <= (float32(r1.X+r1.Width)*TILE_SIZE) && p.Y >= float32(r1.Y)*TILE_SIZE && p.Y <= float32(r1.Y+r1.Height)*TILE_SIZE
}

func GetRandomEnemyType() string {
	return ENEMY_TYPES[int(rl.GetRandomValue(0, int32(len(ENEMY_TYPES)-1)))]
}
