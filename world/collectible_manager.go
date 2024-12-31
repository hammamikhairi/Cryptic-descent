package world

import (
	"crydes/helpers"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollectibleManager struct {
	items       map[int]*CollectibleItem
	effectsChan chan ItemEffectEvent
	playerPos   *rl.Vector2
}

func NewCollectibleManager() *CollectibleManager {
	return &CollectibleManager{
		items:       make(map[int]*CollectibleItem),
		effectsChan: make(chan ItemEffectEvent, 10), // Buffered channel
		playerPos:   nil,
	}
}

func (cm *CollectibleManager) SetPlayerPosition(plPos *rl.Vector2) {
	cm.playerPos = plPos
}

func (cm *CollectibleManager) AddItem(id int, itemType ItemType, x, y float32) {
	animation := LoadItemAnimation(itemType)
	item := NewCollectibleItem(id, itemType, x, y, animation, cm.effectsChan)
	cm.items[id] = item
}

func (cm *CollectibleManager) Update(refreshRate float32) {
	for _, item := range cm.items {
		item.SetPlayerPosition(cm.playerPos)
		item.Update(refreshRate)
		// fmt.Println(cm.playerPos, item.Position)
	}
}

func (cm *CollectibleManager) Render() {
	for _, item := range cm.items {
		item.Render()
	}
}

func (cm *CollectibleManager) UpdatePlayerPosition(pos *rl.Vector2) {
	cm.playerPos = pos
}

func (cm *CollectibleManager) GetEffectsChan() chan ItemEffectEvent {
	return cm.effectsChan
}

func (cm *CollectibleManager) ScatterCollectibles(rooms []helpers.Rectangle, mp *Map) {
	// Clear existing items
	cm.items = make(map[int]*CollectibleItem)

	var itemID int = 1

	for i, room := range rooms {
		// Skip the first room (starting room)
		if i == 0 {
			continue
		}

		// Convert room to proper Room type to get size
		actualRoom := mp.GetRoomByRect(room)
		if actualRoom == nil {
			continue
		}

		// Calculate number of items based on room size
		numItems := calculateItemsForRoom(actualRoom.Size)

		for j := 0; j < numItems; j++ {
			pos := room.GetRandomPosInRect()
			itemType := getRandomItemType()

			cm.AddItem(itemID, itemType, pos.X, pos.Y)
			itemID++
		}
	}
}

func calculateItemsForRoom(size RoomSize) int {
	switch size {
	case SmallRoom:
		return 1 + rand.Intn(2) // 0-1 items
	case MediumRoom:
		return 1 + rand.Intn(2) // 1-2 items
	case LargeRoom:
		return 2 + rand.Intn(2) // 2-3 items
	default:
		return 1
	}
}

func getRandomItemType() ItemType {
	types := []ItemType{
		HealthPotion,
		SpeedPotion,
		Poison,
		// Add more item types here as needed
	}
	return types[rand.Intn(len(types))]
}
