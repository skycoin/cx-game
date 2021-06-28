package world

import (
	"fmt"
)

type TileType uint32

const (
	TileTypeNone TileType = iota
	TileTypeNormal
	TileTypeMulti
	TileTypeChild
)

func (tt TileType) ShouldRender() bool {
	return tt!=TileTypeNone
}

type Tile struct {
	SpriteID uint32
	TileType TileType
	Name     string
	OffsetX  int8
	OffsetY  int8
	Durability int8

	IsBlob   bool
	BlobSpriteID uint32
}

func NewEmptyTile() Tile {
	return Tile {TileType: TileTypeNone}
}

type MultiTile struct {
	Width     int
	Height    int
	TileType  TileType
	SpriteIDs []uint32
	Name      string
}

// gets the base/root/master tile which controls the multi-tile
func (mt MultiTile) Root() Tile {
	return Tile { SpriteID: mt.SpriteIDs[0], TileType: TileTypeMulti }
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
			TileType: TileTypeChild,
			SpriteID: mt.SpriteIDs[idx],
			Name: fmt.Sprintf("%s (child [%d,%d])",mt.Name,x,y),
			OffsetX: int8(x),
			OffsetY: int8(y),
		}
	}
	return tiles
}
