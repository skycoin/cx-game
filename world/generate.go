package world

import (
	"log"

	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/spriteloader"
)

func (planet *Planet) placeTileOnTop(x int, tile Tile) {
	y := planet.GetHeight(x) + 1
	tileIdx := planet.GetTileIndex(x, y)
	planet.Layers.Top[tileIdx] = tile
}

const seed = 2
const persistence = 0.5
const lacunarity = 2

const heightScale = 4

func GeneratePlanet() *Planet {
	planet := NewPlanet(100, 100)
	spriteloader.
		LoadSingleSprite("./assets/tile/dirt.png", "Dirt")
	spriteloader.
		LoadSingleSprite("./assets/tile/bedrock.png", "Bedrock")
	spriteloader.
		LoadSingleSprite("./assets/tile/stone.png", "Stone")
	dirtTile := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
	}

	perlin := perlin.NewPerlin2D(seed, 512, 4, 256)

	for x:=int32(0); x<planet.Width; x++ {
		dirtNoise := perlin.Noise(float32(x), 0, persistence, lacunarity, 8)
		log.Printf("dirt noise = %v",dirtNoise)
		dirtHeight := int((dirtNoise+1) * heightScale)
		for i:=0; i<dirtHeight; i++ {
			planet.placeTileOnTop(int(x),dirtTile)
		}
	}
	
	return planet
}
