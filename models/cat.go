package models

import (
	"image"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
)

type Cat struct {
	physics.Body
	RGBA      *image.RGBA
	ImageSize image.Point
	movSpeed  float32
	jumpSpeed float32
	spriteId  int
}

func NewCat() *Cat {
	spriteId := spriteloader.LoadSingleSprite("./assets/cat.png", "cat")
	cat := Cat{
		Body: physics.Body{
			Size: physics.Vec2{X: 2.0, Y: 2.0},
		},
		movSpeed:  3.0,
		jumpSpeed: 12.0,
		spriteId:  spriteId,
	}
	return &cat
}

func (cat *Cat) Draw(cam *camera.Camera) {

	x := cat.Pos.X - cam.X
	y := cat.Pos.Y - cam.Y
	if !cam.IsInBoundsF(cat.Pos.X, cat.Pos.Y) {
		return
	}

	spriteloader.DrawSpriteQuad(
		x, y,
		cat.Size.X, cat.Size.Y, cat.spriteId,
	)
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func (cat *Cat) Tick(leftPressed, rightPressed, spacePressed bool, planet *world.Planet, dt float32) {
	cat.Vel.X = (boolToFloat(rightPressed) - boolToFloat(leftPressed)) * cat.movSpeed
	cat.Vel.Y -= physics.Gravity * dt

	if spacePressed {
		cat.Vel.Y = cat.jumpSpeed
	}

	cat.Move(planet, dt)
}
