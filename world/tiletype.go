package world

import (
	"github.com/skycoin/cx-game/render/blob"
)

type Placer interface {
	CreateTile(TileType,TileCreationOptions) Tile
	UpdateTile(TileType,TileUpdateOptions)
}

type DirectPlacer struct {
	SpriteID uint32
}
func (placer DirectPlacer) CreateTile(
	tt TileType,opts TileCreationOptions,
) Tile {
	return Tile { Name: tt.Name, SpriteID: placer.SpriteID }
}
// nothing to update
func (placer DirectPlacer) UpdateTile(
	tt TileType, opts TileUpdateOptions ) {}

type TileType struct {
	Name string
	Layer Layer
	Placer Placer
	Invulnerable bool
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

type TileTypeID uint32

// add the null tile type as first element such that tileTypes[0] is empty
var tileTypes = make([]TileType,1)

func RegisterTileType(tileType TileType) TileTypeID {
	id := len(tileTypes)
	tileTypes = append(tileTypes, tileType)
	return TileTypeID(id)
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
