package world

import (
	"crydes/helpers"
	"math"
	"math/rand"

	// "math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Define the Pathfinder struct which will handle pathfinding
type Pathfinder struct {
	grid           [helpers.MAP_WIDTH][helpers.MAP_HEIGHT]int // Copy of the dungeon map
	open           []Node                                     // List of nodes to be evaluated
	closed         []Node                                     // List of nodes already evaluated
	path           []Node                                     // The resulting path
	rooms          []*Room                                    // List of rooms in the map
	currentStarPos rl.Vector2                                 // Add this new field
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
		grid:           m.dungeon,
		rooms:          m.rooms,
		currentStarPos: rl.Vector2{X: 0, Y: 0},
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
	if len(pf.path) < 2 {
		return
	}

	// Get current player position (first node in path)
	playerPos := rl.Vector2{
		X: float32(pf.path[0].x*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
		Y: float32(pf.path[0].y*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
	}

	// Get next waypoint (3-4 tiles ahead)
	nextNode := pf.path[0]
	for i := 1; i < len(pf.path) && i < 5; i++ { // Increased from 4 to 5 for further look-ahead
		nextNode = pf.path[i]
	}

	nextPos := rl.Vector2{
		X: float32(nextNode.x*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
		Y: float32(nextNode.y*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
	}

	// Calculate target position for the star (farther from player)
	direction := rl.Vector2Subtract(nextPos, playerPos)
	direction = rl.Vector2Normalize(direction)
	targetStarPos := rl.Vector2Add(playerPos, rl.Vector2Scale(direction, float32(helpers.TILE_SIZE)*2.5)) // Increased from 1.5 to 2.5

	// Smoothly interpolate current star position towards target position
	smoothFactor := float32(0.1) // Lower value = smoother movement
	pf.currentStarPos = rl.Vector2{
		X: pf.currentStarPos.X + (targetStarPos.X-pf.currentStarPos.X)*smoothFactor,
		Y: pf.currentStarPos.Y + (targetStarPos.Y-pf.currentStarPos.Y)*smoothFactor,
	}

	pf.drawStar(pf.currentStarPos)
}

// New method to get the next significant waypoint
func (pf *Pathfinder) getNextWaypoint() rl.Vector2 {
	// Get the next point in the path that's at least 3 tiles away
	currentIndex := 0
	for i := 1; i < len(pf.path); i++ {
		dist := heuristic(pf.path[currentIndex].x, pf.path[currentIndex].y,
			pf.path[i].x, pf.path[i].y)
		if dist >= 3 {
			return rl.Vector2{
				X: float32(pf.path[i].x*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
				Y: float32(pf.path[i].y*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
			}
		}
	}

	// If no point is far enough, return the last point
	lastNode := pf.path[len(pf.path)-1]
	return rl.Vector2{
		X: float32(lastNode.x*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
		Y: float32(lastNode.y*helpers.TILE_SIZE + helpers.TILE_SIZE/2),
	}
}

func (pf *Pathfinder) drawStar(position rl.Vector2) {
	time := float32(rl.GetTime())
	pulse := (float32(math.Sin(float64(time*3))) + 1) / 2

	// Small star size
	size := float32(helpers.TILE_SIZE) * 0.3
	innerRadius := size * 0.4
	outerRadius := size * (0.8 + pulse*0.2)

	points := 4 // 4-pointed star
	starColor := rl.ColorAlpha(rl.Gold, 0.8+pulse*0.2)

	for i := 0; i < points*2; i++ {
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}

		angle := float32(i) * math.Pi / float32(points)
		nextAngle := float32(i+1) * math.Pi / float32(points)

		p1 := rl.Vector2{
			X: position.X + float32(math.Cos(float64(angle)))*radius,
			Y: position.Y + float32(math.Sin(float64(angle)))*radius,
		}

		p2 := rl.Vector2{
			X: position.X + float32(math.Cos(float64(nextAngle)))*radius,
			Y: position.Y + float32(math.Sin(float64(nextAngle)))*radius,
		}

		rl.DrawLineEx(p1, p2, 2, starColor)
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
func (pf *Pathfinder) CreateSmoothPath() []rl.Vector2 {
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
