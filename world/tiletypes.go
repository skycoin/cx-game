package world

import (
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

var TileTypeIDs struct {
	Air TileTypeID
	Dirt TileTypeID
	Stone TileTypeID
	Bedrock TileTypeID
	DirtWall TileTypeID
	Pipe TileTypeID
}

func RegisterTileTypes() {
	RegisterEmptyTileType()
	RegisterDirtTileType()
	RegisterStoneTileType()
	RegisterBedrockTileType()
	RegisterDirtWallTileType()
	RegisterPipeTileType()
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
		Placer: DirectPlacer{SpriteID:(spriteID)},
		Layer: TopLayer,
	})
}

func RegisterBedrockTileType() {
	spriteID :=
		spriteloader.LoadSingleSprite("./assets/tile/bedrock.png","Stone")
	TileTypeIDs.Bedrock = RegisterTileType(TileType {
		Name: "Bedrock",
		Placer: DirectPlacer{SpriteID:(spriteID)},
		Layer: TopLayer,
		Invulnerable: true,
	})
}

func RegisterDirtTileType() {
	blobSpritesId :=
		blobsprites.LoadFullBlobSprites("./assets/tile/Tiles_1.png")
	TileTypeIDs.Dirt = RegisterTileType(TileType {
		Name: "Dirt",
		Placer: AutoPlacer{
			blobSpritesId: blobSpritesId, TilingType: blob.FullBlobTiling,
		},
		Layer: TopLayer,
	})
}

func RegisterDirtWallTileType() {
	blobSpritesId :=
		blobsprites.LoadFullBlobSprites("./assets/tile/Wall_1.png")
	TileTypeIDs.DirtWall = RegisterTileType(TileType {
		Name: "Dirt Wall",
		Placer: AutoPlacer{
			blobSpritesId: blobSpritesId, TilingType: blob.FullBlobTiling,
		},
		Layer: BgLayer,
	})
}

func RegisterPipeTileType() {
	blobSpritesId := blobsprites.
		LoadSimpleBlobSprites("./assets/tile/VentilationPipes_1.png")
	TileTypeIDs.Pipe = RegisterTileType(TileType {
		Name: "Pipe",
		Placer: AutoPlacer{
			blobSpritesId: blobSpritesId, TilingType: blob.SimpleBlobTiling,
		},
		Layer: MidLayer,
	})
}
