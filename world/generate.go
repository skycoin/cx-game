package world

import (
	"math/rand"

	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/spriteloader"
)

// TODO shove in .yaml file
const persistence = 0.5
const lacunarity = 2
const xs = 3

func (planet *Planet) placeTileOnTop(x int, tile Tile) {
	y := planet.GetHeight(x) + 1
	tileIdx := planet.GetTileIndex(x, y)
	planet.Layers.Top[tileIdx] = tile
}

func (planet *Planet) placeLayer(tile Tile, depth,noiseScale float32) {
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), xs, 256)
	for x:=int32(0); x<planet.Width; x++ {
		noiseSample := perlin.Noise(float32(x), 0, persistence, lacunarity, 8)
		height := int(depth+noiseSample*noiseScale)
		for i:=0; i<height; i++ {
			planet.placeTileOnTop(int(x),tile)
		}
	}
}

const oreLacunarity = lacunarity*50
func (planet *Planet) placeOres(tile Tile, threshold float32) {
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), xs, 256)
	for y:=int32(0); y<planet.Height; y++ {
		for x:=int32(0); x<planet.Width; x++ {
			xf := float32(x)
			yf := float32(y)
			sample := perlin.Noise(xf, yf, persistence, oreLacunarity, 8)
			tileIdx := planet.GetTileIndex(int(x),int(y))
			if sample > threshold && planet.Layers.Top[tileIdx].Name=="Stone" {
				planet.Layers.Top[tileIdx] = tile
			}
		}
	}
}


func GeneratePlanet() *Planet {
	planet := NewPlanet(100, 100)
	oreSheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/tile/gems-ores-stone.png",8,8)
	spriteloader.
		LoadSingleSprite("./assets/tile/dirt.png", "Dirt")
	spriteloader.
		LoadSingleSprite("./assets/tile/bedrock.png", "Bedrock")
	spriteloader.
		LoadSingleSprite("./assets/tile/stone.png", "Stone")

	spriteloader.
		LoadSprite(oreSheetId, "Purple Ore", 1, 2)
	spriteloader.
		LoadSprite(oreSheetId, "Blue Ore", 3, 7)

	dirt := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
		Name: "Dirt",
	}
	stone := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Stone")),
		Name: "Stone",
	}
	bedrock := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Bedrock")),
		Name: "Bedrock",
	}
	purpleOre := Tile {
		TileType: TileTypeNormal,
		SpriteID:  uint32(spriteloader.GetSpriteIdByName("Purple Ore")),
	}
	blueOre := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Blue Ore")),
	}

	planet.placeLayer(bedrock, 4,1)
	planet.placeLayer(stone, 8,2)
	planet.placeLayer(dirt, 4,1)

	planet.placeOres(purpleOre, 0.6)
	planet.placeOres(blueOre, 0.7)
	
	return planet
}
