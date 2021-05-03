package world

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/spriteloader"
)

type Layers struct {
	Background []Tile
	Mid        []Tile
	Top        []Tile
}

type Planet struct {
	Width  int32
	Height int32
	Layers Layers
}

func NewPlanet(x, y int32) *Planet {
	planet := Planet{
		Width:  x,
		Height: y,
		Layers: Layers{
			Background: make([]Tile, x*y),
			Mid:        make([]Tile, x*y),
			Top:        make([]Tile, x*y),
		},
	}

	return &planet
}

func (planet *Planet) DrawLayer(tiles []Tile, cam *camera.Camera) {
	for idx, tile := range tiles {
		y := int32(idx) / planet.Width
		x := int32(idx) % planet.Width

		if tile.TileType != TileTypeNone {
			spriteloader.DrawSpriteQuad(
				float32(x)-cam.X, float32(y)-cam.Y,
				1, 1,
				int(tile.SpriteID),
			)
		}
	}
}

func (planet *Planet) Draw(cam *camera.Camera) {
	planet.DrawLayer(planet.Layers.Background, cam)
	planet.DrawLayer(planet.Layers.Mid, cam)
	planet.DrawLayer(planet.Layers.Top, cam)
}

func (planet *Planet) GetTileIndex(x, y int) int {
	return y*int(planet.Width) + x
}

// gets the y coordinate of the highest solid tile{
// for a given column given by an x coordinate
func (planet *Planet) GetHeight(x int) int {
	for y := int(planet.Height - 1); y >= 0; y-- {
		idx := planet.GetTileIndex(x, y)
		if planet.Layers.Top[idx].TileType != TileTypeNone {
			return y
		}
	}
	return 0
}

func (planet *Planet) GetTopTile(x, y int) *Tile {
	index := planet.GetTileIndex(x, y)
	if index == -1 {
		return nil
	}

	return &planet.Layers.Top[index]
}
