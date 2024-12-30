package effects

import (
	"crydes/helpers"
	"crydes/player"
	"crydes/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightSource struct {
	position rl.Vector2
	radius   float32
	mode     string
	isPlayer bool
	player   *player.Player
}

func (ls LightSource) Position() rl.Vector2 {
	if ls.isPlayer {
		return ls.player.GetPlayerCenterPoint()
	}

	return ls.position
}

func (ls LightSource) Radius() float32 {
	return ls.radius
}
func (ls *LightSource) SetMode(mode string) {
	ls.mode = mode
}
func (ls LightSource) Mode() string {
	return ls.mode
}
func (ls *LightSource) SetRadius(rad float32) {
	ls.radius = rad
}

type LightSourceIf interface {
	Position() rl.Vector2
	Radius() float32
	Mode() string
	SetRadius(float32)
	SetMode(string)
}

type RetroLightingEffect struct {
	lightMask   rl.RenderTexture2D
	lightRadius float32
	pixelSize   int32
	player      *player.Player
	smoothness  float32

	lightSources []LightSourceIf

	modeOrder        []string
	currentModeIndex int
	modes            map[string]func(rl.Vector2, float32)
	noiseMap         [][]float32
	time             float32
}

func NewRetroLightingEffect(screenWidth, screenHeight int32, lightRadius float32, pixelSize int32, p *player.Player) *RetroLightingEffect {
	rle := &RetroLightingEffect{
		lightMask:   rl.LoadRenderTexture(screenWidth, screenHeight),
		lightRadius: lightRadius,
		pixelSize:   pixelSize,
		player:      p,
		smoothness:  0.8,
		noiseMap:    generateNoiseMap(int(screenWidth/pixelSize), int(screenHeight/pixelSize)),
		modeOrder: []string{
			"static", "shimmer", "pulse", "flicker", "noise", "rainbow", "spiral",
			"strobe", "gradient", "ripple", "vortex", "glitch", "heartbeat", "halo",
			"electric", "kaleidoscope",
		},
		currentModeIndex: 0,
		lightSources:     []LightSourceIf{},
	}

	rle.modes = map[string]func(rl.Vector2, float32){
		"static":       rle.HandleStaticLightning,
		"shimmer":      rle.HandleShimmerLightning,
		"pulse":        rle.HandlePulseLighting,
		"flicker":      rle.HandleFlickerLighting,
		"noise":        rle.HandleNoiseLighting,
		"rainbow":      rle.HandleRainbowLighting,
		"spiral":       rle.HandleSpiralLighting,
		"strobe":       rle.HandleStrobeLighting,
		"gradient":     rle.HandleGradientLighting,
		"ripple":       rle.HandleRippleLighting,
		"vortex":       rle.HandleVortexLighting,
		"glitch":       rle.HandleGlitchLighting,
		"heartbeat":    rle.HandleHeartbeatLighting,
		"halo":         rle.HandleHaloLighting,
		"electric":     rle.HandleElectricLighting,
		"kaleidoscope": rle.HandleKaleidoscopeLighting,
	}

	rle.lightSources = append(rle.lightSources, &LightSource{
		position: rle.player.GetPlayerCenterPoint(),
		radius:   helpers.LIGHT_RADIUS,
		mode:     "static",
		isPlayer: true,
		player:   rle.player,
	})

	return rle
}

func (rle *RetroLightingEffect) NextLightningMode() {
	rle.currentModeIndex++
	if rle.currentModeIndex >= len(rle.modeOrder) {
		rle.currentModeIndex = 0
	}

	println(rle.modeOrder[rle.currentModeIndex])
	// for i := range rle.lightSources {
	rle.lightSources[0].SetMode(rle.modeOrder[rle.currentModeIndex])
	// }
	println("---> ", rle.lightSources[0].Mode())
}

func (rle *RetroLightingEffect) SetUpPropsLightning(props *[]*world.Prop) {
	for _, prop := range *props {
		rle.AddLightSource(
			rl.Vector2{X: prop.Position.X + float32(prop.CurrentAnim.Frames[0].Width/2)*prop.Scale, Y: prop.Position.Y + float32(prop.CurrentAnim.Frames[0].Height)*prop.Scale},
			false,
			prop.LTRadius,
			"static",
		)
	}
}

func (rle *RetroLightingEffect) AddLightSource(position rl.Vector2, isPlayer bool, radius float32, mode string) {

	rle.lightSources = append(rle.lightSources, &LightSource{
		position: position,
		radius:   radius,
		mode:     mode,
		// isStatic: isStatic,
		isPlayer: isPlayer,
		player:   nil,
	})
}

func (rle *RetroLightingEffect) Update() {
	rle.time += rl.GetFrameTime()

	rl.BeginTextureMode(rle.lightMask)
	rl.ClearBackground(rl.Black)

	for i := range rle.lightSources {
		// rle.modes[rle.lightSources[i].Mode()](rle.lightSources[i].Position(), rle.lightSources[i].Radius())
		rle.modes[rle.lightSources[i].Mode()](rle.lightSources[i].Position(), helpers.LIGHT_RADIUS)
	}

	rl.EndTextureMode()
}

func (rle *RetroLightingEffect) Render() {
	rl.BeginBlendMode(rl.BlendMultiplied)
	rl.DrawTextureRec(rle.lightMask.Texture,
		rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(rle.lightMask.Texture.Width),
			Height: float32(-rle.lightMask.Texture.Height),
		},
		rl.Vector2{X: 0, Y: 0},
		rl.White)
	rl.EndBlendMode()
}

func (rle *RetroLightingEffect) Unload() {
	rl.UnloadRenderTexture(rle.lightMask)
}
