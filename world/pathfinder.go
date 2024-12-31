package world

import (
	"crydes/helpers"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Define the Pathfinder struct which will handle pathfinding
type Pathfinder struct {
	grid   [helpers.MAP_WIDTH][helpers.MAP_HEIGHT]int // Copy of the dungeon map
	open   []Node                                     // List of nodes to be evaluated
	closed []Node                                     // List of nodes already evaluated
	path   []Node                                     // The resulting path
	rooms  []*Room                                    // List of rooms in the map
}

// Define the Node struct which represents a point in the grid
type Node struct {
	x, y   int     // Coordinates of the node
	costG  float64 // Cost from start to current node
	costH  float64 // Heuristic cost from current node to end
	parent *Node   // Pointer to the parent node (for path reconstruction)
}

// Heuristic function to estimate the cost from the current node to the target node
func heuristic(x1, y1, x2, y2 int) float64 {
	return math.Abs(float64(x2-x1)) + math.Abs(float64(y2-y1))
}

// Returns the total cost (G + H) for the node
func (n Node) TotalCost() float64 {
	return n.costG + n.costH
}

// Initialize Pathfinder with the map grid
func NewPathfinder(m *Map) *Pathfinder {
	return &Pathfinder{
		grid:  m.dungeon,
		rooms: m.rooms,
	}
}

// Update method runs the A* algorithm to find a path from (x1, y1) to (x2, y2)
func (pf *Pathfinder) Update(x1, y1, x2, y2 int) {

	// Reset pathfinding data
	pf.open = make([]Node, 0)
	pf.closed = make([]Node, 0)
	pf.path = make([]Node, 0)

	startNode := Node{x: x1, y: y1, costG: 0, costH: heuristic(x1, y1, x2, y2)}
	pf.open = append(pf.open, startNode)

	for len(pf.open) > 0 {
		// Find node in open list with the lowest cost
		currentIndex := pf.getLowestCostIndex()
		currentNode := pf.open[currentIndex]

		// Check if we have reached the goal
		if currentNode.x == x2 && currentNode.y == y2 {
			pf.path = pf.reconstructPath(currentNode)
			return
		}

		// Remove current node from open and add it to closed
		pf.open = append(pf.open[:currentIndex], pf.open[currentIndex+1:]...)
		pf.closed = append(pf.closed, currentNode)

		// Explore neighbors
		for _, neighbor := range pf.getNeighbors(currentNode) {
			if pf.isInClosedList(neighbor) || pf.grid[neighbor.x][neighbor.y] == 0 {
				continue // Skip if neighbor is in closed list or not walkable
			}

			// Calculate tentative G cost
			tentativeG := currentNode.costG + 1

			// If the neighbor is not in the open list, add it with its G cost
			if !pf.isInOpenList(neighbor) {
				neighbor.costG = tentativeG
				neighbor.costH = heuristic(neighbor.x, neighbor.y, x2, y2)
				neighbor.parent = &currentNode
				pf.open = append(pf.open, neighbor)
			} else if tentativeG < neighbor.costG {
				// If we found a shorter path to the neighbor, update its cost and parent
				neighbor.costG = tentativeG
				neighbor.parent = &currentNode
			}
		}
	}
}

func (pf *Pathfinder) Render(currentRoomIndex int) {
	if currentRoomIndex == -1 || currentRoomIndex >= len(pf.rooms) {
		// No room found or invalid index, fallback to drawing the path normally.
		pf.drawSmoothPath()
		return
	}

	// Get the current room
	// currentRoom := pf.rooms[currentRoomIndex]

	// Check if the room has both entrance and exit
	// if len(currentRoom.Doors) >= 2 {
	// 	// Draw a direct line from entrance to exit
	// 	entrance := currentRoom.Doors[0]
	// 	exit := currentRoom.Doors[1]
	// 	rl.DrawLineEx(entrance, exit, 3, rl.ColorAlpha(rl.Green, 0.7))
	// } else {
	// Draw the smooth path normally if there is no entrance-exit configuration
	pf.drawSmoothPath()
	// }
}

// Helper function to draw the smooth path if no entrance/exit is defined.
func (pf *Pathfinder) drawSmoothPath() {
	if len(pf.path) < 2 {
		return
	}

	// Convert path points to world coordinates
	points := make([]rl.Vector2, len(pf.path))
	for i, node := range pf.path {
		points[i] = rl.Vector2{
			X: float32(node.x*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
			Y: float32(node.y*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
		}
	}

	// Apply Catmull-Rom spline interpolation
	smoothPoints := make([]rl.Vector2, 0)
	segmentPoints := 8 // Number of points per segment

	// Add first point
	smoothPoints = append(smoothPoints, points[0])

	// Interpolate middle points
	for i := 0; i < len(points)-3; i++ {
		p0 := points[i]
		p1 := points[i+1]
		p2 := points[i+2]
		p3 := points[i+3]

		for j := 0; j < segmentPoints; j++ {
			t := float32(j) / float32(segmentPoints)

			// Catmull-Rom interpolation
			point := rl.Vector2{
				X: catmullRom(t, p0.X, p1.X, p2.X, p3.X),
				Y: catmullRom(t, p0.Y, p1.Y, p2.Y, p3.Y),
			}

			smoothPoints = append(smoothPoints, point)
		}
	}

	// Add last point
	smoothPoints = append(smoothPoints, points[len(points)-1])

	// Draw the smooth path with gradient effect
	for i := 0; i < len(smoothPoints)-1; i++ {
		progress := float32(i) / float32(len(smoothPoints)-1)

		// Create a gradient effect from green to blue
		color := rl.ColorAlpha(rl.Color{
			R: uint8(0),
			G: uint8(255 * (1 - progress)),
			B: uint8(255 * progress),
			A: 255,
		}, 0.6)

		// Draw thicker line with smooth fade
		thickness := 4.0 - (progress * 2.0)
		rl.DrawLineEx(smoothPoints[i], smoothPoints[i+1], thickness, color)

		// Add glow effect
		rl.DrawLineEx(smoothPoints[i], smoothPoints[i+1], thickness+2,
			rl.ColorAlpha(color, 0.2))
	}

	// Draw direction indicators
	for i := 0; i < len(smoothPoints); i += 8 {
		if i+1 < len(smoothPoints) {
			drawDirectionIndicator(smoothPoints[i], smoothPoints[i+1])
		}
	}
}

func catmullRom(t, p0, p1, p2, p3 float32) float32 {
	t2 := t * t
	t3 := t2 * t

	return 0.5 * ((2 * p1) +
		(-p0+p2)*t +
		(2*p0-5*p1+4*p2-p3)*t2 +
		(-p0+3*p1-3*p2+p3)*t3)
}

func drawDirectionIndicator(p1, p2 rl.Vector2) {
	dir := rl.Vector2Subtract(p2, p1)
	dir = rl.Vector2Normalize(dir)

	// Calculate perpendicular vector for arrow head
	perp := rl.Vector2{X: -dir.Y, Y: dir.X}
	arrowSize := float32(6.0)

	// Calculate arrow head points
	tip := rl.Vector2Add(p1, rl.Vector2Scale(dir, arrowSize*2))
	left := rl.Vector2Subtract(tip, rl.Vector2Scale(rl.Vector2Add(dir, perp), arrowSize))
	right := rl.Vector2Subtract(tip, rl.Vector2Scale(rl.Vector2Subtract(dir, perp), arrowSize))

	// Draw arrow head
	rl.DrawTriangle(tip, left, right, rl.ColorAlpha(rl.Green, 0.4))
}

// Find the index of the current room based on player position
func (pf *Pathfinder) findCurrentRoom(position rl.Vector2) int {
	for i, room := range pf.rooms {
		if room.ContainsPoint(position) {
			return i
		}
	}
	return -1
}

func (pf *Pathfinder) createSmoothPath() []rl.Vector2 {
	smoothPath := make([]rl.Vector2, 0)

	for i, node := range pf.path {
		baseX := float32(node.x*helpers.TILE_SIZE + helpers.TILE_SIZE/2)
		baseY := float32(node.y*helpers.TILE_SIZE + helpers.TILE_SIZE/2)

		// Add some randomness to the point position
		jitterX := (rand.Float32() - 0.5) * float32(helpers.TILE_SIZE) * 0.5
		jitterY := (rand.Float32() - 0.5) * float32(helpers.TILE_SIZE) * 0.5

		point := rl.Vector2{
			X: baseX + jitterX,
			Y: baseY + jitterY,
		}

		// Add intermediate points for longer segments
		if i > 0 {
			prevPoint := smoothPath[len(smoothPath)-1]
			distance := rl.Vector2Distance(prevPoint, point)

			if distance > float32(helpers.TILE_SIZE)*1.5 {
				numIntermediates := int(distance / float32(helpers.TILE_SIZE))
				for j := 1; j < numIntermediates; j++ {
					t := float32(j) / float32(numIntermediates)
					intermediatePoint := rl.Vector2Lerp(prevPoint, point, t)

					// Add some randomness to intermediate points
					intermediatePoint.X += (rand.Float32() - 0.5) * float32(helpers.TILE_SIZE) * 0.3
					intermediatePoint.Y += (rand.Float32() - 0.5) * float32(helpers.TILE_SIZE) * 0.3

					smoothPath = append(smoothPath, intermediatePoint)
				}
			}
		}

		smoothPath = append(smoothPath, point)
	}

	return smoothPath
}

func (pf *Pathfinder) drawPathDecorations(smoothPath []rl.Vector2) {
	for i := 0; i < len(smoothPath); i++ {
		// Draw small circles at each point
		rl.DrawCircleV(smoothPath[i], 2, rl.ColorAlpha(rl.Green, 0.5))

		// Occasionally draw some "foliage" or other decorative elements
		if rand.Float32() < 0.2 {
			pf.drawFoliage(smoothPath[i])
		}
	}
}

func (pf *Pathfinder) drawFoliage(position rl.Vector2) {
	numLeaves := rand.Intn(3) + 2
	for i := 0; i < numLeaves; i++ {
		angle := rand.Float32() * 2 * math.Pi
		distance := rand.Float32()*5 + 2
		leafPos := rl.Vector2Add(position, rl.Vector2{
			X: float32(math.Cos(float64(angle))) * distance,
			Y: float32(math.Sin(float64(angle))) * distance,
		})
		rl.DrawCircleV(leafPos, 1, rl.ColorAlpha(rl.Green, 0.3))
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Get the index of the node with the lowest total cost (G + H) in the open list
func (pf *Pathfinder) getLowestCostIndex() int {
	lowestIndex := 0
	lowestCost := pf.open[0].TotalCost()
	for i, node := range pf.open {
		if node.TotalCost() < lowestCost {
			lowestIndex = i
			lowestCost = node.TotalCost()
		}
	}
	return lowestIndex
}

// Reconstruct the path by backtracking from the goal node
func (pf *Pathfinder) reconstructPath(goalNode Node) []Node {
	var path []Node
	for node := &goalNode; node != nil; node = node.parent {
		path = append([]Node{*node}, path...)
	}
	return path
}

// Get valid neighbors of the current node (up, down, left, right)
func (pf *Pathfinder) getNeighbors(node Node) []Node {
	var neighbors []Node
	directions := [][2]int{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0}, // Up, Down, Left, Right
	}
	for _, dir := range directions {
		x, y := node.x+dir[0], node.y+dir[1]
		if x >= 0 && x < helpers.MAP_WIDTH && y >= 0 && y < helpers.MAP_HEIGHT {
			neighbors = append(neighbors, Node{x: x, y: y})
		}
	}
	return neighbors
}

// Check if a node is in the closed list
func (pf *Pathfinder) isInClosedList(node Node) bool {
	for _, n := range pf.closed {
		if n.x == node.x && n.y == node.y {
			return true
		}
	}
	return false
}

// Check if a node is in the open list
func (pf *Pathfinder) isInOpenList(node Node) bool {
	for _, n := range pf.open {
		if n.x == node.x && n.y == node.y {
			return true
		}
	}
	return false
}
