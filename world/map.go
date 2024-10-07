package world

import (
	"math/rand"

	helpers "crydes/helpers"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type cornerType int

const (
	noCorner cornerType = iota
	cornerTR
	cornerTL
	cornerBR
	cornerBL
	innerCornerTR
	innerCornerTL
	innerCornerBR
	innerCornerBL
)

type wallDirection int

const (
	noWall wallDirection = iota
	wallTop
	wallBottom
	wallLeft
	wallRight
)

type Textures struct {
	floorTexture   rl.Texture2D
	wallTexture    rl.Texture2D
	cornersTexture map[string]rl.Texture2D
	wallTextures   map[string]rl.Texture2D
}

type Map struct {
	dungeon [helpers.MAP_WIDTH][helpers.MAP_HEIGHT]int
	rooms   []helpers.Rectangle

	Textures
}

func NewMap() *Map {

	m := &Map{
		rooms:   []helpers.Rectangle{},
		dungeon: [helpers.MAP_WIDTH][helpers.MAP_HEIGHT]int{},
		Textures: Textures{
			cornersTexture: make(map[string]rl.Texture2D),
			wallTextures:   make(map[string]rl.Texture2D),
		},
	}

	m.loadTextures()
	m.initDungeon()
	m.generateDungeon()

	return m
}

// Initialize the dungeon with walls.
func (m *Map) initDungeon() {
	for x := 0; x < helpers.MAP_WIDTH; x++ {
		for y := 0; y < helpers.MAP_HEIGHT; y++ {
			m.dungeon[x][y] = 0 // 1 means wall
		}
	}
}

func (m *Map) FirstRoomPosition() (float32, float32) {
	return float32((m.rooms[0].X + m.rooms[0].Width/2) * helpers.TILE_SIZE), float32((m.rooms[0].Y + m.rooms[0].Height/2) * helpers.TILE_SIZE)
}

func (m *Map) GetRooms() []helpers.Rectangle {
	return m.rooms
}

func (m *Map) SwitchMap() (float32, float32) {
	m.initDungeon()
	m.rooms = []helpers.Rectangle{}
	m.generateDungeon()
	return m.FirstRoomPosition()
}

// Generate the dungeon using BSP.
func (m *Map) generateDungeon() {

	for len(m.rooms) < 3 {
		m.rooms = []helpers.Rectangle{}
		m.bspSplit(helpers.Rectangle{4, 4, helpers.MAP_WIDTH - 8, helpers.MAP_HEIGHT - 8}, 0)
	}

	for _, room := range m.rooms {
		m.carveRoom(room)
	}
	m.connectRooms()
}

// Split the map into rooms using Binary Space Partitioning.
func (m *Map) bspSplit(area helpers.Rectangle, depth int) {
	if depth >= helpers.MAX_DEPTH {
		roomWidth := rand.Intn(helpers.MAX_ROOM_SIZE-helpers.MIN_ROOM_SIZE+1) + helpers.MIN_ROOM_SIZE
		roomHeight := rand.Intn(helpers.MAX_ROOM_SIZE-helpers.MIN_ROOM_SIZE+1) + helpers.MIN_ROOM_SIZE

		if int(area.Width)-roomWidth < 3 || int(area.Height)-roomHeight < 3 {
			return
		}

		roomX := rand.Intn(int(area.Width)-roomWidth-2) + int(area.X) + 1
		roomY := rand.Intn(int(area.Height)-roomHeight-2) + int(area.Y) + 1

		m.rooms = append(m.rooms, helpers.Rectangle{int32(roomX), int32(roomY), int32(roomWidth), int32(roomHeight)})
		return
	}

	if rand.Intn(2) == 0 && area.Width > helpers.MIN_ROOM_SIZE*2 {
		split := rand.Intn(int(area.Width)/2) + int(area.Width/4)
		m.bspSplit(helpers.Rectangle{area.X, area.Y, int32(split), area.Height}, depth+1)
		m.bspSplit(helpers.Rectangle{area.X + int32(split), area.Y, area.Width - int32(split), area.Height}, depth+1)
	} else if area.Height > helpers.MIN_ROOM_SIZE*2 {
		split := rand.Intn(int(area.Height)/2) + int(area.Height/4)
		m.bspSplit(helpers.Rectangle{area.X, area.Y, area.Width, int32(split)}, depth+1)
		m.bspSplit(helpers.Rectangle{area.X, area.Y + int32(split), area.Width, area.Height - int32(split)}, depth+1)
	}
}

// Carve out a room in the dungeon by setting its tiles to 1.
func (m *Map) carveRoom(room helpers.Rectangle) {
	for x := room.X; x < room.X+room.Width; x++ {
		for y := room.Y; y < room.Y+room.Height; y++ {
			m.dungeon[x][y] = 1
		}
	}
}

// Connect rooms with corridors.
func (m *Map) connectRooms() {
	for i := 1; i < len(m.rooms); i++ {
		prevRoom := m.rooms[i-1]
		currRoom := m.rooms[i]

		prevCenterX := prevRoom.X + prevRoom.Width/2
		prevCenterY := prevRoom.Y + prevRoom.Height/2
		currCenterX := currRoom.X + currRoom.Width/2
		currCenterY := currRoom.Y + currRoom.Height/2

		// Create horizontal corridor between rooms.
		for x := min(prevCenterX, currCenterX); x <= max(prevCenterX, currCenterX); x++ {
			m.dungeon[x][prevCenterY] = 1
		}

		// Create vertical corridor between rooms.
		for y := min(prevCenterY, currCenterY); y <= max(prevCenterY, currCenterY); y++ {
			m.dungeon[currCenterX][y] = 1
		}
	}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func (m *Map) Render() {
	for x := 0; x < helpers.MAP_WIDTH; x++ {
		for y := 0; y < helpers.MAP_HEIGHT; y++ {
			if m.dungeon[x][y] == 1 {
				rl.DrawTexture(m.floorTexture, int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
			} else {
				if valid, corner := m.isDungeonCorner(x, y); valid {
					switch corner {
					case cornerTL:
						rl.DrawTexture(m.cornersTexture["TL"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case cornerBR:
						rl.DrawTexture(m.cornersTexture["BR"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case cornerTR:
						rl.DrawTexture(m.cornersTexture["TR"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case cornerBL:
						rl.DrawTexture(m.cornersTexture["BL"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case innerCornerTL:
						rl.DrawTexture(m.cornersTexture["TLI"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case innerCornerBR:
						rl.DrawTexture(m.cornersTexture["BRI"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case innerCornerTR:
						rl.DrawTexture(m.cornersTexture["TRI"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case innerCornerBL:
						rl.DrawTexture(m.cornersTexture["BLI"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					}
					continue
				}

				if valid, wall := m.isDungeonWall(x, y); valid {
					switch wall {
					case wallTop:
						rl.DrawTexture(m.wallTextures["T"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case wallBottom:
						rl.DrawTexture(m.wallTextures["B"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case wallLeft:
						rl.DrawTexture(m.wallTextures["L"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					case wallRight:
						rl.DrawTexture(m.wallTextures["R"], int32(x*helpers.TILE_SIZE), int32(y*helpers.TILE_SIZE), rl.White)
					}
					continue
				}
			}
		}
	}
}

// Load textures and other resources.
func (m *Map) loadTextures() {
	m.floorTexture = rl.LoadTexture("assets/ground/88.png")

	m.cornersTexture["BR"] = rl.LoadTexture("assets/walls/16.png")
	m.cornersTexture["TL"] = rl.LoadTexture("assets/walls/14.png")
	m.cornersTexture["TR"] = rl.LoadTexture("assets/walls/1.png")
	m.cornersTexture["BL"] = rl.LoadTexture("assets/walls/9.png")

	m.cornersTexture["BRI"] = rl.LoadTexture("assets/walls/16inner.png")
	m.cornersTexture["TLI"] = rl.LoadTexture("assets/walls/14inner.png")
	m.cornersTexture["TRI"] = rl.LoadTexture("assets/walls/1inner.png")
	m.cornersTexture["BLI"] = rl.LoadTexture("assets/walls/9inner.png")

	m.wallTextures["T"] = rl.LoadTexture("assets/walls/4.png")
	m.wallTextures["B"] = rl.LoadTexture("assets/walls/10.png")
	m.wallTextures["L"] = rl.LoadTexture("assets/walls/12.png")
	m.wallTextures["R"] = rl.LoadTexture("assets/walls/2.png")
}

// Unload textures to free up memory.
func (m *Map) unloadTextures() {
	rl.UnloadTexture(m.floorTexture)
	for _, tex := range m.cornersTexture {
		rl.UnloadTexture(tex)
	}
	for _, tex := range m.wallTextures {
		rl.UnloadTexture(tex)
	}
}

type Player struct {
	X, Y float32
}

func (m *Map) isDungeonCorner(x, y int) (bool, cornerType) {
	// Ensure we're checking within valid dungeon bounds
	if x <= 0 || x >= helpers.MAP_WIDTH || y <= 0 || y >= helpers.MAP_HEIGHT {
		return false, noCorner
	}

	// Check for corner conditions
	// Top-left corner
	if x > 0 && y > 0 && m.dungeon[x-1][y-1] == 1 && m.dungeon[x-1][y] == 0 && m.dungeon[x][y-1] == 0 {
		return true, cornerBR
	}
	// Top-right corner
	if x < helpers.MAP_WIDTH-1 && y > 0 && m.dungeon[x+1][y-1] == 1 && m.dungeon[x+1][y] == 0 && m.dungeon[x][y-1] == 0 {
		return true, cornerBL
	}
	// Bottom-left corner
	if x > 0 && y < helpers.MAP_HEIGHT-1 && m.dungeon[x-1][y+1] == 1 && m.dungeon[x-1][y] == 0 && m.dungeon[x][y+1] == 0 {
		return true, cornerTR
	}
	// Bottom-right corner
	if x < helpers.MAP_WIDTH-1 && y < helpers.MAP_HEIGHT-1 && m.dungeon[x+1][y+1] == 1 && m.dungeon[x+1][y] == 0 && m.dungeon[x][y+1] == 0 {
		return true, cornerTL
	}

	// or 3 sides are 1s // INNER ONES
	if x > 0 && y > 0 && m.dungeon[x-1][y-1] == 1 && m.dungeon[x-1][y] == 1 && m.dungeon[x][y-1] == 1 {
		return true, innerCornerTL
	}

	if x < helpers.MAP_WIDTH-1 && y > 0 && m.dungeon[x+1][y-1] == 1 && m.dungeon[x+1][y] == 1 && m.dungeon[x][y-1] == 1 {
		return true, innerCornerTR
	}

	if x > 0 && y < helpers.MAP_HEIGHT-1 && m.dungeon[x-1][y+1] == 1 && m.dungeon[x-1][y] == 1 && m.dungeon[x][y+1] == 1 {
		return true, innerCornerBL
	}

	if x < helpers.MAP_WIDTH-1 && y < helpers.MAP_HEIGHT-1 && m.dungeon[x+1][y+1] == 1 && m.dungeon[x+1][y] == 1 && m.dungeon[x][y+1] == 1 {
		return true, innerCornerBR
	}

	// If none of the corner conditions are met, return false
	return false, noCorner
}

func (m *Map) isDungeonWall(x, y int) (bool, wallDirection) {
	// Ensure we're checking within valid dungeon bounds
	if x <= 0 || x >= helpers.MAP_WIDTH-1 || y <= 0 || y >= helpers.MAP_HEIGHT-1 {
		return false, noWall
	}

	// if only 1 side is 1

	// Check for wall conditions

	if x > 0 && m.dungeon[x-1][y] == 1 && m.dungeon[x+1][y] == 0 {
		return true, wallLeft
	}

	if x < helpers.MAP_WIDTH-1 && m.dungeon[x+1][y] == 1 && m.dungeon[x-1][y] == 0 {
		return true, wallRight
	}

	if y > 0 && m.dungeon[x][y-1] == 1 && m.dungeon[x][y+1] == 0 {
		return true, wallBottom
	}

	if y < helpers.MAP_HEIGHT-1 && m.dungeon[x][y+1] == 1 && m.dungeon[x][y-1] == 0 {
		return true, wallTop
	}

	// If none of the wall conditions are met, return false
	return false, noWall
}

// IsWalkable checks if a map tile is walkable.
func (m *Map) IsWalkable(x, y int) bool {
	// Check boundaries first
	if x < 0 || x >= helpers.MAP_WIDTH || y < 0 || y >= helpers.MAP_HEIGHT {
		return false
	}
	// Check if the tile is walkable (assuming 0 is not walkable)
	return m.dungeon[x][y] != 0
}

// IsWalkable checks if a map tile is walkable.
func (m *Map) IsWalkableFloat(x, y float32) bool {
	// Check boundaries first
	if int(x) < 0 || int(x) >= helpers.MAP_WIDTH || int(y)/helpers.TILE_SIZE < 0 || int(y)/helpers.TILE_SIZE >= helpers.MAP_HEIGHT {
		return false
	}

	// Check if the tile is walkable / by TILE_SIZE to get the tile position
	return m.dungeon[int(x)/helpers.TILE_SIZE][int(y)/helpers.TILE_SIZE] != 0
}
