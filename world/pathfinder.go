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

	// Convert path to smooth points
	smoothPath := pf.createSmoothPath()

	// Draw the smooth path
	for i := 0; i < len(smoothPath)-1; i++ {
		rl.DrawLineEx(smoothPath[i], smoothPath[i+1], 3, rl.ColorAlpha(rl.Green, 0.7))
	}

	// Optionally, draw some decorative elements
	pf.drawPathDecorations(smoothPath)
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
