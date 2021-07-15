package world

import (
	"github.com/skycoin/cx-game/spriteloader"
)

func NewDevPlanet() *Planet {
	// TODO determine dirt height from perlin
	planet := NewPlanet(100, 100)
	planet.WorldState = NewDevWorldState()

	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/tile/mixed-tileset_00.png")
	spriteloader.
		LoadSingleSprite("./assets/tile/dirt.png", "Dirt")
	spriteloader.
		LoadSingleSprite("./assets/tile/stone.png", "Stone")
	spriteloader.
		LoadSingleSprite("./assets/tile/bedrock.png", "Bedrock")
	spriteloader.
		LoadSprite(spriteSheetId, "RedBlip", 0, 0)

	foregroundTilesSpritesheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/tile/ForegroundTiles.png", 16, 16)

	blueLabSpriteIds := spriteloader.
		LoadSprites(foregroundTilesSpritesheetId, "bluelab", 6, 3, 15, 3)

	// dirt
	for x := 0; x < int(planet.Width); x++ {
		for y := 4; y < 6; y++ {
			tileIdx := planet.GetTileIndex(x, y)
			planet.Layers.Top[tileIdx] = Tile{
				TileCategory: TileCategoryNormal,
				SpriteID:     (spriteloader.GetSpriteIdByName("Dirt")),
			}
		}
	}
	// stone
	for x := 0; x < int(planet.Width); x++ {
		for y := 2; y < 4; y++ {
			tileIdx := planet.GetTileIndex(x, y)
			planet.Layers.Top[tileIdx] = Tile{
				TileCategory: TileCategoryNormal,
				SpriteID:     (spriteloader.GetSpriteIdByName("Stone")),
			}
		}
	}
	// bedrock
	for x := 0; x < int(planet.Width); x++ {
		for y := 0; y < 2; y++ {
			tileIdx := planet.GetTileIndex(x, y)
			planet.Layers.Top[tileIdx] = Tile{
				TileCategory: TileCategoryNormal,
				SpriteID:     (spriteloader.GetSpriteIdByName("Bedrock")),
			}
		}
	}

	// DEBUG: tiles to test collision and physics
	*planet.GetTopLayerTile(27, 7) = Tile{
		TileCategory: TileCategoryNormal,
		SpriteID:     (spriteloader.GetSpriteIdByName("Stone")),
	}
	*planet.GetTopLayerTile(25, 5) = Tile{
		TileCategory: TileCategoryNone,
		SpriteID:     (spriteloader.GetSpriteIdByName("Dirt")),
	}

	// wall to test
	for i := 6; i < 35; i++ {
		*planet.GetTopLayerTile(10, i) = Tile{
			TileCategory: TileCategoryNormal,
			SpriteID:     (spriteloader.GetSpriteIdByName(("Stone"))),
		}
		*planet.GetTopLayerTile(5, i) = Tile{
			TileCategory: TileCategoryNormal,
			SpriteID:     (spriteloader.GetSpriteIdByName(("Stone"))),
		}
	}

	devMultiTile := MultiTile{
		Width: 10, Height: 1,
		TileCategory: TileCategoryNormal,
		SpriteIDs:    blueLabSpriteIds,
		Name:         "dev multi tile",
	}
	planet.PlaceMultiTile(20, 8, BgLayer, devMultiTile)

	return planet
}
