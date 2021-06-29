package world

import (
	"math/rand"

	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

// TODO shove in .yaml file
const persistence = 0.5
const lacunarity = 2
const xs = 3

func (planet *Planet) placeTileOnTop(x int, tile Tile) int {
	y := planet.GetHeight(x) + 1
	tileIdx := planet.GetTileIndex(x, y)
	planet.Layers.Top[tileIdx] = tile
	return y
}

func (planet *Planet) placeLayer(
	tiles []Tile, depth,noiseScale float32,
) []cxmath.Vec2i {
	positions := []cxmath.Vec2i {}
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), xs, 256)
	for x:=int32(0); x<planet.Width; x++ {
		noiseSample := perlin.Noise(float32(x), 0, persistence, lacunarity, 8)
		height := int(depth+noiseSample*noiseScale)
		for i:=0; i<height; i++ {
			tile := tiles[rand.Intn(len(tiles))]
			y := planet.placeTileOnTop(int(x),tile)
			positions = append(positions, cxmath.Vec2i { x,int32(y) } )
		}
	}
	return positions
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

func (planet *Planet) placeBgTile(tile Tile, pos cxmath.Vec2i) {
	tileIdx := planet.GetTileIndex(int(pos.X),int(pos.Y))
	planet.Layers.Mid[tileIdx] = tile
}


func GeneratePlanet() *Planet {
	planet := NewPlanet(100, 100)
	oreSheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/tile/ores-stone.png",3,4)
	spriteloader.
		LoadSingleSprite("./assets/tile/dirt.png", "Dirt")
	spriteloader.
		LoadSingleSprite("./assets/tile/bedrock.png", "Bedrock")
	spriteloader.
		LoadSingleSprite("./assets/tile/stone.png", "Stone")

	spriteloader.
		LoadSprite(oreSheetId, "Big Gold", 2, 0)
	spriteloader.
		LoadSprite(oreSheetId, "Ruby Pentagon", 2, 2)

	dirtBlobSpritesId :=
		blobsprites.LoadBlobSprites("./assets/tile/Tiles_1.png")
	altDirtBlobSpritesId :=
		blobsprites.LoadBlobSprites("./assets/tile/Tiles_1_v1.png")
	dirtWallBlobSpritesId :=
		blobsprites.LoadBlobSprites("./assets/tile/Wall_1.png")

	dirt := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Dirt")),
		Name: "Dirt",
		IsBlob: true,
		BlobSpriteID: dirtBlobSpritesId,
	}
	altDirt := dirt
	altDirt.BlobSpriteID = altDirtBlobSpritesId
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
		SpriteID:  uint32(spriteloader.GetSpriteIdByName("Big Gold")),
	}
	blueOre := Tile {
		TileType: TileTypeNormal,
		SpriteID: uint32(spriteloader.GetSpriteIdByName("Ruby Pentagon")),
	}

	planet.placeLayer([]Tile{bedrock}, 4,1)
	stonePositions := planet.placeLayer([]Tile{stone}, 8,2)

	// todo make dirt wall look different
	dirtWall := Tile {
		TileType: TileTypeNormal,
		Name: "Dirt Wall",
		IsBlob: true,
		BlobSpriteID: dirtWallBlobSpritesId,
	}
	dirtPositions := planet.placeLayer([]Tile{dirt,altDirt}, 4,1)
	for _,pos := range dirtPositions { planet.placeBgTile(dirtWall, pos) }
	for _,pos := range stonePositions { planet.placeBgTile(dirtWall, pos) }

	planet.placeOres(purpleOre, 0.6)
	planet.placeOres(blueOre, 0.7)
	
	return planet
}
