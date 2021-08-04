package world

import (
	"math/rand"

	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
	"github.com/skycoin/cx-game/spriteloader"
)

// place tiles for a given tiletype using an auto-tiling mechanism
type AutoPlacer struct {
	blobSpritesIDs []blobsprites.BlobSpritesID
	TileTypeID TileTypeID
	TilingType blob.TilingType
}

func (placer AutoPlacer) sprite(
		neighbours blob.Neighbours,
) spriteloader.SpriteID {
	blobspritesID :=
		placer.blobSpritesIDs[rand.Intn(len(placer.blobSpritesIDs))]
	sprites := blobsprites.GetBlobSpritesById(blobspritesID)
	idx := blob.ApplyTiling(placer.TilingType, neighbours)
	return sprites[idx]
}

func (placer AutoPlacer) CreateTile(
		tt TileType, createOpts TileCreationOptions,
) Tile {
	tile := Tile{}
	updateOpts := TileUpdateOptions {
		Neighbours: createOpts.Neighbours,
		Tile: &tile,
	}
	placer.UpdateTile(tt,updateOpts)
	return tile
}

func (placer AutoPlacer) UpdateTile(
		tt TileType, opts TileUpdateOptions,
) {
	*opts.Tile = Tile {
		SpriteID: placer.sprite(opts.Neighbours),
		Name: tt.Name,
		TileCategory: TileCategoryNormal,
		TileTypeID: tt.ID,
	}
}

func (placer AutoPlacer) ItemSpriteID() spriteloader.SpriteID {
	return placer.sprite(blob.Neighbours{})
}
