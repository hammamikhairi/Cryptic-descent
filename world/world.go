package world

// World represents the game world
type World struct {
	Map *Map
}

// NewWorld creates a new world instance
func NewWorld() *World {
	return &World{
		Map: NewMap(),
	}
}

func (w *World) PlayerSpawn() (float32, float32) {
	return w.Map.FirstRoomPosition()
}

func (w *World) SwitchMap() (float32, float32) {
	// Future: Update world elements
	return w.Map.SwitchMap()
}

// Render draws the world elements on the screen
func (w *World) Render() {
	w.Map.Render()
}
