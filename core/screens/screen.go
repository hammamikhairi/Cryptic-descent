package screens

import rl "github.com/gen2brain/raylib-go/raylib"

// ScreenType represents different game screens
type ScreenType int

const (
	TITLE ScreenType = iota
	GAME
	PAUSE
	GAME_OVER
	OUTRO
)

// Screen interface defines methods that all screens must implement
type Screen interface {
	Update(deltaTime float32) ScreenType // Returns the next screen to transition to
	Render()
	Init()
	Type() ScreenType
	Unload()
}

// Button represents a clickable UI element
type Button struct {
	Bounds    rl.Rectangle
	Text      string
	IsHovered bool
	Action    func()
}

// NewButton creates a new button with the given parameters
func NewButton(x, y, width, height float32, text string, action func()) *Button {
	return &Button{
		Bounds: rl.NewRectangle(x, y, width, height),
		Text:   text,
		Action: action,
	}
}

// Update checks if the button is being hovered or clicked
func (b *Button) Update() {
	mousePos := rl.GetMousePosition()
	b.IsHovered = rl.CheckCollisionPointRec(mousePos, b.Bounds)

	if b.IsHovered && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		b.Action()
	}
}

// Render draws the button
func (b *Button) Render() {
	color := rl.DarkGray
	if b.IsHovered {
		color = rl.Gray
	}

	rl.DrawRectangleRec(b.Bounds, color)
	rl.DrawRectangleLinesEx(b.Bounds, 2, rl.Black)

	fontSize := int32(30)
	textWidth := rl.MeasureText(b.Text, fontSize)
	textX := b.Bounds.X + (b.Bounds.Width-float32(textWidth))/2
	textY := b.Bounds.Y + (b.Bounds.Height-float32(fontSize))/2

	rl.DrawText(b.Text, int32(textX), int32(textY), fontSize, rl.Black)
}
