package blobsprites

import (
	"fmt"

	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader"
)

var allBlobSprites = make(map[uint32]([]uint32))
var nextBlobSpriteId = uint32(1)

func LoadBlobSprites(fname string) uint32 {
	spritesheetId := spriteloader.LoadSpriteSheetByColRow(
		fname, blob.BlobSheetHeight, blob.BlobSheetWidth )
	blobSprites := []uint32{}
	for idx:=0; idx < blob.BlobSheetWidth*blob.BlobSheetHeight; idx++ {
		y := idx / blob.BlobSheetWidth
		x := idx % blob.BlobSheetWidth
		name := fmt.Sprint("blob_%d",idx)
		blobSprites =
			append(blobSprites,spriteloader.LoadSprite(spritesheetId,name,x,y))
	}
	blobSpriteId := nextBlobSpriteId
	allBlobSprites[blobSpriteId] = blobSprites
	nextBlobSpriteId += 1
	return blobSpriteId
}

func GetBlobSpritesById(blobSpriteID uint32) []uint32 {
	return allBlobSprites[blobSpriteID]
}
