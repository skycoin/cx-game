package blobsprites

import (
	"fmt"

	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader"
)

type BlobSpritesID uint32
var allBlobSprites = make(map[BlobSpritesID]([]spriteloader.SpriteID))
var nextBlobSpriteId = BlobSpritesID(1)

func LoadBlobSprites(fname string, w,h int) BlobSpritesID {
	spritesheetId := spriteloader.LoadSpriteSheetByColRow(
		fname, h, w )
	blobSprites := []spriteloader.SpriteID{}
	for idx:=0; idx < w*h; idx++ {
		y := idx / w
		x := idx % w
		name := fmt.Sprint("blob_%d",idx)
		blobSprites =
			append(blobSprites,spriteloader.LoadSprite(spritesheetId,name,x,y))
	}
	blobSpriteId := nextBlobSpriteId
	allBlobSprites[blobSpriteId] = blobSprites
	nextBlobSpriteId += 1
	return blobSpriteId
}

func LoadFullBlobSprites(fname string) BlobSpritesID {
	return LoadBlobSprites(fname, blob.BlobSheetWidth, blob.BlobSheetHeight)
}

func LoadSimpleBlobSprites(fname string) BlobSpritesID {
	return LoadBlobSprites(
		fname,
		blob.SimpleBlobSheetWidth, blob.SimpleBlobSheetHeight,
	)
}

func GetBlobSpritesById(id BlobSpritesID) []spriteloader.SpriteID {
	return allBlobSprites[id]
}
