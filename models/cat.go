package models

import (
	"image"
	"github.com/skycoin/cx-game/spriteloader"
)

type Cat struct {
	RGBA      *image.RGBA
	Size      image.Point
	width     int
	height    int
	XVelocity float32
	YVelocity float32
	spriteId  int
}

func NewCat() *Cat {
	spriteId := spriteloader.LoadSingleSprite("./assets/cat.png","cat")
	cat := Cat{
		width:  2,
		height: 2,
		spriteId: spriteId,
	}

	return &cat
}

func (cat *Cat) Draw() {
	spriteloader.DrawSpriteQuad(0,0,1,1,cat.spriteId)
}
