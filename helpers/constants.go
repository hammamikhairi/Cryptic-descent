package helpers

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 800
	TILE_SIZE     = 16 // Smaller tile size for a more compact map
	MAP_WIDTH     = 80
	MAP_HEIGHT    = 80
	MAX_DEPTH     = 5  // Number of divisions for the BSP tree
	MIN_ROOM_SIZE = 5  // Minimum size for a room
	MAX_ROOM_SIZE = 12 // Maximum size for a room

	DAMAGE_DURATION = time.Duration(0.1 * float32(time.Second))
)

var DAMAGE_COLOR rl.Color = rl.NewColor(255, 0, 0, 255) // red
