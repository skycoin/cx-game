package world

import (
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

var TileTypeIDs struct {
	Air TileTypeID
	Dirt TileTypeID
	Stone TileTypeID
	Bedrock TileTypeID
	DirtWall TileTypeID
}

func RegisterTileTypes() {
	RegisterEmptyTileType()
	RegisterDirtTileType()
	RegisterStoneTileType()
	RegisterBedrockTileType()
	RegisterDirtWallTileType()
}

func RegisterEmptyTileType() {
	TileTypeIDs.Air = RegisterTileType(TileType {
		Name: "Air",
		Placer: DirectPlacer{},
	})
}

func RegisterStoneTileType() {
	spriteID :=
		spriteloader.LoadSingleSprite("./assets/tile/stone.png","Stone")
	TileTypeIDs.Stone = RegisterTileType(TileType {
		Name: "Stone",
		Placer: DirectPlacer{SpriteID:uint32(spriteID)},
		Layer: TopLayer,
	})
}

func RegisterBedrockTileType() {
	spriteID :=
		spriteloader.LoadSingleSprite("./assets/tile/bedrock.png","Stone")
	TileTypeIDs.Bedrock = RegisterTileType(TileType {
		Name: "Bedrock",
		Placer: DirectPlacer{SpriteID:uint32(spriteID)},
		Layer: TopLayer,
		Invulnerable: true,
	})
}

func RegisterDirtTileType() {
	blobSpritesId :=
		blobsprites.LoadBlobSprites("./assets/tile/Tiles_1.png")
	TileTypeIDs.Dirt = RegisterTileType(TileType {
		Name: "Dirt",
		Placer: AutoPlacer{blobSpritesId: blobSpritesId},
		Layer: TopLayer,
	})
}

func RegisterDirtWallTileType() {
	blobSpritesId :=
		blobsprites.LoadBlobSprites("./assets/tile/Wall_1.png")
	TileTypeIDs.DirtWall = RegisterTileType(TileType {
		Name: "Dirt Wall",
		Placer: AutoPlacer{blobSpritesId: blobSpritesId},
		Layer: TopLayer,
	})

}
