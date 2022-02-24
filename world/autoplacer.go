package world

import (
	"math/rand"

	"github.com/skycoin/cx-game/engine/spriteloader/blobsprites"
	"github.com/skycoin/cx-game/world/tiling"
	"github.com/skycoin/cx-game/world/pipesim"
	"github.com/skycoin/cx-game/render"
)

// place tiles for a given tiletype using an auto-tiling mechanism
type AutoPlacer struct {
    Tile Tile
	blobSpritesIDs []blobsprites.BlobSpritesID
	TilingID     tiling.TilingID
}

func (placer AutoPlacer) sprite(
	neighbours tiling.DetailedNeighbours,
) render.SpriteID {
	blobspritesID :=
		placer.blobSpritesIDs[rand.Intn(len(placer.blobSpritesIDs))]
	sprites := blobsprites.GetBlobSpritesById(blobspritesID)
	idx := tiling.ApplyTiling(placer.TilingID, neighbours)
	return sprites[idx]
}

func (placer AutoPlacer) CreateTile(
	tt TileType, createOpts TileCreationOptions,
) Tile {
	tile := placer.Tile
	updateOpts := TileUpdateOptions{
		Neighbours: createOpts.Neighbours,
		Tile:       &tile,
	}
	placer.UpdateTile(tt, updateOpts)


	return tile
}

func (placer AutoPlacer) UpdateTile(
	tt TileType, opts TileUpdateOptions,
) {
	if opts.Cycling {
		connectedNeighbours :=
			pipesim.ConnectedNeighbours(opts.Tile.Connections, opts.Neighbours)
		opts.Tile.SpriteID = placer.sprite(connectedNeighbours)
	} else {
		opts.Tile.Connections =
			pipesim.ConnectionsFromNeighbours(opts.Neighbours)
		opts.Tile.SpriteID = placer.sprite(opts.Neighbours)
	}
}

func (placer AutoPlacer) ItemSpriteID() render.SpriteID {
	return placer.sprite(tiling.DetailedNeighbours{})
}
