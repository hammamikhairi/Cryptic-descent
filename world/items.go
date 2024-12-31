package world

import (
	"crydes/helpers"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ItemType represents different types of collectible items
type ItemType string

const (
	HealthPotion ItemType = "health_potion"
	SpeedPotion  ItemType = "speed_potion"
	Key          ItemType = "key"
	Coin         ItemType = "coin"
	Poison       ItemType = "poison"
)

// ItemEffect represents the effect an item has when collected
type ItemEffect struct {
	Type     string        // Type of effect (e.g., "heal", "boost_speed")
	Value    float32       // Magnitude of the effect
	Duration time.Duration // Duration of the effect in seconds (0 for instant effects)
}

// ItemEffectEvent represents the event triggered by collecting an item
type ItemEffectEvent struct {
	SourceID int         // ID of the item that caused the effect
	Effect   *ItemEffect // The effect to apply
}

// CollectibleItem extends the base Prop type with item-specific properties
type CollectibleItem struct {
	*Prop                              // Embed the base Prop type
	ItemType    ItemType               // Type of item
	Effect      *ItemEffect            // Effect when collected
	Collected   bool                   // Whether the item has been collected
	HoverOffset float32                // Offset for hover animation
	HoverSpeed  float32                // Speed of hover animation
	Time        float32                // Time tracker for animations
	EffectsChan chan<- ItemEffectEvent // Add this field
	playerPos   *rl.Vector2
}

// NewCollectibleItem creates a new collectible item
func NewCollectibleItem(id int, itemType ItemType, x, y float32, animation *helpers.Animation, effectsChan chan<- ItemEffectEvent) *CollectibleItem {
	var effect *ItemEffect
	scale := float32(1.0)
	size := rl.NewVector2(16, 16)

	// Configure effect based on item type
	switch itemType {
	case HealthPotion:
		effect = &ItemEffect{Type: "heal", Value: 2, Duration: 0}
	case SpeedPotion:
		effect = &ItemEffect{Type: "speed", Value: 2, Duration: time.Second * 30}
	case Poison:
		effect = &ItemEffect{Type: "poison", Value: 2, Duration: time.Second * 1}
	case Key:
		effect = &ItemEffect{Type: "key", Value: 1, Duration: 0}
	case Coin:
		effect = &ItemEffect{Type: "coin", Value: 1, Duration: 0}
	}

	baseProp := NewProp(
		id,
		string(itemType),
		x, y,
		scale,
		0, // No light radius for items
		size,
		animation,
		true,
	)

	return &CollectibleItem{
		Prop:        baseProp,
		ItemType:    itemType,
		Effect:      effect,
		HoverOffset: 0,
		HoverSpeed:  4.0,
		Collected:   false,
		EffectsChan: effectsChan,
		playerPos:   &rl.Vector2{},
	}
}

// Update handles item animations and effects
func (ci *CollectibleItem) Update(refreshRate float32) {
	if ci.Collected {
		return
	}

	// Check if player is close enough to collect
	if ci.playerPos != nil {
		collectionRadius := float32(10.0) // Adjust this value as needed
		distance := helpers.Distance(ci.Position, *ci.playerPos)

		if distance <= collectionRadius {
			ci.Collect()
			return
		}
	}

	// Update base prop animations
	ci.Prop.Update(refreshRate)

	// Update hover animation
	// ci.Time += refreshRate
	// ci.HoverOffset = float32(math.Sin(float64(ci.Time*ci.HoverSpeed))) * 5.0

	// Update position with hover offset
	// originalY := ci.Prop.Position.Y
	// ci.Prop.Position.Y = originalY + ci.HoverOffset
}

// Render draws the item with any special effects
func (ci *CollectibleItem) Render() {
	if ci.Collected {
		return
	}

	// Add a subtle glow effect
	// glowColor := rl.ColorAlpha(rl.White, 0.3)
	// glowScale := ci.Scale * 1.2
	// rl.DrawTextureEx(
	// 	ci.CurrentAnim.Frames[ci.CurrentAnim.CurrentFrame],
	// 	rl.Vector2{X: ci.Position.X - 5, Y: ci.Position.Y - 5},
	// 	ci.Rotation,
	// 	glowScale,
	// 	glowColor,
	// )

	// Draw the actual item
	ci.Prop.Render()
}

// Collect handles the collection of the item
func (ci *CollectibleItem) Collect() {
	if ci.Collected {
		return
	}
	ci.Collected = true

	// Send the effect through the channel
	if ci.Effect != nil && ci.EffectsChan != nil {
		ci.EffectsChan <- ItemEffectEvent{
			SourceID: ci.ID,
			Effect:   ci.Effect,
		}
	}
}

// LoadItemAnimation loads the appropriate animation for an item type
func LoadItemAnimation(itemType ItemType) *helpers.Animation {
	switch itemType {
	case HealthPotion:
		return helpers.LoadAnimation("health_potion",
			"assets/health_potion/1.png",
			"assets/health_potion/2.png",
			"assets/health_potion/3.png",
			"assets/health_potion/4.png",
		)
	case SpeedPotion:
		return helpers.LoadAnimation("speed_potion",
			"assets/speed_potion/9.png",
			"assets/speed_potion/10.png",
			"assets/speed_potion/11.png",
			"assets/speed_potion/12.png",
		)
	case Key:
		return helpers.LoadAnimation("key",
			"assets/items/key/1.png",
			"assets/items/key/2.png",
			"assets/items/key/3.png",
		)
	case Poison:
		return helpers.LoadAnimation("key",
			"assets/speed_potion/9.png",
			"assets/speed_potion/10.png",
			"assets/speed_potion/11.png",
			"assets/speed_potion/12.png",
		)
	case Coin:
		return helpers.LoadAnimation("coin",
			"assets/items/coin/1.png",
			"assets/items/coin/2.png",
			"assets/items/coin/3.png",
		)
	default:
		return nil
	}
}

// Add method to update player position
func (ci *CollectibleItem) SetPlayerPosition(pos *rl.Vector2) {
	if ci.playerPos != nil {
		ci.playerPos = pos
	}
}
