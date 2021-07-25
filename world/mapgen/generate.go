package mapgen

import (
	"math/rand"

	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

// TODO shove in .yaml file
const persistence = 0.5
const lacunarity = 2
const xs = 3

func placeTileOnTop(planet *world.Planet, x int, tile world.Tile) int {
	y := planet.GetHeight(x) + 1
	tileIdx := planet.GetTileIndex(x, y)
	planet.Layers.Top[tileIdx] = tile
	return y
}

func placeLayer(
	planet *world.Planet, tileTypeID world.TileTypeID,
	depth,noiseScale float32,
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
func placeOres(planet *world.Planet,tile world.Tile, threshold float32) {
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

func placeBgTile(
		planet *world.Planet, tileTypeID world.TileTypeID, pos cxmath.Vec2i,
) {
	planet.PlaceTileType(tileTypeID, int(pos.X),int(pos.Y))
}

func placePoles(planet *world.Planet) {
	w := float32(planet.Width)
	north := int( w*(1.0/2.0 + 1.0/4.0) )
	south := int( w*(1.0/2.0 - 1.0/4.0) )
	placePole(planet,north)
	placePole(planet,south)
}

const poleRadius int = 4
func placePole(planet *world.Planet, origin int) {
	for x := origin - poleRadius ; x < origin + poleRadius ; x++ {
		for y := 0 ; y < int(planet.Height) ; y++ {
			tile := planet.GetTile(x,y,world.TopLayer)
			if tile.TileTypeID == world.TileTypeIDs.Dirt {
				planet.PlaceTileType(world.TileTypeIDs.MethaneIce, x,y)
			}
		}
	}
}


func GeneratePlanet() *world.Planet {
	planet := world.NewPlanet(100, 100)
	placeLayer(planet,world.TileTypeIDs.Air, 4,1)
	stonePositions := placeLayer(planet,world.TileTypeIDs.Stone, 8,2)
	dirtPositions := placeLayer(planet,world.TileTypeIDs.Dirt, 4,1)

	for _,pos := range dirtPositions {
		placeBgTile(planet,world.TileTypeIDs.DirtWall, pos)
	}
	for _,pos := range stonePositions {
		placeBgTile(planet,world.TileTypeIDs.DirtWall, pos)
	}

	placePoles(planet)

	return planet
}
