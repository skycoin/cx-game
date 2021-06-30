package ui

import (
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/spriteloader"
)

func NewDevTilePaleteSelector() TilePaletteSelector {
	selector := MakeTilePaleteSelector(11,11)

	/*
	// solid tiles
	selector.AddTile(
		world.Tile {
			TileCategory: world.TileCategoryNormal,
			SpriteID: uint32(spriteloader.GetSpriteIdByName("Bedrock")),
		},
		0,1, world.TopLayer,
	)
	selector.AddTile(
		world.Tile {
			TileCategory: world.TileCategoryNormal,
			SpriteID: uint32(spriteloader.GetSpriteIdByName("Stone")),
		},
		0,2, world.TopLayer,
	)
	selector.AddTile(
		world.Tile {
			TileCategory: world.TileCategoryNormal,
			SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
		},
		0,3, world.TopLayer,
	)
	*/

	// mid-layer objects
	// TODO: share texture resources between world and UI
	foregroundTilesSpritesheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/tile/ForegroundTiles.png",16,16)
	blueLabSpriteIds := spriteloader.
		LoadSprites(foregroundTilesSpritesheetId,"bluelab",6,3,15,3)
	selector.AddMultiTile(
		world.MultiTile {
			Width: 10, Height: 1, TileCategory: world.TileCategoryNormal,
			SpriteIDs: blueLabSpriteIds,
			Name: "dev multi tile",
		},
		1, 1, world.BgLayer,
	)

	return selector
}
