package world

import (
	"crydes/helpers"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PropsManager struct {
	rooms *[]*Room
	props []*Prop
	Map   *Map
}

// Prop represents an interactive or static item in the game.
type Prop struct {
	ID          int                // Unique identifier for the prop
	Type        string             // Type of the prop (e.g., "chest", "door", "key")
	Position    rl.Vector2         // Position of the prop
	Size        rl.Vector2         // Size of the prop for collision detection
	Scale       float32            // Scale for rendering
	Visible     bool               // Visibility of the prop
	CurrentAnim *helpers.Animation // Current animation to display if animated
	Animation   *helpers.Animation // Map of animations
	IsAnimated  bool               // Flag indicating if the prop is animated
	Color       rl.Color           // Base color for the prop (e.g., for shading effects)
	LTRadius    float32            // Light radius for light sources

	// Optional properties to control prop behavior.
	Rotation float32 // Rotation in degrees
	Opacity  float32 // Opacity for rendering (0.0 to 1.0)
	Friction float32 // Friction to apply when interacting with other objects
}

func newPropsManager(rooms *[]*Room, mp *Map) *PropsManager {
	return &PropsManager{
		rooms: rooms,
		props: []*Prop{},
		Map:   mp,
	}
}

func (pm *PropsManager) SetUpProps() {
	// First set up room props (lights etc)
	pm.setupRoomProps()
	// Then set up corridor props
	pm.setupCorridorProps()
}

func (pm *PropsManager) setupRoomProps() {
	const minDistance = 100.0 // Minimum distance between props

	for _, room := range *pm.rooms {
		lightPos := room.GetLightPositions()
		scale, radius := room.ProperRoomLightning()

		// Filter positions that are too close to existing props
		var validPositions []rl.Vector2
		for _, pos := range lightPos {
			if pm.isPositionValid(pos.X, pos.Y, minDistance) {
				validPositions = append(validPositions, pos)
			}
		}

		// Place props at valid positions
		for _, pos := range validPositions {
			pm.props = append(pm.props, NewProp(
				1,
				"fire",
				pos.X,
				pos.Y,
				scale,
				radius,
				rl.NewVector2(16, 16),
				helpers.LoadAnimation(
					"assets/fireplace/1.png",
					"assets/fireplace/2.png",
					"assets/fireplace/3.png",
					"assets/fireplace/4.png",
				),
				true,
			))
		}
	}
}

func (pm *PropsManager) setupCorridorProps() {
	corridorTiles := pm.Map.GetCorridorTiles()
	const (
		propFrequency = 25    // Adjust this value to control density
		minDistance   = 100.0 // Minimum distance between props
	)

	// Shuffle corridor tiles for random placement
	rand.Shuffle(len(corridorTiles), func(i, j int) {
		corridorTiles[i], corridorTiles[j] = corridorTiles[j], corridorTiles[i]
	})

	for _, pos := range corridorTiles {
		if rand.Float32() < 0.4 { // 40% chance to try spawning
			if pm.isPositionValid(pos.X, pos.Y, minDistance) {
				// Create a smaller light source for corridors
				pm.props = append(pm.props, NewProp(
					1,
					"fire",
					pos.X,
					pos.Y,
					0.5, // Smaller scale
					20,  // Smaller light radius
					rl.NewVector2(16, 16),
					helpers.LoadAnimation(
						"assets/fireplace/1.png",
						"assets/fireplace/2.png",
						"assets/fireplace/3.png",
						"assets/fireplace/4.png",
					),
					true,
				))
			}
		}
	}
}

func (pm *PropsManager) Update(refreshRate float32) {
	for _, prop := range pm.props {
		prop.Update(refreshRate)
	}
}

func (pm *PropsManager) Render() {
	for _, prop := range pm.props {
		prop.Render()
	}
}

func (pm *PropsManager) GetProps() *[]*Prop {
	return &pm.props
}

// NewProp initializes and returns a new Prop instance.
func NewProp(id int, tp string, x, y float32, scale, radius float32, size rl.Vector2, animations *helpers.Animation, isAnimated bool) *Prop {
	return &Prop{
		ID:          id,
		Type:        tp,
		Position:    rl.NewVector2(x, y),
		Size:        size,
		Scale:       scale,
		LTRadius:    radius,
		Visible:     true,
		IsAnimated:  isAnimated,
		Animation:   animations,
		CurrentAnim: animations, // Set a default animation if available
		Color:       rl.White,   // Default color
		Rotation:    0,
		Opacity:     1.0,
		Friction:    1.0,
	}
}

// Update handles animation and other dynamic properties.
func (p *Prop) Update(refreshRate float32) {
	if !p.Visible {
		return
	}

	if p.IsAnimated && p.CurrentAnim != nil {
		// helpers.DEBUG("====Updating animation for prop %d", p.ID)
		p.UpdateAnimation(refreshRate)
	}
}

// UpdateAnimation updates the current animation frame of the prop.
func (p *Prop) UpdateAnimation(refreshRate float32) {
	if p.CurrentAnim == nil {
		return
	}

	p.CurrentAnim.Timer += refreshRate
	if p.CurrentAnim.Timer >= p.CurrentAnim.FrameTime {
		p.CurrentAnim.CurrentFrame = (p.CurrentAnim.CurrentFrame + 1) % len(p.CurrentAnim.Frames)
		p.CurrentAnim.Timer = 0
	}
}

// Render draws the prop based on its properties and state.
func (p *Prop) Render() {
	// if !p.Visible {
	// 	return
	// }

	// Determine color with opacity
	finalColor := rl.Fade(p.Color, p.Opacity)

	if p.IsAnimated && p.CurrentAnim != nil {
		rl.DrawTextureEx(p.CurrentAnim.Frames[p.CurrentAnim.CurrentFrame], p.Position, p.Rotation, p.Scale, finalColor)
	} else {
		rl.DrawTextureEx(p.Animation.Frames[0], p.Position, p.Rotation, p.Scale, finalColor)
	}
}

// SetPosition updates the position of the prop.
func (p *Prop) SetPosition(x, y float32) {
	p.Position = rl.NewVector2(x, y)
}

// SetVisibility toggles the visibility of the prop.
func (p *Prop) SetVisibility(visible bool) {
	p.Visible = visible
}

// RandomizePosition places the prop at a random position within the given bounds.
func (p *Prop) RandomizePosition(bounds rl.Rectangle) {
	x := bounds.X + rand.Float32()*(bounds.Width-p.Size.X*p.Scale)
	y := bounds.Y + rand.Float32()*(bounds.Height-p.Size.Y*p.Scale)
	p.Position = rl.NewVector2(x, y)
}

// CheckCollision checks if the prop collides with another rectangle.
// func (p *Prop) CheckCollision(target rl.Rectangle) bool {
// 	propRect := rl.NewRectangle(p.Position.X, p.Position.Y, p.Size.X*p.Scale, p.Size.Y*p.Scale)
// 	return rl.CheckCollisionRecs(propRect, target)
// }

// ApplyFriction reduces the velocity of the prop based on its friction coefficient.
func (p *Prop) ApplyFriction(velocity rl.Vector2) rl.Vector2 {
	velocity.X *= p.Friction
	velocity.Y *= p.Friction
	return velocity
}

// Add this helper function to check distance between two points
func distanceBetweenPoints(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

// Add this method to PropsManager
func (pm *PropsManager) isPositionValid(x, y float32, minDistance float32) bool {
	for _, prop := range pm.props {
		dist := distanceBetweenPoints(x, y, prop.Position.X, prop.Position.Y)
		if dist < minDistance {
			return false
		}
	}
	return true
}
