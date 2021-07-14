package world

import (
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
)

type Placer interface {
	CreateTile(TileType,TileCreationOptions) Tile
	UpdateTile(TileType,TileUpdateOptions)
	ItemSpriteID() spriteloader.SpriteID
}

// place tiles for a tiletype which has a single sprite
type DirectPlacer struct {
	SpriteID spriteloader.SpriteID
}
func (placer DirectPlacer) CreateTile(
	tt TileType,opts TileCreationOptions,
) Tile {
	return Tile {
		Name: tt.Name,
		SpriteID: placer.SpriteID,
		TileTypeID: tt.ID,
		TileCategory: TileCategoryNormal,
	}
}
// nothing to update
func (placer DirectPlacer) UpdateTile(
	tt TileType, opts TileUpdateOptions ) {}

func (placer DirectPlacer) ItemSpriteID() spriteloader.SpriteID {
	return placer.SpriteID
}

type TileTypeID uint32
type TileType struct {
	Name string
	Layer Layer
	Placer Placer
	Invulnerable bool
	ID TileTypeID
	MaterialID MaterialID
	Width,Height int32
	Drops Drops
	ItemSpriteID spriteloader.SpriteID
}

func (tt TileType) Size() cxmath.Vec2i {
	return cxmath.Vec2i{tt.Width,tt.Height}
}

type TileCreationOptions struct {
	Neighbours blob.Neighbours
}
type TileUpdateOptions struct {
	Neighbours blob.Neighbours
	Tile *Tile
}

func (tt TileType) CreateTile(opts TileCreationOptions) Tile {
	return tt.Placer.CreateTile(tt,opts)
}

func (tt TileType) UpdateTile(opts TileUpdateOptions) {
	tt.Placer.UpdateTile(tt,opts)
}

// add the null tile type as first element such that tileTypes[0] is empty
var tileTypes = make([]TileType,1)

func RegisterTileType(tileType TileType) TileTypeID {
	id := TileTypeID(len(tileTypes))
	tileType.ID = id
	// fill in default size
	if tileType.Width==0 { tileType.Width=1 }
	if tileType.Height==0 { tileType.Height=1 }
	tileType.ItemSpriteID = tileType.Placer.ItemSpriteID()
	if tileType.Drops==nil { tileType.Drops=Drops{} }
	tileTypes = append(tileTypes, tileType)
	return id
}

func NextTileTypeID() TileTypeID {
	return TileTypeID(len(tileTypes))
}

func GetTileTypeByID(id TileTypeID) (TileType,bool) {
	ok :=  id >=1 && id < TileTypeID(len(tileTypes))
	if ok {
		return tileTypes[id],true
	} else {
		return TileType{},false
	}
}

func (id TileTypeID) Get() *TileType {
	return &tileTypes[id]
}

func AddDrop(id TileTypeID, drop Drop) {
	tileTypes[id].Drops = append(tileTypes[id].Drops, drop)
}
