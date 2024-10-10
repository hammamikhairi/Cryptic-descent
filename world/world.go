package world

import (
	helpers "crydes/helpers"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// World represents the game world
type World struct {
	Map        *Map
	Props      []*Prop
	Pathfinder *Pathfinder
}

// NewWorld creates a new world instance
func NewWorld() *World {
	mp := NewMap()
	return &World{
		Map:        mp,
		Pathfinder: NewPathfinder(mp),
	}
}

func (w *World) PlayerSpawn() (float32, float32) {
	return w.Map.FirstRoomPosition()
}

func (w *World) SwitchMap() (float32, float32) {
	// Future: Update world elements
	return w.Map.SwitchMap()
}

func (w *World) Update(deltaTime float32) {
	for _, p := range w.Props {
		p.Update(deltaTime)
	}
}

// Render draws the world elements on the screen
func (w *World) Render() {
	w.Map.Render()
	for _, p := range w.Props {
		p.Render()
	}
	// w.Pathfinder.Render(
	// 	currentRoomIndex,
	// )
}

func (w *World) NewProp(id int, x, y float32, scale float32, size rl.Vector2, animations *helpers.Animation, isAnimated bool) {
	w.Props = append(w.Props, NewProp(id, x, y, scale, size, animations, isAnimated))
}
