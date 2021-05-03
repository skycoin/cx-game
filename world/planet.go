package world

import (
	"math"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
)

type Layers struct {
	Background []Tile
	Mid []Tile
	Top []Tile
}

type Planet struct {
	Width  int32
	Height int32
	Layers Layers
}

func NewPlanet(x, y int32) *Planet {
	planet := Planet{
		Width:	x,
		Height: y,
		Layers:  Layers{
			Background: make([]Tile,x*y),
			Mid: make([]Tile,x*y),
			Top: make([]Tile,x*y),
		},
	}

	return &planet
}

func (planet *Planet) DrawLayer(tiles []Tile, cam *camera.Camera) {
	for idx,tile := range tiles {
		y := int32(idx)/planet.Width
		x := int32(idx)%planet.Width

		if (tile.TileType != TileTypeNone) {
			spriteloader.DrawSpriteQuad(
				float32(x)-cam.X,float32(y)-cam.Y,
				1,1,
				int(tile.SpriteID),
			)
		}
	}
}

func (planet *Planet) Draw(cam *camera.Camera) {
	planet.DrawLayer(planet.Layers.Background, cam)
	planet.DrawLayer(planet.Layers.Mid, cam)
	planet.DrawLayer(planet.Layers.Top, cam)
}

func (planet *Planet) GetTileIndex(x,y int) int {
	return y*int(planet.Width)+x
}

func (planet *Planet) GetAllTilesUnique() []Tile {
	tiles := []Tile {}
	seenTiles := map[Tile]bool {}
	for _,tile := range planet.Layers.Top {
		_,seen := seenTiles[tile]
		if !seen {
			tiles = append(tiles,tile)
		}
		seenTiles[tile] = true
	}
	return tiles
}

func (planet *Planet) TryPlaceTile(
	x,y float32, projection mgl32.Mat4,
	tile Tile,
	cam *camera.Camera,
) {
	// tilemap is drawn directly on the world - no need to convert further
	worldCoords := cxmath.ConvertScreenCoordsToWorld(x,y,projection)
	// FIXME dirty workaround for broken view matrx
	worldCoords[0] += cam.X
	worldCoords[1] += cam.Y
	tileX := int32(math.Floor((float64(worldCoords.X()+0.5))))
	tileY := int32(math.Floor((float64)(worldCoords.Y()+0.5)))
	if tileX>=0 && tileX<planet.Width && tileY>=0 && tileY<planet.Width {
		tileIdx := planet.GetTileIndex(int(tileX),int(tileY))
		// TODO allow placing on background and mid layers
		planet.Layers.Top[tileIdx] = tile
	}

// gets the y coordinate of the highest solid tile
// for a given column given by an x coordinate
func (planet *Planet) GetHeight(x int) int {
    for y:=int(planet.Height-1); y>=0; y-- {
        idx := planet.GetTileIndex(x,y)
        if planet.Layers.Top[idx].TileType != TileTypeNone {
            return y
        }
    }
    return 0
}
