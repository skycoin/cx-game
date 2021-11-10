package world

import (
	"github.com/skycoin/cx-game/render"
)

type Placer interface {
	CreateTile(TileType, TileCreationOptions) Tile
	UpdateTile(TileType, TileUpdateOptions)
	ItemSpriteID() render.SpriteID
}

// place tiles for a tiletype which has a single sprite
type DirectPlacer struct {
	SpriteID render.SpriteID
	Tile     Tile
}

func (placer DirectPlacer) CreateTile(
	tt TileType, opts TileCreationOptions,
) Tile {
	tile := placer.Tile
	tile.SpriteID = placer.SpriteID
	tile.FlipTransform = opts.FlipTransform
	return tile
}

// nothing to update
func (placer DirectPlacer) UpdateTile(
	tt TileType, opts TileUpdateOptions) {
}

func (placer DirectPlacer) ItemSpriteID() render.SpriteID {
	return placer.SpriteID
}
