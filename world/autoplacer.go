package world

import (
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

// place tiles for a given tiletype using an auto-tiling mechanism
type AutoPlacer struct {
	name string
	blobSpritesId blobsprites.BlobSpritesID
	TileTypeID TileTypeID
}

func (ap AutoPlacer) blobSprites() []uint32 {
	return blobsprites.GetBlobSpritesById(ap.blobSpritesId)
}

func (ap AutoPlacer) CreateTile(
		tt TileType, createOpts TileCreationOptions,
) Tile {
	tile := Tile{}
	updateOpts := TileUpdateOptions {
		Neighbours: createOpts.Neighbours,
		Tile: &tile,
	}
	ap.UpdateTile(tt,updateOpts)
	return tile
}

func (ap AutoPlacer) UpdateTile(
		tt TileType, opts TileUpdateOptions,
) {
	blobSpriteIdx := blob.ApplyBlobTiling(opts.Neighbours)
	*opts.Tile = Tile {
		SpriteID: ap.blobSprites()[blobSpriteIdx],
		Name: tt.Name,
		TileCategory: TileCategoryNormal,
		TileTypeID: tt.ID,
	}
}
