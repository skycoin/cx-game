package world

import (
	"log"

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
	log.Printf("LightPlacer set spriteID to %v", opts.Tile.SpriteID)
	log.Printf("On=%v, Off=%v", placer.OnSpriteID, placer.OffSpriteID)
}

func (placer LightPlacer) ItemSpriteID() render.SpriteID {
	return placer.OnSpriteID
}
