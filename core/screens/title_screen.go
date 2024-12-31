package screens

import (
	"crydes/audio"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TitleScreen struct {
	buttons      []*Button
	soundManager *audio.SoundManager
	nextScreen   ScreenType
}

func NewTitleScreen(soundManager *audio.SoundManager) *TitleScreen {
	ts := &TitleScreen{
		soundManager: soundManager,
		nextScreen:   TITLE,
	}
	ts.Init()
	return ts
}

func (ts *TitleScreen) Init() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	buttonWidth := float32(200)
	buttonHeight := float32(50)
	startX := (screenWidth - buttonWidth) / 2
	startY := screenHeight/2 - buttonHeight

	ts.buttons = []*Button{
		NewButton(startX, startY, buttonWidth, buttonHeight, "Play", func() {
			ts.soundManager.RequestSound("menu_select", 1.0, 1.0)
			ts.nextScreen = GAME
		}),
		NewButton(startX, startY+buttonHeight+20, buttonWidth, buttonHeight, "Quit", func() {
			ts.soundManager.RequestSound("menu_select", 1.0, 1.0)
			rl.CloseWindow()
		}),
	}

	// Start title screen music
	ts.soundManager.RequestMusic("title_theme", true)
}

func (ts *TitleScreen) Update(deltaTime float32) bool {
	for _, button := range ts.buttons {
		button.Update()
	}

	if rl.IsKeyDown(rl.KeyP) {
		ts.soundManager.RequestMusic("dungeon_theme", true)
		return true
	}

	if ts.nextScreen == GAME {
		ts.soundManager.RequestMusic("dungeon_theme", true)
	}

	return ts.nextScreen == GAME
}

func (ts *TitleScreen) Type() ScreenType {
	return TITLE
}

func (ts *TitleScreen) Render() {
	rl.ClearBackground(rl.RayWhite)

	// Draw title
	titleText := "Cryptic Descent"
	fontSize := int32(60)
	textWidth := rl.MeasureText(titleText, fontSize)
	rl.DrawText(titleText, int32(float32(rl.GetScreenWidth()-int(textWidth))/2), 100, fontSize, rl.Black)

	// Draw buttons
	for _, button := range ts.buttons {
		button.Render()
	}
}

func (ts *TitleScreen) Unload() {
	// Cleanup if needed
}
