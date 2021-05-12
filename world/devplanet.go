package world

import (
	"github.com/skycoin/cx-game/spriteloader"
)

func NewDevPlanet() *Planet {
	// TODO determine dirt height from perlin
	planet := NewPlanet(100, 100)

	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/tile/mixed-tileset_00.png")
	spriteloader.
		LoadSingleSprite("./assets/tile/dirt.png","Dirt")
	spriteloader.
		LoadSingleSprite("./assets/tile/stone.png", "Stone")
	spriteloader.
		LoadSingleSprite("./assets/tile/bedrock.png", "Bedrock")
	spriteloader.
		LoadSprite(spriteSheetId, "RedBlip", 0, 0)

	// dirt
	for x := 0; x < int(planet.Width); x++ {
		for y := 4; y < 6; y++ {
			tileIdx := planet.GetTileIndex(x, y)
			planet.Layers.Top[tileIdx] = Tile{
				TileType: TileTypeNormal,
				SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
			}
		}
	}
	// stone
	for x := 0; x < int(planet.Width); x++ {
		for y := 2; y < 4; y++ {
			tileIdx := planet.GetTileIndex(x, y)
			planet.Layers.Top[tileIdx] = Tile{
				TileType: TileTypeNormal,
				SpriteID: uint32(spriteloader.GetSpriteIdByName("Stone")),
			}
		}
	}
	// bedrock
	for x := 0; x < int(planet.Width); x++ {
		for y := 0; y < 2; y++ {
			tileIdx := planet.GetTileIndex(x, y)
			planet.Layers.Top[tileIdx] = Tile{
				TileType: TileTypeNormal,
				SpriteID: uint32(spriteloader.GetSpriteIdByName("Bedrock")),
			}
		}
	}

	// DEBUG: tiles to test collision and physics
	*planet.GetTopLayerTile(27, 7) = Tile{
		TileType: TileTypeNormal,
		SpriteID: 1,
	}
	*planet.GetTopLayerTile(25, 5) = Tile{
		TileType: TileTypeNone,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
	}

	return planet
}
