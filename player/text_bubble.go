package player

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// Bubble appearance
	BUBBLE_PADDING    = 10
	BUBBLE_FONT_SIZE  = 20
	BUBBLE_MAX_WIDTH  = 200
	BUBBLE_FADE_SPEED = 2.0
	BUBBLE_SHOW_TIME  = 5.0 // How long to show the bubble
	BUBBLE_FADE_TIME  = 0.5 // How long to fade in/out
	BUBBLE_OFFSET_Y   = -30 // Offset above player
)

var (
	BUBBLE_BACKGROUND = rl.NewColor(0, 0, 0, 180)
	BUBBLE_TEXT_COLOR = rl.White
	BUBBLE_BORDER     = rl.NewColor(255, 255, 255, 80)
)

// Message categories
const (
	// Tutorial messages
	MSG_MOVEMENT = "Use WASD or arrow keys to move"
	MSG_ATTACK   = "Press SPACE to attack"

	// Story messages
	MSG_DUNGEON_ENTER = "This place feels ancient..."
	MSG_DUNGEON_SHIFT = "Something's not right..."
	MSG_FOUND_ENEMY   = "Dangerous creatures lurk here..."

	// Status messages
	MSG_LOW_HEALTH  = "I need to find healing..."
	MSG_POISONED    = "This poison burns..."
	MSG_SPEED_BOOST = "I feel faster!"

	// SHIFTING
	FIRST_SHIFT = "What happened?"
)

type TextBubble struct {
	text          string
	isVisible     bool
	alpha         float32
	fadeIn        bool
	showStartTime time.Time
	width         float32
	height        float32
	wrappedText   []string
}

func NewTextBubble() *TextBubble {
	return &TextBubble{
		alpha:     0,
		isVisible: false,
	}
}

func (tb *TextBubble) ShowMessage(text string) {
	tb.text = text
	tb.isVisible = true
	tb.fadeIn = true
	tb.alpha = 0
	tb.showStartTime = time.Now()

	// Calculate wrapped text and bubble dimensions
	tb.wrappedText = wrapText(text, BUBBLE_MAX_WIDTH, BUBBLE_FONT_SIZE)
	tb.width = float32(BUBBLE_MAX_WIDTH + BUBBLE_PADDING*2)
	tb.height = float32(len(tb.wrappedText)*BUBBLE_FONT_SIZE + BUBBLE_PADDING*2)
}

func (tb *TextBubble) Update(deltaTime float32) {
	if !tb.isVisible {
		return
	}

	elapsed := time.Since(tb.showStartTime).Seconds()

	// Handle fade in
	if tb.fadeIn {
		tb.alpha += BUBBLE_FADE_SPEED * deltaTime
		if tb.alpha >= 1.0 {
			tb.alpha = 1.0
			tb.fadeIn = false
		}
	} else if elapsed >= BUBBLE_SHOW_TIME {
		// Handle fade out
		tb.alpha -= BUBBLE_FADE_SPEED * deltaTime
		if tb.alpha <= 0 {
			tb.alpha = 0
			tb.isVisible = false
		}
	}
}

func (tb *TextBubble) Render(playerPos rl.Vector2) {
	if !tb.isVisible {
		return
	}

	// Calculate bubble position (centered on screen)
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	x := (screenWidth - tb.width) / 2
	y := (screenHeight-tb.height)/2 + BUBBLE_OFFSET_Y

	// Draw bubble background with alpha
	bgColor := rl.ColorAlpha(BUBBLE_BACKGROUND, tb.alpha)
	borderColor := rl.ColorAlpha(BUBBLE_BORDER, tb.alpha)
	textColor := rl.ColorAlpha(BUBBLE_TEXT_COLOR, tb.alpha)

	// Draw background and border
	rl.DrawRectangle(int32(x), int32(y), int32(tb.width), int32(tb.height), bgColor)
	rl.DrawRectangleLinesEx(
		rl.Rectangle{X: x, Y: y, Width: tb.width, Height: tb.height},
		2,
		borderColor,
	)

	// Draw text lines
	for i, line := range tb.wrappedText {
		textX := x + BUBBLE_PADDING
		textY := y + BUBBLE_PADDING + float32(i*BUBBLE_FONT_SIZE)
		rl.DrawText(line, int32(textX), int32(textY), BUBBLE_FONT_SIZE, textColor)
	}
}

// Helper function to wrap text
func wrapText(text string, maxWidth float32, fontSize int32) []string {
	var lines []string
	words := splitWords(text)
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if rl.MeasureText(testLine, fontSize) < int32(maxWidth) {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func splitWords(text string) []string {
	var words []string
	var currentWord string

	for _, char := range text {
		if char == ' ' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}

	if currentWord != "" {
		words = append(words, currentWord)
	}

	return words
}
