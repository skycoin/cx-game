package render

import (
	"strconv"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/constants"
)

var (
	//active sprite program
	spriteProgram *Program
	//no pixel snapping
	spriteProgram1 Program
	//pixel snap to 1.0/32
	spriteProgram2 Program
	//pixel snap to 1.0/32+ 1.0/64
	spriteProgram3 Program
)

func initSprite() {
	config1 := NewShaderConfig(
		"./assets/shader/spritesheet.vert", "./assets/shader/spritesheet.frag",
	)
	config2 := NewShaderConfig(
		"./assets/shader/spritesheetv2.vert", "./assets/shader/spritesheetv2.frag",
	)
	config3 := NewShaderConfig(
		"./assets/shader/spritesheetv3.vert", "./assets/shader/spritesheetv3.frag",
	)

	config1.Define("NUM_INSTANCES",
		strconv.Itoa(int(constants.DRAW_SPRITE_BATCH_SIZE)))

	config2.Define("NUM_INSTANCES",
		strconv.Itoa(int(constants.DRAW_SPRITE_BATCH_SIZE)))

	config3.Define("NUM_INSTANCES",
		strconv.Itoa(int(constants.DRAW_SPRITE_BATCH_SIZE)))

	spriteProgram1 = config1.Compile()
	spriteProgram2 = config2.Compile()
	spriteProgram3 = config3.Compile()

	//default program
	spriteProgram = &spriteProgram1
}


type Sprite struct {
	Name      string
	Transform mgl32.Mat3
	Model     mgl32.Mat4
	Texture   Texture
}

// deprecate at some point
type SpriteSheet struct {
	Texture Texture
	Sprites []Sprite
}
