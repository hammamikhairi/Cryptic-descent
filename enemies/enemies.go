package enemies

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"crydes/helpers"
	"crydes/player"
	"crydes/world"
)

type EnemiesManager struct {
	Enemies    []*Enemy
	Animations map[string]*map[string]*helpers.Animation

	Map *world.Map

	inComingDamage chan rl.Rectangle
}

func NewEnemiesManager(pX, pY float32, mp *world.Map, playerAttackChan chan rl.Rectangle) *EnemiesManager {
	manager := &EnemiesManager{
		Enemies:        []*Enemy{},
		Animations:     map[string]*map[string]*helpers.Animation{},
		Map:            mp,
		inComingDamage: playerAttackChan,
	}

	manager.LoadAnimations()

	manager.Enemies = append(
		manager.Enemies,
		NewEnemy(1, pX+20, pY+20, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["spider"], 3),
		NewEnemy(1, pX-20, pY-20, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["skeleton"], 3),
		NewEnemy(1, pX-20, pY+20, 0.5, rl.NewVector2(16, 16), 200, manager.Animations["skeleton"], 3),
	)

	go manager.ListenForDamage()

	return manager
}

var lastSpawnTime time.Time

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
	if time.Since(lastSpawnTime) > time.Second*5 {
		lastSpawnTime = time.Now()
		em.Enemies = append(
			em.Enemies,
			NewEnemy(1, p.Position.X+float32(rand.Intn(100)), p.Position.Y+float32(rand.Intn(100)), 0.5, rl.NewVector2(16, 16), 200, em.Animations["goblin"], 3))
	}

}

func (em *EnemiesManager) Render() {
	for _, e := range em.Enemies {
		e.Render()
	}
}

func (em *EnemiesManager) PlayerAttack(area rl.Rectangle) {
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
