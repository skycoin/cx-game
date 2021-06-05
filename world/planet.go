package world

import (
	"log"

	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	//"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

type Layer int

const (
	BgLayer  Layer = 0
	MidLayer Layer = 1
	TopLayer Layer = 2
)

type Layers struct {
	Background []Tile
	Mid        []Tile
	Top        []Tile
}

type Planet struct {
	Width           int32
	Height          int32
	Layers          Layers
	collidingLines  []float32
	collidingLinesX int
	collidingLinesY int
}

func NewPlanet(x, y int32) *Planet {
	planet := Planet{
		Width:  x,
		Height: y,
		Layers: Layers{
			Background: make([]Tile, x*y),
			Mid:        make([]Tile, x*y),
			Top:        make([]Tile, x*y),
		},
	}

	return &planet
}

func (planet *Planet) DrawLayer(tiles []Tile, cam *camera.Camera) {
	for idx, tile := range tiles {
		y := int32(idx) / planet.Width
		x := int32(idx) % planet.Width

		if tile.TileType.ShouldRender() && cam.IsInBounds(int(x),int(y)) {
			spriteloader.DrawSpriteQuad(
				float32(x)-cam.X, float32(y)-cam.Y,
				1, 1,
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

func (planet *Planet) GetTileIndex(x, y int) int {
	if x >= int(planet.Width) || x < 0 || y >= int(planet.Height) || y < 0 {
		log.Panicln("trying to get tile out of the defined tile array")
	}
	return y*int(planet.Width) + x
}

func (planet *Planet) GetTile(x,y int, layer Layer) *Tile {
	return &planet.GetLayer(layer)[planet.GetTileIndex(x,y)]
}

func (planet *Planet) GetAllTilesUnique() []Tile {
	tiles := []Tile{}
	seenTiles := map[Tile]bool{}
	for _, tile := range planet.Layers.Top {
		_, seen := seenTiles[tile]
		if !seen {
			tiles = append(tiles, tile)
		}
		seenTiles[tile] = true
	}
	return tiles
}

func (planet *Planet) TryPlaceTile(
	x, y float32,
	layer Layer,
	tile Tile,
	cam *camera.Camera,
) bool {
	// click relative to camera
	camCoords := mgl32.Vec4{x / render.PixelsPerTile, y / render.PixelsPerTile, 0, 1}
	// click relative to world
	worldCoords := cam.GetTransform().Mul4x1(camCoords)
	tileX := int32(math.Round((float64(worldCoords.X()))))
	tileY := int32(math.Round((float64(worldCoords.Y()))))
	if tileX >= 0 && tileX < planet.Width && tileY >= 0 && tileY < planet.Width {
		tileIdx := planet.GetTileIndex(int(tileX), int(tileY))
		planetLayer := planet.GetLayer(layer)
		if len(planetLayer) == 0 {
			return false
		}
		if planetLayer[tileIdx].TileType == TileTypeChild ||
		 	planetLayer[tileIdx].TileType == TileTypeMulti {
			planet.RemoveParentTile(planetLayer,tileIdx)
		}
		planetLayer[tileIdx] = tile
		return true
	}
	return false
}

func (planet *Planet) TryPlaceMultiTile(
		x,y float32, layer Layer, multiTile MultiTile, cam *camera.Camera,
) bool {
	// click relative to camera
	camCoords := mgl32.Vec4{x / render.PixelsPerTile, y / render.PixelsPerTile, 0, 1}
	// click relative to world
	worldCoords := cam.GetTransform().Mul4x1(camCoords)
	tileX := int32(math.Round((float64(worldCoords.X()))))
	tileY := int32(math.Round((float64(worldCoords.Y()))))
	if tileX >= 0 && tileX < planet.Width && tileY >= 0 && tileY < planet.Width {
		tileIdx := planet.GetTileIndex(int(tileX), int(tileY))
		planetLayer := planet.GetLayer(layer)
		if len(planetLayer) == 0 {
			return false
		}
		if planetLayer[tileIdx].TileType == TileTypeChild ||
		    planetLayer[tileIdx].TileType == TileTypeMulti {
			planet.RemoveParentTile(planetLayer,tileIdx)
		}
		planet.PlaceMultiTile(int(tileX),int(tileY),layer,multiTile)
		return true
	}
	return false
}

// note that multi-tiles are assumed to be rectangular
func (planet *Planet) getMultiTileWidth(layer []Tile, x,y int) int {
	offsetX := 1
	for layer[planet.GetTileIndex(x+offsetX,y)].OffsetX == int8(offsetX) {
		offsetX++
	}
	return offsetX
}

func (planet *Planet) getMultiTileHeight(layer []Tile, x,y int) int {
	offsetY := 1
	for layer[planet.GetTileIndex(x,y+offsetY)].OffsetY == int8(offsetY) {
		offsetY++
	}
	return offsetY
}

func (planet *Planet) RemoveParentTile(layer []Tile, idx int) {
	tile := layer[idx]
	parentY := idx / int(planet.Width) - int(tile.OffsetY)
	parentX := idx % int(planet.Width) - int(tile.OffsetX)

	width := planet.getMultiTileWidth(layer,parentX,parentY)
	height := planet.getMultiTileHeight(layer,parentX,parentY)

	for y:=parentY; y<parentY+height; y++ {
		for x:=parentX; x<parentX+width; x++ {

			layer[planet.GetTileIndex(x,y)] = Tile {
				TileType: TileTypeNone,
			}

		}
	}

}


func (planet *Planet) GetLayer(layer Layer) []Tile {
	switch layer {
		case BgLayer: return planet.Layers.Background;
		case MidLayer: return planet.Layers.Mid;
		case TopLayer: return planet.Layers.Top;
	}
	return []Tile{}
}

func (planet *Planet) PlaceMultiTile(
	left,bottom int, layer Layer, mt MultiTile,
) {
	planetLayer := planet.GetLayer(layer)

	// place master tile
	planetLayer[planet.GetTileIndex(left,bottom)] = Tile {
		SpriteID: mt.SpriteIDs[0],
		TileType: TileTypeMulti,
		Name: mt.Name,
		// (0,0) offset indicates master / standalone tile
		OffsetX: 0, OffsetY: 0,
	}

	for spriteIdIdx:=1; spriteIdIdx<len(mt.SpriteIDs); spriteIdIdx++ {
		localY := spriteIdIdx / mt.Width
		localX := spriteIdIdx % mt.Width

		x := left + localX
		y := bottom + localY
		tileIdx := planet.GetTileIndex(x,y)

		planetLayer[tileIdx] = Tile {
			SpriteID: mt.SpriteIDs[spriteIdIdx],
			TileType: TileTypeChild,
			Name:     mt.Name,
			OffsetX: int8(localX), OffsetY: int8(localY),
		}
	}
}

// gets the y coordinate of the highest solid tile
// for a given column given by an x coordinate
func (planet *Planet) GetHeight(x int) int {
	for y := int(planet.Height - 1); y >= 0; y-- {
		idx := planet.GetTileIndex(x, y)
		if planet.Layers.Top[idx].TileType != TileTypeNone {
			return y
		}
	}
	return 0
}

func (planet *Planet) GetTopLayerTile(x, y int) *Tile {
	index := planet.GetTileIndex(x, y)
	return &planet.Layers.Top[index]
}

func (planet *Planet) GetCollidingTilesLinesRelative(x, y int) []float32 {
	// avoid from getting outside the world
	if x < 2 {
		x = 2
	}
	if x >= int(planet.Width)-3 {
		x = int(planet.Width) - 4
	}
	if y < 2 {
		y = 0
	}
	if y >= int(planet.Height)-3 {
		y = int(planet.Height) - 4
	}

	// return stored lines if it's in the range
	if x > planet.collidingLinesX-2 && x < planet.collidingLinesX+2 &&
		y > planet.collidingLinesY-2 && y < planet.collidingLinesY+2 {
		return planet.collidingLines
	}

	// store the new position for the newly calculated lines
	planet.collidingLinesX = x
	planet.collidingLinesY = y

	lines := []float32{}

	// calcule all the lines
	for j := y - 2; j < y+4; j++ {
		for i := x - 2; i < x+4; i++ {
			if planet.GetTopLayerTile(i, j).TileType != TileTypeNone {
				fx := float32(i) - 0.5
				fy := float32(j) - 0.5
				fxw := fx + 1.0
				fyh := fy + 1.0

				// only draw the tiles outline instead of every single one
				if planet.GetTopLayerTile(i+1, j).TileType == TileTypeNone { // right
					lines = append(lines, []float32{
						fxw, fyh, 0.0,
						fxw, fy, 0.0,
					}...)
				}
				if planet.GetTopLayerTile(i-1, j).TileType == TileTypeNone { // left
					lines = append(lines, []float32{
						fx, fyh, 0.0,
						fx, fy, 0.0,
					}...)
				}
				if planet.GetTopLayerTile(i, j+1).TileType == TileTypeNone { // up
					lines = append(lines, []float32{
						fx, fyh, 0.0,
						fxw, fyh, 0.0,
					}...)
				}
				if planet.GetTopLayerTile(i, j-1).TileType == TileTypeNone { // down
					lines = append(lines, []float32{
						fx, fy, 0.0,
						fxw, fy, 0.0,
					}...)
				}
			}
		}
	}

	// save array
	planet.collidingLines = lines
	return planet.collidingLines
}

func (planet *Planet) DamageTile(x,y int, layer Layer) {
	tileIdx := planet.GetTileIndex(x,y)
	tile := &planet.GetLayer(layer)[tileIdx]
	tile.Durability--
	if tile.Durability <= 0 {
		*tile = NewEmptyTile()
	}
}
