package helpers

import (
	"image/color"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	SCREEN_WIDTH  int32 = 1500
	SCREEN_HEIGHT int32 = 1000
)

const (
	FULLSCREEN    = false
	TILE_SIZE     = 16 // Smaller tile size for a more compact map
	MAP_WIDTH     = 100
	MAP_HEIGHT    = 100
	MAX_DEPTH     = 5  // Number of divisions for the BSP tree
	MIN_ROOM_SIZE = 5  // Minimum size for a room
	MAX_ROOM_SIZE = 12 // Maximum size for a room

	DAMAGE_DURATION = time.Duration(0.1 * float32(time.Second))

	ENEMIES_PLAYER_RANGE = 100
	ENEMIES_MOV_SPEED    = 0.001
	// ENEMIES_EPSILON                    = 0.001
	ENEMIES_BOUNCE_BACK_DISTANCE       = 6
	ENEMIES_DIRECTION_CHANGE_THRESHOLD = 5.0

	//LIGHTNING

	CAM_ZOOM = 6.5
)

var (
	VOID_COLOR   = color.RGBA{0, 0, 0, 255}
	DECAY_FACTOR = 5.5
	LIGHT_RADIUS = float32(98.0)
	// VOID_COLOR = color.RGBA{0, 0, 0, 255}
)

var DAMAGE_COLOR rl.Color = rl.NewColor(255, 0, 0, 255) // red
