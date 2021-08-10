package render

import (
	"strconv"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/constants"
)

var spriteProgram Program

func initSprite() {
	config := NewShaderConfig(
		"./assets/shader/spritesheet.vert", "./assets/shader/spritesheet.frag",
	)
	config.Define("NUM_INSTANCES",
		strconv.Itoa(int(constants.DRAW_SPRITE_BATCH_SIZE)) )

	spriteProgram = config.Compile()
}

func Init() { initSprite() }

type Sprite struct {
	Name string
	Transform mgl32.Mat3
	Model mgl32.Mat4
	Texture Texture
}

// deprecate at some point
type SpriteSheet struct {
	Texture Texture
	Sprites []Sprite
}
