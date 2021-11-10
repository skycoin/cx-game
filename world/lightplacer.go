package world

import (
	"github.com/skycoin/cx-game/render"
)

type LightPlacer struct {
	Tile                    Tile
	OnSpriteID, OffSpriteID render.SpriteID
}

func (placer LightPlacer) CreateTile(
	tt TileType, opts TileCreationOptions,
) Tile {
	tile := placer.Tile
	updateOpts := TileUpdateOptions{
		Tile: &tile,
	}
	placer.UpdateTile(tt, updateOpts)
	return tile
}

func (placer LightPlacer) UpdateTile(
	tt TileType, opts TileUpdateOptions,
) {
	if opts.Tile.Power.On {
		opts.Tile.SpriteID = placer.OnSpriteID
	} else {
		opts.Tile.SpriteID = placer.OffSpriteID
	}
}

func (placer LightPlacer) ItemSpriteID() render.SpriteID {
	return placer.OnSpriteID
}
