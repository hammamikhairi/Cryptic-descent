package effects

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Particle struct {
	Position rl.Vector2
	Velocity rl.Vector2
	Color    rl.Color
	Life     float32
	MaxLife  float32
	Size     float32
	Alpha    float32
}

type ParticleSystem struct {
	Particles []Particle
	Active    bool
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		Particles: make([]Particle, 0),
		Active:    true,
	}
}

func (ps *ParticleSystem) EmitParticles(position rl.Vector2, count int, color rl.Color, particleType string) {
	for i := 0; i < count; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float32()*1.5 + 0.5

		var life float32
		var size float32

		switch particleType {
		case "hit":
			life = 0.3
			size = rand.Float32()*1.0 + 0.5
			color = rl.Color{
				R: 139,
				G: uint8(rand.Intn(20)),
				B: uint8(rand.Intn(20)),
				A: 255,
			}
		case "death":
			life = 0.6
			size = rand.Float32()*1.5 + 1.0
			color = rl.Color{
				R: 120,
				G: uint8(rand.Intn(10)),
				B: uint8(rand.Intn(10)),
				A: 255,
			}
		case "heart":
			life = 0.3
			size = rand.Float32()*1.0 + 0.8
			color = rl.Color{
				R: 139,
				G: uint8(rand.Intn(20)),
				B: uint8(rand.Intn(20)),
				A: 255,
			}
		default:
			life = 0.3
			size = rand.Float32()*1.0 + 0.5
		}

		particle := Particle{
			Position: position,
			Velocity: rl.Vector2{
				X: float32(math.Cos(angle)) * speed,
				Y: float32(math.Sin(angle)) * speed,
			},
			Color:   color,
			Life:    life,
			MaxLife: life,
			Size:    size,
			Alpha:   1.0,
		}
		ps.Particles = append(ps.Particles, particle)
	}
}

func (ps *ParticleSystem) Update(dt float32) {
	activeParticles := make([]Particle, 0, len(ps.Particles))

	for _, p := range ps.Particles {
		p.Life -= dt
		if p.Life > 0 {
			p.Velocity.Y += dt * 2.0

			p.Position.X += p.Velocity.X
			p.Position.Y += p.Velocity.Y
			p.Alpha = (p.Life / p.MaxLife) * 0.8
			activeParticles = append(activeParticles, p)
		}
	}

	if len(activeParticles) > 1000 {
		activeParticles = activeParticles[len(activeParticles)-1000:]
	}

	ps.Particles = activeParticles
}

func (ps *ParticleSystem) Draw() {
	for _, p := range ps.Particles {
		color := p.Color
		color.A = uint8(255 * p.Alpha)
		rl.DrawCircleV(p.Position, p.Size, color)
	}
}
