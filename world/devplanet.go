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
	planet.Layers.Background[0] = Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("red")),
	}
	return planet
}
