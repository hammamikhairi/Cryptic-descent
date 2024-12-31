package world

import (
	"math"
	"math/rand"
	"sort"

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
	cornersTexture map[string]rl.Texture2D
	wallTextures   map[string]rl.Texture2D
}

// 0 means not walkable, 1 means walkable
type Map struct {
	dungeon [helpers.MAP_WIDTH][helpers.MAP_HEIGHT]int
	rooms   []*Room
	// corridors [][]rl.Vector2

	Textures
}

type RoomSize int

const (
	SmallRoom RoomSize = iota
	MediumRoom
	LargeRoom
)

type Room struct {
	helpers.Rectangle
	Size RoomSize
}

func NewMap() *Map {
	m := &Map{
		rooms:   []*Room{},
		dungeon: [helpers.MAP_WIDTH][helpers.MAP_HEIGHT]int{},
		Textures: Textures{
			cornersTexture: make(map[string]rl.Texture2D),
			wallTextures:   make(map[string]rl.Texture2D),
		},
	}

	m.loadTextures()
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
	// Choose a random room that's not the last room
	roomIndex := rand.Intn(len(m.rooms) - 1)
	room := m.rooms[roomIndex]

	// Move the chosen room to the front of the slice
	m.rooms[0], m.rooms[roomIndex] = m.rooms[roomIndex], m.rooms[0]

	return float32((room.X + room.Width/2) * helpers.TILE_SIZE),
		float32((room.Y + room.Height/2) * helpers.TILE_SIZE)
}

func (m *Map) GetRoomsRects() []helpers.Rectangle {
	rectangles := make([]helpers.Rectangle, len(m.rooms))
	for i, room := range m.rooms {
		rectangles[i] = room.Rectangle
	}
	return rectangles
}

func (m *Map) GetRooms() *[]*Room {
	return &m.rooms
}

func (m *Map) SwitchMap() (float32, float32) {
	m.initDungeon()
	m.rooms = []*Room{}
	m.generateDungeon()
	return m.FirstRoomPosition()
}

// Generate the dungeon using BSP.
func (m *Map) generateDungeon() {
	m.initDungeon()
	for len(m.rooms) < 3 {
		m.rooms = []*Room{}
		m.bspSplit(helpers.Rectangle{X: 1, Y: 1, Width: helpers.MAP_WIDTH - 2, Height: helpers.MAP_HEIGHT - 2}, 0)
	}

	for _, room := range m.rooms {
		m.carveRoom(room.Rectangle)
	}
	m.connectRooms()
}

func (m *Map) GetRoomsBySize(size int) []Room {
	var result []Room
	for _, room := range m.rooms {
		if room.Size == RoomSize(size) {
			result = append(result, *room)
		}
	}
	return result
}

// Split the map into rooms using Binary Space Partitioning.
func (m *Map) bspSplit(area helpers.Rectangle, depth int) {
	if depth >= helpers.MAX_DEPTH {
		roomSize := m.chooseRoomSize(depth)
		roomWidth, roomHeight := m.getRoomDimensions(roomSize)

		if int(area.Width)-int(roomWidth) < 3 || int(area.Height)-int(roomHeight) < 3 {
			return
		}

		maxAttempts := 100
		for attempt := 0; attempt < maxAttempts; attempt++ {
			roomX := int(area.X) + (int(area.Width)-int(roomWidth))/2 + rand.Intn(3) - 1
			roomY := int(area.Y) + (int(area.Height)-int(roomHeight))/2 + rand.Intn(3) - 1

			newRoom := Room{
				Rectangle: helpers.Rectangle{X: int32(roomX), Y: int32(roomY), Width: roomWidth, Height: roomHeight},
				Size:      roomSize,
			}

			if m.isValidRoomPlacement(newRoom.Rectangle) {
				m.rooms = append(m.rooms, &newRoom)
				return
			}
		}
		return
	}

	splitRatio := 0.4 + rand.Float64()*0.2

	if rand.Intn(2) == 0 && area.Width > helpers.MIN_ROOM_SIZE*2 {
		split := int(float64(area.Width) * splitRatio)
		m.bspSplit(helpers.Rectangle{area.X, area.Y, int32(split), area.Height}, depth+1)
		m.bspSplit(helpers.Rectangle{area.X + int32(split), area.Y, area.Width - int32(split), area.Height}, depth+1)
	} else if area.Height > helpers.MIN_ROOM_SIZE*2 {
		split := int(float64(area.Height) * splitRatio)
		m.bspSplit(helpers.Rectangle{area.X, area.Y, area.Width, int32(split)}, depth+1)
		m.bspSplit(helpers.Rectangle{area.X, area.Y + int32(split), area.Width, area.Height - int32(split)}, depth+1)
	}
}

