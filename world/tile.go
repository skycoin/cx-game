package world

import (
	"fmt"
)

type TileCategory uint32

const (
	TileCategoryNone TileCategory = iota
	TileCategoryNormal
	TileCategoryMulti
	TileCategoryChild
)

func (tt TileCategory) ShouldRender() bool {
	return tt!=TileCategoryNone
}

type Tile struct {
	SpriteID uint32
	TileCategory TileCategory
	TileTypeID TileTypeID
	Name     string
	OffsetX  int8
	OffsetY  int8
	Durability int8
}

func NewEmptyTile() Tile {
	return Tile {TileCategory: TileCategoryNone}
}

type MultiTile struct {
	Width     int
	Height    int
	TileCategory  TileCategory
	SpriteIDs []uint32
	Name      string
}

// gets the base/root/master tile which controls the multi-tile
func (mt MultiTile) Root() Tile {
	return Tile { SpriteID: mt.SpriteIDs[0], TileCategory: TileCategoryMulti }
}

func (mt MultiTile) Area() int {
	return mt.Width * mt.Height
}

func (mt MultiTile) Tiles() []Tile {
	tiles := make([]Tile,mt.Area())
	tiles[0] = mt.Root()
	for idx:=1; idx<mt.Width*mt.Height; idx++ {
		y := idx / mt.Width
		x := idx % mt.Width
		tiles[idx] = Tile {
			TileCategory: TileCategoryChild,
			SpriteID: mt.SpriteIDs[idx],
			Name: fmt.Sprintf("%s (child [%d,%d])",mt.Name,x,y),
			OffsetX: int8(x),
			OffsetY: int8(y),
		}
	}
	return tiles
}
