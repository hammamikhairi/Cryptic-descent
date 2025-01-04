package screens

import (
	"crydes/audio"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type OutroScreen struct {
	buttons      []*Button
	soundManager *audio.SoundManager
	nextScreen   ScreenType
	startTime    float32
	fadeAlpha    float32
	credits      []CreditEntry
	scrollY      float32
}

type CreditEntry struct {
	role string
	name string
}

func NewOutroScreen(soundManager *audio.SoundManager) *OutroScreen {
	os := &OutroScreen{
		soundManager: soundManager,
		nextScreen:   OUTRO,
		startTime:    0,
		fadeAlpha:    0,
		scrollY:      float32(rl.GetScreenHeight())/2 + 100,
		credits: []CreditEntry{
			{"Game Design", "Khairi Hammami"},
			{"Development", "Khairi Hammami"},
			{"Music", "Mohamed Amine Dridi"},
			{"Sound Effects", "Khairi Hammami"},
			{"Writing", "Khairi Hammami"},
			{"Sprites", "Khairi Hammami"},
			{"Animations", "Khairi Hammami"},
			{"Physics Engine", "Khairi Hammami"},
			// {"AI Programming", "Khairi Hammami"},
			{"Level Design", "Khairi Hammami"},
			{"UI/UX Design", "Khairi Hammami"},
			// {"Networking Code", "Khairi Hammami"},
			// {"Localization", "Khairi Hammami"},
			// {"Marketing", "Khairi Hammami"},
			// {"Documentation", "Khairi Hammami"},
			{"QA Testing", "Khairi Hammami"},
			{"Debugging", "Khairi Hammami"},
			{"Optimization", "Khairi Hammami"},
			{"Coffee Supply", "9ahwet l7ouma"},
			{"Special Thanks", "The Raylib Community"},
			{"Extra Special Thanks", "Khairi Hammami"},
		},
	}
	os.Init()
	return os
}

func (os *OutroScreen) Type() ScreenType {
	return OUTRO
}

func (os *OutroScreen) Init() {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	buttonWidth := float32(200)
	buttonHeight := float32(50)
	startX := (screenWidth - buttonWidth) / 2
	startY := ((screenHeight * 4) / 5) + buttonHeight*2 // Position button lower on

	os.buttons = []*Button{
		NewButton(startX, startY, buttonWidth, buttonHeight, "Back to Title", func() {
			os.soundManager.RequestSound("menu_select", 1.0, 1.0)
			os.nextScreen = TITLE
		}),
	}

	// Start outro music
}

func (os *OutroScreen) Update(deltaTime float32) bool {
	// Check if outro music is playing
	if os.soundManager.GetCurrentMusic() != "outro" {
		os.soundManager.RequestMusic("outro", true)
	}

	os.startTime += deltaTime

	// Update scroll position
	scrollSpeed := float32(50.0) // Adjust speed as needed
	os.scrollY -= scrollSpeed * deltaTime

	// Reset scroll position when all credits have scrolled
	totalHeight := float32(len(os.credits)*30 + 200)
	screenHeight := float32(rl.GetScreenHeight())
	if os.scrollY < screenHeight/2-totalHeight {
		os.scrollY = screenHeight/2 + 100
	}

	// Fade in effect
	if os.fadeAlpha < 1.0 {
		os.fadeAlpha += deltaTime * 0.5 // Adjust fade speed as needed
		if os.fadeAlpha > 1.0 {
			os.fadeAlpha = 1.0
		}
	}

	// Update buttons
	for _, button := range os.buttons {
		button.Update()
	}

	// Return to title screen if requested
	return os.nextScreen == TITLE
}

func (os *OutroScreen) Render() {
	rl.ClearBackground(rl.Black)

	// Draw title with fade effect
	color := rl.ColorAlpha(rl.White, os.fadeAlpha)
	titleText := "Thank you for playing"
	fontSize := int32(60)
	textWidth := rl.MeasureText(titleText, fontSize)
	rl.DrawText(
		titleText,
		int32(float32(rl.GetScreenWidth()-int(textWidth))/2),
		100,
		fontSize,
		color,
	)

	// Draw game name
	subtitleText := "Cryptic Descent"
	subFontSize := int32(40)
	subTextWidth := rl.MeasureText(subtitleText, subFontSize)
	rl.DrawText(
		subtitleText,
		int32(float32(rl.GetScreenWidth()-int(subTextWidth))/2),
		180,
		subFontSize,
		rl.ColorAlpha(rl.Gray, os.fadeAlpha),
	)

	// Draw scrolling credits
	creditsFontSize := int32(20)
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	// Define scroll boundaries
	scrollAreaTop := screenHeight/2 - 100    // Top boundary
	scrollAreaBottom := screenHeight/2 + 200 // Bottom boundary
	fadeDistance := float32(40.0)            // Distance over which the fade occurs

	// Draw scrolling credits
	for i, credit := range os.credits {
		yPos := os.scrollY + float32(i*30)

		// Only draw credits within or near the scroll area
		if yPos >= scrollAreaTop-fadeDistance && yPos <= scrollAreaBottom+fadeDistance {
			// Calculate fade alpha based on position
			fadeAlpha := float32(1.0)

			// Fade in at top
			if yPos < scrollAreaTop {
				fadeAlpha = (scrollAreaTop - yPos) / fadeDistance
				fadeAlpha = 1.0 - fadeAlpha
			}
			// Fade out at bottom
			if yPos > scrollAreaBottom {
				fadeAlpha = (yPos - scrollAreaBottom) / fadeDistance
				fadeAlpha = 1.0 - fadeAlpha
			}

			// Clamp alpha value between 0 and 1
			if fadeAlpha < 0 {
				fadeAlpha = 0
			} else if fadeAlpha > 1 {
				fadeAlpha = 1
			}

			// Apply both the screen fade and position fade
			finalAlpha := fadeAlpha * os.fadeAlpha

			creditText := credit.role + ": " + credit.name
			textWidth := rl.MeasureText(creditText, creditsFontSize)

			rl.DrawText(
				creditText,
				int32((screenWidth-float32(textWidth))/2),
				int32(yPos),
				creditsFontSize,
				rl.ColorAlpha(rl.White, finalAlpha),
			)
		}
	}

	// Draw buttons only after initial fade
	if os.fadeAlpha >= 0.5 {
		for _, button := range os.buttons {
			button.Render()
		}
	}
}

func (os *OutroScreen) Unload() {
	// Stop the outro music when leaving the screen
	os.soundManager.RequestMusic("title_theme", true)
}
