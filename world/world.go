package world

// World represents the game world
type World struct {
	Map          *Map
	PropsManager *PropsManager

	Pathfinder *Pathfinder
}

// NewWorld creates a new world instance
func NewWorld() *World {
	mp := NewMap()
	wrld := &World{
		Map:          mp,
		Pathfinder:   NewPathfinder(mp),
		PropsManager: newPropsManager(mp.GetRooms()),
	}

	wrld.PropsManager.SetUpProps()

	return wrld
}

func (w *World) PlayerSpawn() (float32, float32) {
	return w.Map.FirstRoomPosition()
}

func (w *World) SwitchMap() (float32, float32) {
	// Future: Update world elements
	x, y := w.Map.SwitchMap()

	w.Pathfinder = NewPathfinder(
		w.Map,
	)

	return x, y
}

func (w *World) Update(deltaTime float32) {
	w.PropsManager.Update(deltaTime)
}

// Render draws the world elements on the screen
func (w *World) Render() {
	w.Map.Render()
	w.PropsManager.Render()
}
