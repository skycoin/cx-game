package models

import (
	"image"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
	"github.com/skycoin/cx-game/spriteloader"
)

type CatBlack struct {
	RGBA                 *image.RGBA
	Size                 image.Point
	walk                 func()
	width                float32
	height               float32
	X                    float32
	Y                    float32
	XVelocity, YVelocity float32
	movSpeed             float32
	jumpSpeed            float32
	spriteId             int
}

var lspriteSheetId int
var lwindow *glfw.Window

func NewCatBlack(window *glfw.Window) *CatBlack {
	lwindow = window
	spriteloader.InitSpriteloader(lwindow)
	lspriteSheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/blackcat_sprite.png", 13, 4)
	catBlack := CatBlack{
		width:     2,
		height:    1,
		movSpeed:  0.05,
		jumpSpeed: 0.2,
		walk: func() {
			j := 0
			for {
				spriteloader.LoadSprite(lspriteSheetId, "blackcat", 0, j)
				spriteId := spriteloader.GetSpriteIdByName("blackcat")
				gl.ClearColor(1, 1, 1, 1)
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
				time.Sleep(100 * time.Millisecond)
				spriteloader.DrawSpriteQuad(0, 0, 2, 1, spriteId)
				glfw.PollEvents()
				lwindow.SwapBuffers()
				j++
				if j == 11 {
					j = 0
				}
			}
		},
	}

	return &catBlack
}
