package world

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
)

type TileCategory uint32

const (
	TileCategoryNone TileCategory = iota
	TileCategoryNormal
	TileCategoryMulti
	TileCategoryChild
	TileCategoryLiquid
)

type TileCollisionType uint32

const (
	TileCollisionTypeSolid TileCollisionType = iota
	TileCollisionTypePlatform
)

func (tt TileCategory) ShouldRender() bool {
	return tt != TileCategoryNone
}

type Tile struct {
	SpriteID          render.SpriteID
	TileCategory      TileCategory
	TileCollisionType TileCollisionType
	TileTypeID        TileTypeID
	Name              string
	OffsetX           int8
	OffsetY           int8
	Durability        int8
	Connections       Connections
	LightSource       bool
	NeedsGround       bool
	Power             TilePower

	FlipTransform     mgl32.Mat4
}

type TilePower struct {
	On bool
	Wattage int
}

func NewEmptyTile() Tile {
	return Tile{TileCategory: TileCategoryNone}
}

func NewNormalTile() Tile {
	return Tile{
		TileCategory:      TileCategoryNormal,
		TileCollisionType: TileCollisionTypeSolid,
		Durability:        1,
		FlipTransform:     mgl32.Ident4(),
	}
}
