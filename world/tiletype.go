package world

import (
	"github.com/skycoin/cx-game/render/blob"
)

type TileCreator func(neighbours blob.Neighbours) Tile
type TileUpdater func(tile *Tile, neighbours blob.Neighbours)
type TileType struct {
	Name string
	Layer Layer
	CreateTile TileCreator
	UpdateTile TileUpdater 
	Invulnerable bool
}

func NewTileType(name string, layer Layer, createTile TileCreator) TileType {
	return TileType { Name: name, Layer: layer, CreateTile: createTile }
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
