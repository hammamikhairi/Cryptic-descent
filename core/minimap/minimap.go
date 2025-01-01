package minimap

import (
	"crydes/helpers"
	"crydes/world"

	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Minimap struct {
	scale            float32
	cornerPos        rl.Vector2 // Position for corner minimap
	cornerSize       rl.Vector2 // Size for corner minimap
	centerPos        rl.Vector2 // Position for centered map view
	centerSize       rl.Vector2 // Size for centered map view
	mapData          *world.Map
	texture          rl.RenderTexture2D
	isDirty          bool
	borderPad        int32
	isFullscreen     bool // Toggle for full map view
	initialZoom      float32
	destinationX     int
	destinationY     int
	hasDestination   bool
	lastScreenWidth  float32
	lastScreenHeight float32
}

func NewMinimap(mapData *world.Map) *Minimap {
	screenWidth := float32(helpers.SCREEN_WIDTH)
	screenHeight := float32(helpers.SCREEN_HEIGHT)

	// Corner minimap dimensions (20% of screen)
	cornerSize := rl.Vector2{
		X: screenWidth * 0.2,
		Y: screenHeight * 0.2,
	}

	cornerPos := rl.Vector2{
		X: screenWidth - cornerSize.X - 20,
		Y: 20,
	}

	// Center map view dimensions (70% of screen)
	centerSize := rl.Vector2{
		X: screenWidth * 0.7,
		Y: screenHeight * 0.7,
	}

	centerPos := rl.Vector2{
		X: (screenWidth - centerSize.X) / 2,
		Y: (screenHeight - centerSize.Y) / 2,
	}

	// Calculate scale based on map size and desired minimap size
	scaleX := centerSize.X / float32(helpers.MAP_WIDTH*helpers.TILE_SIZE)
	scaleY := centerSize.Y / float32(helpers.MAP_HEIGHT*helpers.TILE_SIZE)
	scale := min(scaleX, scaleY)

	println("KJSDHKJQHDQJHDS")
	// Use the same scale for both dimensions to maintain aspect ratio
	texture := rl.LoadRenderTexture(
		int32(float32(helpers.MAP_WIDTH*helpers.TILE_SIZE)*scale),  // Match map dimensions
		int32(float32(helpers.MAP_HEIGHT*helpers.TILE_SIZE)*scale), // Match map dimensions
	)

	return &Minimap{
		scale:            scale,
		cornerPos:        cornerPos,
		cornerSize:       cornerSize,
		centerPos:        centerPos,
		centerSize:       centerSize,
		mapData:          mapData,
		texture:          texture,
		isDirty:          true,
		borderPad:        2,
		isFullscreen:     false,
		lastScreenWidth:  screenWidth,
		lastScreenHeight: screenHeight,
	}
}

func (m *Minimap) ToggleView() {
	m.isFullscreen = !m.isFullscreen
}

func (m *Minimap) Update(playerPos rl.Vector2) {
	m.UpdateDimensions()

	if m.isDirty {
		m.RenderToTexture()
		m.isDirty = false
	}
	// Mark dirty if the player moves
	m.SetDirty()
}

func (m *Minimap) RenderToTexture() {
	rl.BeginTextureMode(m.texture)
	rl.ClearBackground(rl.Black)

	// Draw rooms
	for x := 0; x < helpers.MAP_WIDTH; x++ {
		for y := 0; y < helpers.MAP_HEIGHT; y++ {
			if m.mapData.IsWalkable(x, y) {
				posX := float32(x) * helpers.TILE_SIZE * m.scale
				posY := float32(y) * helpers.TILE_SIZE * m.scale
				size := float32(helpers.TILE_SIZE) * m.scale

				rl.DrawRectangle(
					int32(posX),
					int32(posY),
					int32(size),
					int32(size),
					rl.Gray,
				)
			}
		}
	}

	rl.EndTextureMode()
}

func (m *Minimap) Render(playerPos rl.Vector2) {
	if !m.isFullscreen {
		// Draw corner minimap
		m.renderAt(m.cornerPos, m.cornerSize, playerPos, 3)
	} else {
		// Draw semi-transparent background
		rl.DrawRectangle(0, 0, helpers.SCREEN_WIDTH, helpers.SCREEN_HEIGHT,
			rl.ColorAlpha(rl.Black, 0.7))

		// Draw centered map
		m.renderAt(m.centerPos, m.centerSize, playerPos, 6)

		// Draw instructions
		text := "Press T to close map"
		fontSize := int32(20)
		textWidth := rl.MeasureText(text, fontSize)
		rl.DrawText(text,
			int32(m.centerPos.X+(m.centerSize.X-float32(textWidth))/2),
			int32(m.centerPos.Y+m.centerSize.Y+10),
			fontSize,
			rl.White)
	}
}
func (m *Minimap) renderAt(pos, size rl.Vector2, playerPos rl.Vector2, playerDotSize float32) {
	// Draw outer border (dark gray)
	outerPad := m.borderPad + 1
	rl.DrawRectangle(
		int32(pos.X)-outerPad,
		int32(pos.Y)-outerPad,
		int32(size.X)+outerPad*2,
		int32(size.Y)+outerPad*2,
		rl.ColorAlpha(rl.DarkGray, 0.2),
	)

	// Draw inner background (black)
	rl.DrawRectangle(
		int32(pos.X)-m.borderPad,
		int32(pos.Y)-m.borderPad,
		int32(size.X)+m.borderPad*2,
		int32(size.Y)+m.borderPad*2,
		rl.ColorAlpha(rl.Black, 0.2),
	)

	// Draw the map texture with slight transparency
	rl.DrawTexturePro(
		m.texture.Texture,
		rl.NewRectangle(0, 0, float32(m.texture.Texture.Width), float32(-m.texture.Texture.Height)),
		rl.NewRectangle(pos.X, pos.Y, size.X, size.Y),
		rl.Vector2{X: 0, Y: 0},
		0,
		rl.ColorAlpha(rl.White, 0.4),
	)

	// Calculate player position relative to map size
	relativeX := (playerPos.X / float32(helpers.MAP_WIDTH*helpers.TILE_SIZE)) * size.X
	relativeY := (playerPos.Y / float32(helpers.MAP_HEIGHT*helpers.TILE_SIZE)) * size.Y

	// Draw player dot glow effect
	glowSize := playerDotSize * 2
	rl.DrawCircle(
		int32(pos.X+relativeX),
		int32(pos.Y+relativeY),
		glowSize,
		rl.ColorAlpha(rl.Red, 0.4),
	)

	// Draw player dot with a white outline
	rl.DrawCircle(
		int32(pos.X+relativeX),
		int32(pos.Y+relativeY),
		playerDotSize+1,
		rl.White,
	)
	rl.DrawCircle(
		int32(pos.X+relativeX),
		int32(pos.Y+relativeY),
		playerDotSize,
		rl.Red,
	)

	// Draw the destination marker if it exists
	if m.hasDestination {
		// Calculate destination position relative to map dimensions
		destRelativeX := (float32(m.destinationX) * helpers.TILE_SIZE) / float32(helpers.MAP_WIDTH*helpers.TILE_SIZE)
		destRelativeY := (float32(m.destinationY) * helpers.TILE_SIZE) / float32(helpers.MAP_HEIGHT*helpers.TILE_SIZE)

		destMapX := destRelativeX * size.X
		destMapY := destRelativeY * size.Y

		// Draw pulsing destination marker
		pulseScale := 1.0 + 0.2*float32(math.Sin(float64(rl.GetTime()*4)))
		markerSize := playerDotSize * pulseScale

		// Draw outer glow
		rl.DrawCircle(
			int32(pos.X+destMapX),
			int32(pos.Y+destMapY),
			markerSize*2,
			rl.ColorAlpha(rl.Yellow, 0.2),
		)

		// Draw inner circle
		rl.DrawCircle(
			int32(pos.X+destMapX),
			int32(pos.Y+destMapY),
			markerSize,
			rl.ColorAlpha(rl.Yellow, 0.7),
		)

		// Draw center dot
		rl.DrawCircle(
			int32(pos.X+destMapX),
			int32(pos.Y+destMapY),
			markerSize/2,
			rl.Yellow,
		)
	}
}

func (m *Minimap) SetDirty() {
	m.isDirty = true
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (m *Minimap) SetDestination(x, y int) {
	m.destinationX = x
	m.destinationY = y
	m.hasDestination = true
}

func (m *Minimap) UpdateDimensions() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Skip update if dimensions haven't changed
	if screenWidth == m.lastScreenWidth && screenHeight == m.lastScreenHeight {
		return
	}

	// Store new dimensions
	m.lastScreenWidth = screenWidth
	m.lastScreenHeight = screenHeight

	// Update corner minimap dimensions (20% of screen)
	m.cornerSize = rl.Vector2{
		X: screenWidth * 0.2,
		Y: screenHeight * 0.2,
	}

	m.cornerPos = rl.Vector2{
		X: screenWidth - m.cornerSize.X - 20,
		Y: 20,
	}

	// Update center map view dimensions (70% of screen)
	m.centerSize = rl.Vector2{
		X: screenWidth * 0.7,
		Y: screenHeight * 0.7,
	}

	m.centerPos = rl.Vector2{
		X: (screenWidth - m.centerSize.X) / 2,
		Y: (screenHeight - m.centerSize.Y) / 2,
	}

	// Recalculate scale
	scaleX := m.centerSize.X / float32(helpers.MAP_WIDTH*helpers.TILE_SIZE)
	scaleY := m.centerSize.Y / float32(helpers.MAP_HEIGHT*helpers.TILE_SIZE)
	m.scale = min(scaleX, scaleY)

	// Update texture size
	m.texture = rl.LoadRenderTexture(
		int32(float32(helpers.MAP_WIDTH*helpers.TILE_SIZE)*m.scale),
		int32(float32(helpers.MAP_HEIGHT*helpers.TILE_SIZE)*m.scale),
	)
	m.isDirty = true
}
