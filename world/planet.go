package world

import (
	//"log"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/camera"
)

type Layers struct {
	Background []Tile
	Mid []Tile
	Top []Tile
}

type Planet struct {
	Width  int32
	Height int32
	Layers Layers
}

func NewPlanet(x, y int32) *Planet {
	tiles := make([]Tile,x*y)
	planet := Planet{
		Width:  x,
		Height: y,
		Layers:  Layers{
			Background: tiles,
			Mid: tiles,
			Top: tiles,
		},
	}

	return &planet
}

func (planet *Planet) DrawLayer(tiles []Tile, cam *camera.Camera) {
	for idx,tile := range tiles {
		y := int32(idx)/planet.Width
		x := int32(idx)%planet.Width

		spriteloader.DrawSpriteQuad(
			float32(x)-cam.X,float32(y)-cam.Y,
			1,1,
			int(tile.SpriteID),
		)
	}
}

func (planet *Planet) Draw(cam *camera.Camera) {
	planet.DrawLayer(planet.Layers.Background, cam)
	planet.DrawLayer(planet.Layers.Mid, cam)
	planet.DrawLayer(planet.Layers.Top, cam)
}
