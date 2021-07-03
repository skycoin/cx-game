package world

import (
	"math/rand"

	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/cxmath"
	//"github.com/skycoin/cx-game/spriteloader"
	//"github.com/skycoin/cx-game/spriteloader/blobsprites"
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
	tileTypeID TileTypeID, depth,noiseScale float32,
) []cxmath.Vec2i {
	positions := []cxmath.Vec2i {}
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), xs, 256)
	for x:=int32(0); x<planet.Width; x++ {
		noiseSample := perlin.Noise(float32(x), 0, persistence, lacunarity, 8)
		height := int(depth+noiseSample*noiseScale)
		for i:=0; i<height; i++ {
			y := planet.GetHeight(int(x)) + 1
			planet.PlaceTileType(tileTypeID,int(x),int(y))
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

func (planet *Planet) placeBgTile(tileTypeID TileTypeID, pos cxmath.Vec2i) {
	planet.PlaceTileType(tileTypeID, int(pos.X),int(pos.Y))
}


func GeneratePlanet() *Planet {
	planet := NewPlanet(100, 100)
	planet.placeLayer(TileTypeIDs.Air, 4,1)
	stonePositions := planet.placeLayer(TileTypeIDs.Stone, 8,2)
	dirtPositions := planet.placeLayer(TileTypeIDs.Dirt, 4,1)

	for _,pos := range dirtPositions {
		planet.placeBgTile(TileTypeIDs.DirtWall, pos)
	}
	for _,pos := range stonePositions {
		planet.placeBgTile(TileTypeIDs.DirtWall, pos)
	}

	return planet
}
