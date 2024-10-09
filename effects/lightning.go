package effects

import (
	"crydes/helpers"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (rle *RetroLightingEffect) drawLightCircle(center rl.Vector2, radius float32, colorFunc func(float32) rl.Color) {
	for x := int32(0); x < int32(rle.lightMask.Texture.Width); x += rle.pixelSize {
		for y := int32(0); y < int32(rle.lightMask.Texture.Height); y += rle.pixelSize {
			dx := float32(x) - center.X
			dy := float32(y) - center.Y
			distance := rl.Vector2Length(rl.Vector2{X: dx, Y: dy})

			if distance < radius {
				intensity := float32(math.Exp(-helpers.DECAY_FACTOR * float64(distance/radius)))
				intensity = float32(math.Min(float64(intensity)*1.5, 1.0))

				color := colorFunc(intensity)
				rl.DrawRectangle(x, y, rle.pixelSize, rle.pixelSize, color)
			}
		}
	}
}

func (rle *RetroLightingEffect) HandleStaticLightning(center rl.Vector2, radius float32) {
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		return rl.ColorAlpha(rl.White, intensity)
	})
}

func (rle *RetroLightingEffect) HandleShimmerLightning(center rl.Vector2, radius float32) {
	shimmerSpeed := 2.0
	shimmerIntensity := 0.1
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		shimmer := float32(math.Sin(float64(rle.time*float32(shimmerSpeed) + intensity*10)))
		shimmer = (shimmer + 1) * 0.5
		finalIntensity := intensity + shimmer*float32(shimmerIntensity)
		finalIntensity = float32(math.Max(0, math.Min(1, float64(finalIntensity))))
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

func (rle *RetroLightingEffect) HandlePulseLighting(center rl.Vector2, radius float32) {
	pulseSpeed := 6.0
	pulseIntensity := 0.2
	pulse := float32(math.Sin(float64(rle.time * float32(pulseSpeed))))
	pulse = (pulse + 1) * 0.5 * float32(pulseIntensity)
	rle.drawLightCircle(center, radius*(1+pulse), func(intensity float32) rl.Color {
		return rl.ColorAlpha(rl.White, intensity)
	})
}

func (rle *RetroLightingEffect) HandleFlickerLighting(center rl.Vector2, radius float32) {
	flickerSpeed := 10.0
	flickerIntensity := 0.3
	flicker := rand.Float32() * float32(flickerIntensity)
	if math.Sin(float64(rle.time*float32(flickerSpeed))) > 0 {
		flicker = 0
	}
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		finalIntensity := intensity + flicker
		finalIntensity = float32(math.Max(0, math.Min(1, float64(finalIntensity))))
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

func (rle *RetroLightingEffect) HandleNoiseLighting(center rl.Vector2, radius float32) {
	noiseScale := 0.1
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		x := int((center.X + float32(math.Cos(float64(intensity*2*math.Pi))*float64(radius))) / float32(rle.pixelSize))
		y := int((center.Y + float32(math.Sin(float64(intensity*2*math.Pi))*float64(radius))) / float32(rle.pixelSize))
		noise := rle.noiseMap[y%len(rle.noiseMap)][x%len(rle.noiseMap[0])]
		finalIntensity := intensity + float32(noiseScale)*noise
		finalIntensity = float32(math.Max(0, math.Min(1, float64(finalIntensity))))
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

func (rle *RetroLightingEffect) HandleRainbowLighting(center rl.Vector2, radius float32) {
	rainbowSpeed := 1.0
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		hue := float32(math.Mod(float64(rle.time*float32(rainbowSpeed)+intensity*360), 360))
		return rl.ColorAlpha(rl.ColorFromHSV(hue, 1, 1), intensity)
	})
}

func (rle *RetroLightingEffect) HandleSpiralLighting(center rl.Vector2, radius float32) {
	spiralSpeed := 2.0
	spiralTightness := 10.0
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		angle := intensity*2*math.Pi + float32(rle.time*float32(spiralSpeed))
		spiral := float32(math.Sin(float64(angle * float32(spiralTightness))))
		finalIntensity := intensity * (0.5 + 0.5*spiral)
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

func generateNoiseMap(width, height int) [][]float32 {
	noiseMap := make([][]float32, height)
	for y := range noiseMap {
		noiseMap[y] = make([]float32, width)
		for x := range noiseMap[y] {
			noiseMap[y][x] = rand.Float32()
		}
	}
	return noiseMap
}

// Strobe effect: Alternates between bright and dark
func (rle *RetroLightingEffect) HandleStrobeLighting(center rl.Vector2, radius float32) {
	strobeSpeed := 5.0
	strobe := float32(math.Sin(float64(rle.time * float32(strobeSpeed))))
	intensity := float32(math.Max(0, math.Sin(float64(strobe))))
	rle.drawLightCircle(center, radius, func(_ float32) rl.Color {
		return rl.ColorAlpha(rl.White, intensity)
	})
}

// Gradient effect: Creates a gradient from center to edge
func (rle *RetroLightingEffect) HandleGradientLighting(center rl.Vector2, radius float32) {
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		return rl.ColorAlpha(rl.ColorFromHSV(120*(1-intensity), 1, 1), intensity)
	})
}

