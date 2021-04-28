package world

import (
	"github.com/skycoin/cx-game/spriteloader"
)

func NewDevPlanet() *Planet {
	planet := NewPlanet(100,100)

	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/tile/mixed-tileset_00.png")
	spriteloader.
		LoadSprite(spriteSheetId,"red",0,0)

	for x := 0; x < int(planet.Width); x++ {
		tileIdx := planet.GetTileIndex(x,0)
		planet.Layers.Background[tileIdx] = Tile {
			TileType: TileTypeNormal,
			SpriteID: uint32(spriteloader.GetSpriteIdByName("red")),
		}
	}
	return planet
}
