package screens

import (
	"crydes/audio"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VictoryScreen struct {
	soundManager *audio.SoundManager
	nextScreen   ScreenType
	timer        float32
	fadeAlpha    float32
}

func NewVictoryScreen(soundManager *audio.SoundManager) *VictoryScreen {
	vs := &VictoryScreen{
		soundManager: soundManager,
		nextScreen:   VICTORY,
		timer:        0,
		fadeAlpha:    0,
	}
	vs.Init()
	return vs
}

func (vs *VictoryScreen) Type() ScreenType {
	return VICTORY
}

func (vs *VictoryScreen) Init() {
	// Play victory sound/music
	vs.soundManager.RequestSound("victory", 1.0, 1.0)
}

func (vs *VictoryScreen) Update(deltaTime float32) bool {
	vs.timer += deltaTime

	// Fade in effect
	if vs.fadeAlpha < 1.0 {
		vs.fadeAlpha += deltaTime * 2.0 // Adjust fade speed as needed
		if vs.fadeAlpha > 1.0 {
			vs.fadeAlpha = 1.0
		}
	}

	// After 5 seconds, transition to outro screen
	if vs.timer >= 5.0 {
		return true
	}

	return false
}

func (vs *VictoryScreen) Render() {
	rl.ClearBackground(rl.Black)

	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Draw victory message
	titleText := "Victory!"
	fontSize := int32(80)
	textWidth := rl.MeasureText(titleText, fontSize)
	color := rl.ColorAlpha(rl.Gold, vs.fadeAlpha)

	rl.DrawText(
		titleText,
		int32(screenWidth/2-float32(textWidth)/2),
		int32(screenHeight/2-float32(fontSize)),
		fontSize,
		color,
	)

	// Draw sub-message
	subText := "You've collected all the keys!"
	subFontSize := int32(40)
	subTextWidth := rl.MeasureText(subText, subFontSize)
	subColor := rl.ColorAlpha(rl.White, vs.fadeAlpha)

	rl.DrawText(
		subText,
		int32(screenWidth/2-float32(subTextWidth)/2),
		int32(screenHeight/2+float32(fontSize)/2),
		subFontSize,
		subColor,
	)
}

func (vs *VictoryScreen) Unload() {
	// Clean up any resources if needed
}