// Ripple effect: Creates expanding circular waves
func (rle *RetroLightingEffect) HandleRippleLighting(center rl.Vector2, radius float32) {
	rippleSpeed := 2.0
	rippleFrequency := float32(0.1)
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		distance := (1 - intensity) * radius
		ripple := float32(math.Sin(float64(distance*rippleFrequency - rle.time*float32(rippleSpeed))))
		finalIntensity := intensity * (0.7 + 0.3*ripple)
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

// Vortex effect: Creates a swirling vortex pattern
func (rle *RetroLightingEffect) HandleVortexLighting(center rl.Vector2, radius float32) {
	vortexSpeed := 2.0
	vortexTightness := 5.0
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		angle := float32(math.Atan2(float64(center.Y-center.Y), float64(center.X-center.X)))
		vortex := float32(math.Sin(float64(angle*float32(vortexTightness) + rle.time*float32(vortexSpeed))))
		finalIntensity := intensity * (0.7 + 0.3*vortex)
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

// Glitch effect: Creates random glitch-like artifacts
func (rle *RetroLightingEffect) HandleGlitchLighting(center rl.Vector2, radius float32) {
	glitchIntensity := 0.2
	glitchSpeed := 10.0
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		if rand.Float32() < float32(glitchIntensity) && int(rle.time*float32(glitchSpeed))%2 == 0 {
			return rl.ColorAlpha(rl.White, rand.Float32())
		}
		return rl.ColorAlpha(rl.White, intensity)
	})
}

// Heartbeat effect: Simulates a heartbeat pattern
func (rle *RetroLightingEffect) HandleHeartbeatLighting(center rl.Vector2, radius float32) {
	beatSpeed := 1.0
	beatIntensity := 0.3
	time := float64(rle.time * float32(beatSpeed))
	beat := float32(math.Pow(math.Sin(time*math.Pi), 63) + math.Pow(math.Sin((time+0.25)*math.Pi), 63))
	rle.drawLightCircle(center, radius*(1+beat*float32(beatIntensity)), func(intensity float32) rl.Color {
		return rl.ColorAlpha(rl.White, intensity*(1+beat*float32(beatIntensity)))
	})
}

// Halo effect: Creates a halo around the edge of the light
func (rle *RetroLightingEffect) HandleHaloLighting(center rl.Vector2, radius float32) {
	haloWidth := 0.1
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		haloIntensity := float32(math.Max(0, 1-math.Abs(float64(intensity-float32(haloWidth))/float64(haloWidth))))
		finalIntensity := float32(math.Max(float64(intensity), float64(haloIntensity)))
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

// Electric effect: Creates lightning-like tendrils
func (rle *RetroLightingEffect) HandleElectricLighting(center rl.Vector2, radius float32) {
	electricSpeed := 5.0
	electricIntensity := 0.3
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		angle := float32(math.Atan2(float64(center.Y-center.Y), float64(center.X-center.X)))
		electric := float32(math.Sin(float64(angle*10 + rle.time*float32(electricSpeed))))
		finalIntensity := intensity + electric*float32(electricIntensity)*(1-intensity)
		return rl.ColorAlpha(rl.White, finalIntensity)
	})
}

// Kaleidoscope effect: Creates a kaleidoscope-like pattern
func (rle *RetroLightingEffect) HandleKaleidoscopeLighting(center rl.Vector2, radius float32) {
	segments := 6
	rotationSpeed := 1.0
	rle.drawLightCircle(center, radius, func(intensity float32) rl.Color {
		angle := float32(math.Atan2(float64(center.Y-center.Y), float64(center.X-center.X)))
		segment := float32(math.Floor(float64(angle/(2*math.Pi/float32(segments)) + rle.time*float32(rotationSpeed))))
		hue := segment * 360 / float32(segments)
		return rl.ColorAlpha(rl.ColorFromHSV(hue, 1, 1), intensity)
	})
}