func (m *Map) chooseRoomSize(depth int) RoomSize {
	if depth > 5 {
		return RoomSize(rand.Intn(3))
	}
	if rand.Float32() < 0.6 {
		return SmallRoom
	}
	return MediumRoom
}

func (m *Map) getRoomDimensions(size RoomSize) (width, height int32) {
	switch size {
	case SmallRoom:
		width = int32(rand.Intn(4) + 3)  // 3-5
		height = int32(rand.Intn(4) + 3) // 3-5
	case MediumRoom:
		width = int32(rand.Intn(4) + 6)  // 6-8
		height = int32(rand.Intn(4) + 6) // 6-8
	case LargeRoom:
		width = int32(rand.Intn(7) + 9)  // 9-13
		height = int32(rand.Intn(7) + 9) // 9-13
	}

	width += 3
	height += 3

	return
}

func (m *Map) isValidRoomPlacement(newRoom helpers.Rectangle) bool {
	for _, room := range m.rooms {
		if room.Intersects(newRoom) {
			return false
		}
	}
	return true
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
	numRooms := len(m.rooms)
	if numRooms == 0 {
		return
	}

	// Create a list of all possible connections
	type connection struct {
		room1, room2 int
		distance     float64
	}

	var connections []connection

	// Generate all possible connections between rooms
	for i := 0; i < numRooms; i++ {
		for j := i + 1; j < numRooms; j++ {
			dist := distanceBetweenRooms(m.rooms[i], m.rooms[j])
			connections = append(connections, connection{i, j, dist})
		}
	}

	// Sort connections by distance
	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distance < connections[j].distance
	})

	// Union-Find data structure for detecting cycles
	parent := make([]int, numRooms)
	for i := range parent {
		parent[i] = i
	}

	// Find with path compression
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	// Union by rank
	union := func(x, y int) {
		parent[find(x)] = find(y)
	}

	// Create minimum spanning tree
	connected := make(map[int]map[int]bool)
	for i := 0; i < numRooms; i++ {
		connected[i] = make(map[int]bool)
	}

	// Connect rooms using Kruskal's algorithm
	for _, conn := range connections {
		if find(conn.room1) != find(conn.room2) {
			union(conn.room1, conn.room2)
			connected[conn.room1][conn.room2] = true
			connected[conn.room2][conn.room1] = true
			m.createCorridor(m.rooms[conn.room1], m.rooms[conn.room2])
		}
	}

	// Add a few extra connections for loops (optional)
	for _, conn := range connections {
		if !connected[conn.room1][conn.room2] && rand.Float64() < 0.2 { // 20% chance for extra connections
			m.createCorridor(m.rooms[conn.room1], m.rooms[conn.room2])
			connected[conn.room1][conn.room2] = true
			connected[conn.room2][conn.room1] = true
		}
	}
}

func (m *Map) createCorridor(room1, room2 *Room) {
	// Get room centers
	start := rl.Vector2{
		X: float32(room1.X + room1.Width/2),
		Y: float32(room1.Y + room1.Height/2),
	}
	end := rl.Vector2{
		X: float32(room2.X + room2.Width/2),
		Y: float32(room2.Y + room2.Height/2),
	}

	// Create 2-3 control points for more organic paths
	numPoints := rand.Intn(2) + 2
	controlPoints := make([]rl.Vector2, numPoints)
	controlPoints[0] = start
	controlPoints[numPoints-1] = end

	// Generate intermediate control points
	for i := 1; i < numPoints-1; i++ {
		controlPoints[i] = rl.Vector2{
			X: start.X + (end.X-start.X)*float32(i)/float32(numPoints-1) + float32(rand.Intn(5)-2),
			Y: start.Y + (end.Y-start.Y)*float32(i)/float32(numPoints-1) + float32(rand.Intn(5)-2),
		}
	}

	// Carve paths through all control points
	for i := 0; i < len(controlPoints)-1; i++ {
		m.carvePath(controlPoints[i], controlPoints[i+1])
	}
}

