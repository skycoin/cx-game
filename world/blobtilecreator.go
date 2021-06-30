package world

import (
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

type BlobTileCreator struct {
	name string
	blobSpritesId blobsprites.BlobSpritesID
	TileTypeID TileTypeID
}

func NewBlobTileCreator(
		name string, blobSpritesId blobsprites.BlobSpritesID,
) *BlobTileCreator {
	return &BlobTileCreator { name: name, blobSpritesId: blobSpritesId }
}

func (creator BlobTileCreator) blobSprites() []uint32 {
	return blobsprites.GetBlobSpritesById(creator.blobSpritesId)
}

func (creator BlobTileCreator) CreateTile(neighbours blob.Neighbours) Tile {
	tile := Tile{}
	creator.UpdateTile(&tile,neighbours)
	return tile
}

func (creator BlobTileCreator) UpdateTile(
		tile *Tile, neighbours blob.Neighbours,
) {
	blobSpriteIdx := blob.ApplyBlobTiling(neighbours)
	*tile =  Tile {
		SpriteID: creator.blobSprites()[blobSpriteIdx],
		Name: creator.name,
		TileCategory: TileCategoryNormal,
		TileTypeID: creator.TileTypeID,
	}
}
