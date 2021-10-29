package mapgen

import (
	"log"
	"math/rand"

	"github.com/skycoin/cx-game/cxmath"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/world"
)

// TODO shove in .yaml file
const (
	persistence = 0.5
	lacunarity  = 2
	xs          = 3
)

func placeTileOnTop(planet *world.Planet, x int, tile world.Tile) int {
	y := planet.GetHeight(x) + 1
	tileIdx := planet.GetTileIndex(x, y)
	planet.Layers[world.TopLayer].Tiles[tileIdx] = tile
	return y
}

func placeLayer(
	planet *world.Planet, tileTypeID world.TileTypeID,
	depth, noiseScale float32,
) []cxmath.Vec2i {
	positions := []cxmath.Vec2i{}
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), xs, 256)
	for x := int32(0); x < planet.Width; x++ {
		noiseSample := perlin.Noise(float32(x), 0, persistence, lacunarity, 8)
		height := int(depth + noiseSample*noiseScale)
		for i := 0; i < height; i++ {
			y := planet.GetHeight(int(x)) + 1
			planet.PlaceTileType(tileTypeID, int(x), int(y))
			positions = append(positions, cxmath.Vec2i{x, int32(y)})
		}
	}
	return positions
}

const oreLacunarity = lacunarity * 50

func placeOres(planet *world.Planet, tile world.Tile, threshold float32) {
	perlin := perlin.NewPerlin2D(rand.Int63(), int(planet.Width), xs, 256)
	for y := int32(0); y < planet.Height; y++ {
		for x := int32(0); x < planet.Width; x++ {
			xf := float32(x)
			yf := float32(y)
			sample := perlin.Noise(xf, yf, persistence, oreLacunarity, 8)
			tileIdx := planet.GetTileIndex(int(x), int(y))
			if sample > threshold &&
				planet.Layers[world.TopLayer].Tiles[tileIdx].Name == "Stone" {
				planet.Layers[world.TopLayer].Tiles[tileIdx] = tile
			}
		}
	}
}

func placeBgTile(
	planet *world.Planet, tileTypeID world.TileTypeID, pos cxmath.Vec2i,
) {
	planet.PlaceTileType(tileTypeID, int(pos.X), int(pos.Y))
}

func placePoles(planet *world.Planet) {
	w := float32(planet.Width)
	north := int(w * (1.0/2.0 + 1.0/4.0))
	south := int(w * (1.0/2.0 - 1.0/4.0))
	placePole(planet, north)
	placePole(planet, south)
}

func idFor(name string) world.TileTypeID {
	id, ok := world.IDFor(name)
	if !ok {
		log.Fatalf("cannot find tile type ID for \"%s\"", name)
	}
	return id
}

const poleRadius int = 4

func placePole(planet *world.Planet, origin int) {
	for x := origin - poleRadius; x < origin+poleRadius; x++ {
		for y := 0; y < int(planet.Height); y++ {
			tile := planet.GetTile(x, y, world.TopLayer)
			if tile.TileTypeID == idFor("regolith") {
				planet.PlaceTileType(idFor("methane-ice"), x, y)
			}
		}
	}
}

func GeneratePlanet() *world.Planet {
	planet := world.NewPlanet(100, 100)
	placeLayer(planet, idFor("air"), 4, 1)
	stonePositions := placeLayer(planet, idFor("stone"), 8, 2)
	dirtPositions := placeLayer(planet, idFor("regolith"), 4, 1)

	for _, pos := range dirtPositions {
		placeBgTile(planet, idFor("regolith-wall"), pos)
	}
	for _, pos := range stonePositions {
		placeBgTile(planet, idFor("regolith-wall"), pos)
	}

	placePoles(planet)

	return planet
}
