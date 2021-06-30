package world

import (
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

var (
	emptyTileTypeID TileTypeID
	dirtTileTypeID TileTypeID
	stoneTileTypeID TileTypeID
	bedrockTileTypeID TileTypeID
)

func RegisterTileTypes() {
	RegisterEmptyTileType()
	RegisterDirtTileType()
	RegisterStoneTileType()
	RegisterBedrockTileType()
}

func RegisterEmptyTileType() {
	emptyTileTypeID = RegisterTileType(TileType {
		Name: "Air",
		CreateTile: func (n blob.Neighbours) Tile {
			return Tile {}
		},
		UpdateTile: func (t *Tile,n blob.Neighbours) {},
	})
}

func RegisterSimpleTileType(name string, spriteID uint32) TileTypeID {
	id := NextTileTypeID()
	return RegisterTileType(TileType {
		Name: name,
		Layer: TopLayer,
		CreateTile: func (n blob.Neighbours) Tile { return Tile {
			SpriteID: spriteID,
			Name: name,
			TileCategory: TileCategoryNormal,
			TileTypeID: id,
			Durability: 1,
		} },
		UpdateTile: func (t *Tile,n blob.Neighbours) {},
	})
}

func RegisterStoneTileType() {
	spriteID :=
		spriteloader.LoadSingleSprite("./assets/tile/stone.png","Stone")
	stoneTileTypeID = RegisterSimpleTileType("Stone",uint32(spriteID))
}

func RegisterBedrockTileType() {
	spriteID :=
		spriteloader.LoadSingleSprite("./assets/tile/bedrock.png","Stone")
	bedrockTileTypeID = RegisterSimpleTileType("Bedrock",uint32(spriteID))
	tileTypes[bedrockTileTypeID].Invulnerable = true
}

func RegisterDirtTileType() {
	blobSpritesId :=
		blobsprites.LoadBlobSprites("./assets/tile/Tiles_1.png")
	creator := NewBlobTileCreator("Dirt", blobSpritesId)
	creator.TileTypeID = NextTileTypeID()
	tileType := NewTileType("Dirt", TopLayer, creator.CreateTile)
	tileType.UpdateTile = creator.UpdateTile
	dirtTileTypeID = RegisterTileType(tileType)
}
