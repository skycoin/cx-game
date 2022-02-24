package constants

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DRAW_SPRITE_BATCH_SIZE  int32 = 100
	DRAW_COLOR_BATCH_SIZE   int32 = 100
	PIXELS_PER_TILE         int   = 32
	VIRTUAL_VIEWPORT_WIDTH  int32 = 1440
	VIRTUAL_VIEWPORT_HEIGHT int32 = 900
)

const (
	BACKGROUND_Z float32 = -100
	//2_gap between layers so we can draw something between them
	BGLAYER_Z     float32 = -12
	WINDOWLAYER_Z float32 = -11.5
	PIPELAYER_Z   float32 = -11
	MIDLAYER_Z    float32 = -10
	FRONTLAYER_Z  float32 = -8
	SUPERLAYER_Z  float32 = -6
	AGENT_Z       float32 = -5
	PLAYER_Z      float32 = -3
	HUD_Z         float32 = 2
)

var (
	OUTLINE_BORDER_COLOR = mgl32.Vec4{0.106, 0.106, 0.106, 1}
)

// func ChangeResolution(resolutionTyep int32) (int, int) {
// 	if resolutionTyep == 1 {
// 		Resolution_Selected++
// 	} else {
// 		Resolution_Selected--
// 	}
// 	if Resolution_Selected == 1 {
// 		VIRTUAL_VIEWPORT_WIDTH2 = 1280
// 		VIRTUAL_VIEWPORT_HEIGHT2 = 720
// 	} else if Resolution_Selected == 2 {
// 		VIRTUAL_VIEWPORT_WIDTH2 = 1366
// 		VIRTUAL_VIEWPORT_HEIGHT2 = 768
// 	} else if Resolution_Selected == 3 {
// 		VIRTUAL_VIEWPORT_WIDTH2 = 1440
// 		VIRTUAL_VIEWPORT_HEIGHT2 = 900
// 	} else if Resolution_Selected == 4 {
// 		VIRTUAL_VIEWPORT_WIDTH2 = 1536
// 		VIRTUAL_VIEWPORT_HEIGHT2 = 864
// 	} else if Resolution_Selected == 5 {
// 		VIRTUAL_VIEWPORT_WIDTH2 = 1600
// 		VIRTUAL_VIEWPORT_HEIGHT2 = 900
// 	} else if Resolution_Selected == 6 {
// 		VIRTUAL_VIEWPORT_WIDTH2 = 1920
// 		VIRTUAL_VIEWPORT_HEIGHT2 = 1080
// 	}

// 	return int(VIRTUAL_VIEWPORT_HEIGHT2), int(VIRTUAL_VIEWPORT_WIDTH2)
// }
