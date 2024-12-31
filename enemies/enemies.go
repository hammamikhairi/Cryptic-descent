package enemies

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"

	"crydes/audio"
	"crydes/helpers"
	"crydes/player"
	"crydes/world"
)

type EnemiesManager struct {
	Enemies    []*Enemy
	Animations map[string]*map[string]*helpers.Animation

	Map   *world.Map
	Rooms []helpers.Rectangle

	inComingDamage chan rl.Rectangle
	soundManager   *audio.SoundManager
}

func NewEnemiesManager(pX, pY float32, mp *world.Map, playerAttackChan chan rl.Rectangle, rooms []helpers.Rectangle, soundManager *audio.SoundManager) *EnemiesManager {
	manager := &EnemiesManager{
		Enemies:        []*Enemy{},
		Animations:     map[string]*map[string]*helpers.Animation{},
		Map:            mp,
		inComingDamage: playerAttackChan,
		Rooms:          rooms,
		soundManager:   soundManager,
	}

	manager.LoadAnimations()

	// manager.Enemies = append(
	// 	manager.Enemies,
	// 	// NewEnemy(1, pX+20, pY+20, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["spider"], 3),
	// 	// NewEnemy(1, pX-20, pY-20, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["skeleton"], 3),
	// 	// NewEnemy(1, pX-20, pY+20, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["skeleton"], 3),
	// )

	go manager.ListenForDamage()

	return manager
}

// var lastSpawnTime time.Time

func (em *EnemiesManager) Update(refreshRate float32, p *player.Player) {
	for _, e := range em.Enemies {

		if e.isDead {
			continue
		}

		e.Update(refreshRate, p)
	}

	// clean up dead enemies
	// for i, e := range em.Enemies {
	// 	if e.isDead {
	// 		em.Enemies = append(em.Enemies[:i], em.Enemies[i+1:]...)
	// 	}
	// }

	// spawn enemies every 5 seconds
	// if time.Since(lastSpawnTime) > time.Second*5 {
	// 	lastSpawnTime = time.Now()
	// 	em.Enemies = append(
	// 		em.Enemies,
	// 		NewEnemy(1, p.Position.X+float32(rand.Intn(100)), p.Position.Y+float32(rand.Intn(100)), 0.5, rl.NewVector2(16, 16), 200, em.Animations["goblin"], 3))
	// }
}

func (em *EnemiesManager) ResetEnemies() {
	em.Enemies = []*Enemy{}
	em.SpawnEnemies()
}

func (em *EnemiesManager) SpawnEnemies() {
	// First spawn in rooms as before
	em.spawnEnemiesInRooms()
	// Then spawn in corridors
	// em.spawnEnemiesInCorridors()
}

// Move the existing room spawning logic to this method
func (em *EnemiesManager) spawnEnemiesInRooms() {
	for i, room := range em.Rooms {
		if i == 0 {
			continue
		} // Skip starting room

		actualRoom := em.Map.GetRoomByRect(room)
		if actualRoom == nil {
			continue
		}

		numEnemies := calculateEnemiesForRoom(actualRoom.Size)

		for j := 0; j < numEnemies; j++ {
			ePos := room.GetRandomPosInRect()
			eType := helpers.GetRandomEnemyType()
			scale, speed, health := getEnemyAttributes(actualRoom.Size)

			enemy := NewEnemy(
				j,
				ePos.X,
				ePos.Y,
				scale,
				rl.NewVector2(16, 16),
				speed,
				em.Animations[eType],
				health,
				i,
				em.soundManager,
			)

			em.Enemies = append(em.Enemies, enemy)
		}
	}
}

// Add new method for corridor spawning
func (em *EnemiesManager) spawnEnemiesInCorridors() {
	corridorTiles := em.Map.GetCorridorTiles()

	// Spawn an enemy every N tiles in corridors (adjust as needed)
	spawnFrequency := 20 // Adjust this value to control density

	for i := 0; i < len(corridorTiles); i += spawnFrequency {
		if rand.Float32() < 0.3 { // 30% chance to spawn at each valid location
			pos := corridorTiles[i]
			eType := helpers.GetRandomEnemyType()

			// Corridor enemies are slightly weaker
			enemy := NewEnemy(
				i,
				pos.X,
				pos.Y,
				0.6, // Smaller scale
				rl.NewVector2(16, 16),
				150, // Slower speed
				em.Animations[eType],
				2,  // Less health
				-1, // No specific room
				em.soundManager,
			)

			em.Enemies = append(em.Enemies, enemy)
		}
	}
}

func calculateEnemiesForRoom(size world.RoomSize) int {
	switch size {
	case world.SmallRoom:
		return 1 + rand.Intn(2) // 1-2 enemies
	case world.MediumRoom:
		return 2 + rand.Intn(6) // 2-4 enemies
	case world.LargeRoom:
		return 4 + rand.Intn(10) // 4-7 enemies
	default:
		return 2
	}
}

func getEnemyAttributes(size world.RoomSize) (scale, speed float32, health int) {
	switch size {
	case world.SmallRoom:
		return 0.6, 150, 2
	case world.MediumRoom:
		return 0.8, 200, 3
	case world.LargeRoom:
		return 1.1, 250, 4
	default:
		return 0.5, 200, 3
	}
}

func (em *EnemiesManager) AddEnemy(e *Enemy) {
	em.Enemies = append(em.Enemies, e)
}

func (em *EnemiesManager) Render() {
	for _, e := range em.Enemies {
		e.Render()
	}
}

func (em *EnemiesManager) PlayerAttack(area rl.Rectangle) {
	em.soundManager.RequestSound("sword_swing", 1.0, 1.0)
	for _, e := range em.Enemies {
		e.DamageChan <- area
	}
}

func (em *EnemiesManager) ListenForDamage() {
	for {
		select {
		case damageArea := <-em.inComingDamage:
			em.PlayerAttack(damageArea)
		}
	}
}
