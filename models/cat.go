package models

import (
	"image"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/physics"
)

type Cat struct {
	RGBA      *image.RGBA
	Size      image.Point
	width     float32
	height    float32
	X         float32
	Y         float32
	XVelocity,YVelocity float32
	movSpeed  float32
	jumpSpeed float32
	spriteId  int
}

func NewCat() *Cat {
	spriteId := spriteloader.LoadSingleSprite("./assets/cat.png","cat")
	cat := Cat{
		width:  1,
		height: 1,
		movSpeed: 0.05,
		jumpSpeed: 0.2,
		spriteId: spriteId,
	}

	return &cat
}

func (cat *Cat) Draw(cam *camera.Camera) {
	spriteloader.DrawSpriteQuad(
		cat.X-cam.X,cat.Y-cam.Y,
		cat.width,cat.height,cat.spriteId,
	)
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func (cat *Cat) Tick(leftPressed,rightPressed,spacePressed bool) {
	// TODO
	if cat.Y > 6.5 {
		cat.YVelocity -= physics.Gravity
	} else {
		cat.YVelocity = 0

		if spacePressed {
			cat.YVelocity = 0.2
		}
	}
	// FIXME speeds are currently based on ticks instead of time
	xDirection := boolToFloat(rightPressed)-boolToFloat(leftPressed)
	cat.XVelocity = xDirection*cat.movSpeed

	cat.YVelocity = boolToFloat(spacePressed)*cat.jumpSpeed

	cat.X += cat.XVelocity
	cat.Y += cat.YVelocity
}