func (m *Map) carvePath(start, end rl.Vector2) {
	x := int32(start.X)
	y := int32(start.Y)

	// Make the corridor wider at the start
	m.carveArea(x, y, 3)

	for x != int32(end.X) || y != int32(end.Y) {
		if x < int32(end.X) {
			x++
		} else if x > int32(end.X) {
			x--
		}

		if y < int32(end.Y) {
			y++
		} else if y > int32(end.Y) {
			y--
		}

		// Always carve a wider path (minimum width)
		m.carveArea(x, y, 2)

		// Randomly make even wider corridors at some points
		if rand.Float32() < 0.3 {
			m.carveArea(x, y, 3)
		}
	}

	// Make the corridor wider at the end
	m.carveArea(int32(end.X), int32(end.Y), 3)
}

func (m *Map) carveArea(x, y int32, radius int32) {
	// Cap the radius to a maximum value (e.g., 3)
	maxRadius := int32(2)
	if radius > maxRadius {
		radius = maxRadius
	}

	// First pass: carve the main area
	for dx := -radius; dx <= radius; dx++ {
		for dy := -radius; dy <= radius; dy++ {
			newX := x + dx
			newY := y + dy
			if m.isValidPosition(newX, newY) {
				m.dungeon[newX][newY] = 1
			}
		}
	}

	// Second pass: smooth out corners to prevent 1-tile gaps
	for dx := -radius - 1; dx <= radius+1; dx++ {
		for dy := -radius - 1; dy <= radius+1; dy++ {
			newX := x + dx
			newY := y + dy
			if m.isValidPosition(newX, newY) {
				// If surrounded by walkable tiles, make this tile walkable too
				if m.countAdjacentWalkable(newX, newY) >= 5 {
					m.dungeon[newX][newY] = 1
				}
			}
		}
	}
}

func (m *Map) countAdjacentWalkable(x, y int32) int {
	count := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			newX := x + int32(dx)
			newY := y + int32(dy)
			if m.isValidPosition(newX, newY) && m.dungeon[newX][newY] == 1 {
				count++
			}
		}
	}
	return count
}

func (m *Map) isValidPosition(x, y int32) bool {
	return x > 0 && x < helpers.MAP_WIDTH-1 && y > 0 && y < helpers.MAP_HEIGHT-1
}

func distanceBetweenRooms(r1, r2 *Room) float64 {
	c1x := float64(r1.X) + float64(r1.Width)/2
	c1y := float64(r1.Y) + float64(r1.Height)/2
	c2x := float64(r2.X) + float64(r2.Width)/2
	c2y := float64(r2.Y) + float64(r2.Height)/2

	dx := c1x - c2x
	dy := c1y - c2y
	return math.Sqrt(dx*dx + dy*dy)
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

	m.cornersTexture["BR"] = rl.LoadTexture("assets/walls/6.png")
	m.cornersTexture["TL"] = rl.LoadTexture("assets/walls/8.png")
	m.cornersTexture["TR"] = rl.LoadTexture("assets/walls/11.png")
	m.cornersTexture["BL"] = rl.LoadTexture("assets/walls/3.png")

	m.cornersTexture["BRI"] = rl.LoadTexture("assets/walls/16inner.png")
	m.cornersTexture["TLI"] = rl.LoadTexture("assets/walls/14inner.png")
	m.cornersTexture["TRI"] = rl.LoadTexture("assets/walls/1inner.png")
	m.cornersTexture["BLI"] = rl.LoadTexture("assets/walls/9inner.png")

	m.wallTextures["B"] = rl.LoadTexture("assets/walls/4.png")
	m.wallTextures["T"] = rl.LoadTexture("assets/walls/10.png")
	m.wallTextures["R"] = rl.LoadTexture("assets/walls/12.png")
	m.wallTextures["L"] = rl.LoadTexture("assets/walls/2.png")
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
	tileX, tileY := int(x)/helpers.TILE_SIZE, int(y)/helpers.TILE_SIZE
	return m.IsWalkable(tileX, tileY)
}

func (m *Map) CurrentRoomIndex(p rl.Vector2) int {

	for i, room := range m.rooms {
		if room.ContainsPoint(p) {
			return i
		}
	}

	return -1
}

