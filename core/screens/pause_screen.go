package screens

import (
	"crydes/audio"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PauseScreen struct {
	buttons      []*Button
	soundManager *audio.SoundManager
	nextScreen   ScreenType
}

func NewPauseScreen(soundManager *audio.SoundManager) *PauseScreen {
	ps := &PauseScreen{
		soundManager: soundManager,
		nextScreen:   PAUSE,
	}
	ps.Init()
	return ps
}

func (ps *PauseScreen) Type() ScreenType {
	return PAUSE
}

func (ps *PauseScreen) Init() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	buttonWidth := float32(200)
	buttonHeight := float32(50)
	startX := (screenWidth - buttonWidth) / 2
	startY := screenHeight/2 - buttonHeight

	ps.buttons = []*Button{
		NewButton(startX, startY, buttonWidth, buttonHeight, "Resume", func() {
			ps.soundManager.RequestSound("menu_select", 1.0, 1.0)
			ps.nextScreen = GAME
		}),
		NewButton(startX, startY+buttonHeight+20, buttonWidth, buttonHeight, "Quit", func() {
			ps.soundManager.RequestSound("menu_select", 1.0, 1.0)
			ps.nextScreen = TITLE
		}),
	}
}

func (ps *PauseScreen) Update(deltaTime float32) bool {
	for _, button := range ps.buttons {
		button.Update()
	}
	return ps.nextScreen == GAME
}

func (ps *PauseScreen) Render() {
	// Draw semi-transparent background
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()),
		rl.ColorAlpha(rl.Black, 0.7))

	// Draw pause text
	pauseText := "PAUSED"
	fontSize := int32(60)
	textWidth := rl.MeasureText(pauseText, fontSize)
	rl.DrawText(pauseText, int32(float32(rl.GetScreenWidth()-int(textWidth))/2), 100, fontSize, rl.White)

	// Draw buttons
	for _, button := range ps.buttons {
		button.Render()
	}
}

func (ps *PauseScreen) Unload() {
	// Cleanup if needed
}