func (r *Room) GetLightPositions() []rl.Vector2 {

	helpers.DEBUG("Room size", r.Size)

	switch r.Size {
	case SmallRoom:
		helpers.DEBUG("Small room", "")
		lights := make([]rl.Vector2, 1)

		centerX := r.X*helpers.TILE_SIZE + r.Width*helpers.TILE_SIZE/2
		centerY := r.Y*helpers.TILE_SIZE + r.Height*helpers.TILE_SIZE/2

		lights[0] = rl.NewVector2(float32(centerX), float32(centerY))
		return lights

	case MediumRoom:
		helpers.DEBUG("Medium room", "")
		lights := make([]rl.Vector2, 4)
		centerX := r.X*helpers.TILE_SIZE + r.Width*helpers.TILE_SIZE/2
		centerY := r.Y*helpers.TILE_SIZE + r.Height*helpers.TILE_SIZE/2
		lightOffsetX := r.Width * helpers.TILE_SIZE / 4
		lightOffsetY := r.Height * helpers.TILE_SIZE / 4

		lights[0] = rl.NewVector2(float32(centerX+lightOffsetX), float32(centerY+lightOffsetY))
		lights[1] = rl.NewVector2(float32(centerX-lightOffsetX), float32(centerY+lightOffsetY))
		lights[2] = rl.NewVector2(float32(centerX+lightOffsetX), float32(centerY-lightOffsetY))
		lights[3] = rl.NewVector2(float32(centerX-lightOffsetX), float32(centerY-lightOffsetY))
		return lights
	case LargeRoom:
		helpers.DEBUG("Large room", "")
		lights := make([]rl.Vector2, 4)
		centerX := r.X*helpers.TILE_SIZE + r.Width*helpers.TILE_SIZE/2
		centerY := r.Y*helpers.TILE_SIZE + r.Height*helpers.TILE_SIZE/2
		lightOffsetX := r.Width * helpers.TILE_SIZE / 4
		lightOffsetY := r.Height * helpers.TILE_SIZE / 4

		lights[0] = rl.NewVector2(float32(centerX+lightOffsetX), float32(centerY+lightOffsetY))
		lights[1] = rl.NewVector2(float32(centerX-lightOffsetX), float32(centerY+lightOffsetY))
		lights[2] = rl.NewVector2(float32(centerX+lightOffsetX), float32(centerY-lightOffsetY))
		lights[3] = rl.NewVector2(float32(centerX-lightOffsetX), float32(centerY-lightOffsetY))
		return lights
	}

	return nil
}

func (mp *Map) GetRandomRoomCenterBySize(size int) (float32, float32) {

	rooms := mp.GetRoomsBySize(size)

	if len(rooms) == 0 {
		return 0, 0
	}

	room := rooms[rand.Intn(len(rooms))]

	centerX := room.X*helpers.TILE_SIZE + room.Width*helpers.TILE_SIZE/2
	centerY := room.Y*helpers.TILE_SIZE + room.Height*helpers.TILE_SIZE/2

	return float32(centerX), float32(centerY)
}

func (r *Room) ProperRoomLightning() (scale float32, radius float32) {
	switch r.Size {
	case SmallRoom:
		return 0.6, 30 + 10
	case MediumRoom:
		return 0.8, 40 + 10
	case LargeRoom:
		return 1, 50 + 10
	}
	return 0, 0

}

func (m *Map) GetRoomByRect(rect helpers.Rectangle) *Room {
	for _, room := range m.rooms {
		if room.X == rect.X && room.Y == rect.Y &&
			room.Width == rect.Width && room.Height == rect.Height {
			return room
		}
	}
	return nil
}

// Add this new method to get walkable areas that aren't rooms
func (m *Map) GetCorridorTiles() []rl.Vector2 {
	var corridorTiles []rl.Vector2

	// Check each tile in the map
	for x := 0; x < helpers.MAP_WIDTH; x++ {
		for y := 0; y < helpers.MAP_HEIGHT; y++ {
			if m.dungeon[x][y] == 1 { // If it's a walkable tile
				isInRoom := false
				// Check if this tile is in any room
				for _, room := range m.rooms {
					if room.ContainsPoint(rl.Vector2{
						X: float32(x * helpers.TILE_SIZE),
						Y: float32(y * helpers.TILE_SIZE),
					}) {
						isInRoom = true
						break
					}
				}
				// If it's not in any room, it's a corridor tile
				if !isInRoom {
					corridorTiles = append(corridorTiles, rl.Vector2{
						X: float32(x * helpers.TILE_SIZE),
						Y: float32(y * helpers.TILE_SIZE),
					})
				}
			}
		}
	}
	return corridorTiles
}
