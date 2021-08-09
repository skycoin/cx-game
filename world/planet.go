package world

import (
	"log"
	"math"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/render/blob"
)

const NUM_INSTANCES = 100

type LayerID int

const (
	BgLayer LayerID = iota
	MidLayer
	TopLayer

	NumLayers // DO NOT SET MANUALLY
)

func (id LayerID) Valid() bool {
	return id >= 0 && id < NumLayers
}

type Layer struct {
	Tiles       []Tile
	Spritesheet spriteloader.Spritesheet
}

type Layers [NumLayers]Layer

func NewLayer(numTiles int32) Layer {
	// TODO pack spritesheet
	return Layer{
		Tiles: make([]Tile, numTiles),
	}
}

func NewLayers(numTiles int32) Layers {
	layers := Layers{}
	for i := LayerID(0); i < NumLayers; i++ {
		layers[i] = NewLayer(numTiles)
	}
	return layers
}

type Planet struct {
	//WorldState      *WorldState
	Width           int32
	Height          int32
	Layers          Layers
	collidingLines  []float32
	collidingLinesX int
	collidingLinesY int
	Time            float32

	program, liquidProgram render.Program
}

func newPlanetProgram() render.Program {
	config := render.NewShaderConfig(
		"./assets/shader/world.vert", "./assets/shader/sprite.frag")
	config.Define("NUM_INSTANCES", strconv.Itoa(NUM_INSTANCES))
	program := config.Compile()

	program.Use()
	defer program.StopUsing()
	// TODO only need one of these
	program.SetVec4F("color", 1, 1, 1, 1)
	program.SetVec4F("colour", 1, 0, 0, 1)
	return program
}

func newPlanetLiquidProgram() render.Program {
	config := render.NewShaderConfig(
		"./assets/shader/world.vert", "./assets/shader/liquid.frag")
	config.Define("NUM_INSTANCES", strconv.Itoa(NUM_INSTANCES))
	program := config.Compile()

	return program
}

func NewPlanet(x, y int32) *Planet {
	planet := Planet{
		//WorldState: NewDevWorldState(),
		Width:         x,
		Height:        y,
		Layers:        NewLayers(x * y),
		program:       newPlanetProgram(),
		liquidProgram: newPlanetLiquidProgram(),
	}
	return &planet
}

func (planet *Planet) GetNeighbours(layer []Tile, x, y int) blob.Neighbours {
	return blob.Neighbours{
		Up:        planet.TileExists(layer, x, y+1),
		UpRight:   planet.TileExists(layer, x+1, y+1),
		Right:     planet.TileExists(layer, x+1, y),
		DownRight: planet.TileExists(layer, x+1, y-1),
		Down:      planet.TileExists(layer, x, y-1),
		DownLeft:  planet.TileExists(layer, x-1, y-1),
		Left:      planet.TileExists(layer, x-1, y),
		UpLeft:    planet.TileExists(layer, x-1, y+1),
	}
}

func (planet *Planet) GetTileIndex(x, y int) int {
	// apply wrap-around
	x = cxmath.PositiveModulo(x, int(planet.Width))
	if x >= int(planet.Width) || x < 0 || y >= int(planet.Height) || y < 0 {
		return -1
	}
	return y*int(planet.Width) + x
}

func (planet *Planet) ShortestDisplacement(from, to mgl32.Vec2) mgl32.Vec2 {
	disp := to.Sub(from)
	w := float32(planet.Width)
	disp[0] = math32.Mod(disp[0], w)
	if disp.X() > w/2 {
		to = to.Add(mgl32.Vec2{-w, 0})
	}
	if disp.X() < -w/2 {
		to = to.Add(mgl32.Vec2{w, 0})
	}
	disp = to.Sub(from)
	disp[0] = math32.Mod(disp[0], w)
	return disp
}

func (planet *Planet) GetTile(x, y int, layerID LayerID) *Tile {
	idx := planet.GetTileIndex(x, y)
	if idx >= 0 {
		return &planet.GetLayerTiles(layerID)[planet.GetTileIndex(x, y)]
	} else {
		return nil
	}
}

func (planet *Planet) GetAllTilesUnique() []Tile {
	tiles := []Tile{}
	seenTiles := map[Tile]bool{}
	for _, tile := range planet.Layers[TopLayer].Tiles {
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
	layer LayerID,
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
		if !layer.Valid() {
			return false
		}
		planetLayer := planet.GetLayerTiles(layer)
		if planetLayer[tileIdx].TileCategory == TileCategoryChild ||
			planetLayer[tileIdx].TileCategory == TileCategoryMulti {
			planet.RemoveParentTile(planetLayer, tileIdx)
		}
		planetLayer[tileIdx] = tile
		return true
	}
	return false
}

func (planet *Planet) PlaceTileType(tileTypeID TileTypeID, x, y int) {
	tileType, ok := GetTileTypeByID(tileTypeID)
	if !ok {
		log.Fatalf("cannot find tile type for id [%v]", tileTypeID)
	}
	tilesInLayer := planet.GetLayerTiles(tileType.Layer)
	tilesInLayer[planet.GetTileIndex(x, y)] =
		tileType.CreateTile(TileCreationOptions{
			Neighbours: planet.GetNeighbours(tilesInLayer, x, y),
		})
	planet.updateSurroundingTiles(tilesInLayer, x, y)
}

func (planet *Planet) updateSurroundingTiles(
	tilesInLayer []Tile, x, y int,
) {
	for xOffset := -1; xOffset <= 1; xOffset++ {
		for yOffset := -1; yOffset <= 1; yOffset++ {
			if xOffset != 0 || yOffset != 0 {
				planet.updateTile(tilesInLayer, x+xOffset, y+yOffset)
			}
		}
	}
}

func (planet *Planet) updateTile(tilesInLayer []Tile, x, y int) {
	tileIdx := planet.GetTileIndex(x, y)
	if tileIdx >= 0 {
		tile := &tilesInLayer[tileIdx]
		tileType, ok := GetTileTypeByID(tile.TileTypeID)
		if ok {
			neighbours := planet.GetNeighbours(tilesInLayer, x, y)
			tileType.UpdateTile(TileUpdateOptions{
				Tile: tile, Neighbours: neighbours,
			})
		}
	}
}

func (planet *Planet) TryPlaceMultiTile(
	x, y float32, layerID LayerID, multiTile MultiTile, cam *camera.Camera,
) bool {
	// click relative to camera
	camCoords := mgl32.Vec4{x / render.PixelsPerTile, y / render.PixelsPerTile, 0, 1}
	// click relative to world
	worldCoords := cam.GetTransform().Mul4x1(camCoords)
	tileX := int32(math.Round((float64(worldCoords.X()))))
	tileY := int32(math.Round((float64(worldCoords.Y()))))
	if tileX >= 0 && tileX < planet.Width && tileY >= 0 && tileY < planet.Width {
		tileIdx := planet.GetTileIndex(int(tileX), int(tileY))
		planetLayer := planet.GetLayerTiles(layerID)
		if len(planetLayer) == 0 {
			return false
		}
		if planetLayer[tileIdx].TileCategory == TileCategoryChild ||
			planetLayer[tileIdx].TileCategory == TileCategoryMulti {
			planet.RemoveParentTile(planetLayer, tileIdx)
		}
		planet.PlaceMultiTile(int(tileX), int(tileY), layerID, multiTile)
		return true
	}
	return false
}

// note that multi-tiles are assumed to be rectangular
func (planet *Planet) getMultiTileWidth(layer []Tile, x, y int) int {
	offsetX := 1
	for layer[planet.GetTileIndex(x+offsetX, y)].OffsetX == int8(offsetX) {
		offsetX++
	}
	return offsetX
}

func (planet *Planet) getMultiTileHeight(layer []Tile, x, y int) int {
	offsetY := 1
	for layer[planet.GetTileIndex(x, y+offsetY)].OffsetY == int8(offsetY) {
		offsetY++
	}
	return offsetY
}

func (planet *Planet) RemoveParentTile(layer []Tile, idx int) {
	tile := layer[idx]
	parentY := idx/int(planet.Width) - int(tile.OffsetY)
	parentX := idx%int(planet.Width) - int(tile.OffsetX)

	width := planet.getMultiTileWidth(layer, parentX, parentY)
	height := planet.getMultiTileHeight(layer, parentX, parentY)

	for y := parentY; y < parentY+height; y++ {
		for x := parentX; x < parentX+width; x++ {

			layer[planet.GetTileIndex(x, y)] = Tile{
				TileCategory: TileCategoryNone,
			}

		}
	}

}

func (planet *Planet) GetLayerTiles(layerID LayerID) []Tile {
	return planet.Layers[layerID].Tiles
}

func (planet *Planet) PlaceMultiTile(
	left, bottom int, layerID LayerID, mt MultiTile,
) {
	planetLayer := planet.GetLayerTiles(layerID)

	// place master tile
	planetLayer[planet.GetTileIndex(left, bottom)] = Tile{
		SpriteID:     mt.SpriteIDs[0],
		TileCategory: TileCategoryMulti,
		Name:         mt.Name,
		// (0,0) offset indicates master / standalone tile
		OffsetX: 0, OffsetY: 0,
	}

	for spriteIdIdx := 1; spriteIdIdx < len(mt.SpriteIDs); spriteIdIdx++ {
		localY := spriteIdIdx / mt.Width
		localX := spriteIdIdx % mt.Width

		x := left + localX
		y := bottom + localY
		tileIdx := planet.GetTileIndex(x, y)

		planetLayer[tileIdx] = Tile{
			SpriteID:     mt.SpriteIDs[spriteIdIdx],
			TileCategory: TileCategoryChild,
			Name:         mt.Name,
			OffsetX:      int8(localX), OffsetY: int8(localY),
		}
	}
}

// gets the y coordinate of the highest solid tile
// for a given column given by an x coordinate
func (planet *Planet) GetHeight(x int) int {
	for y := int(planet.Height - 1); y >= 0; y-- {
		idx := planet.GetTileIndex(x, y)
		if planet.Layers[TopLayer].Tiles[idx].TileCategory != TileCategoryNone {
			return y
		}
	}
	return 0
}

func (planet *Planet) GetTopLayerTile(x, y int) *Tile {
	index := planet.GetTileIndex(x, y)
	if index >= 0 {
		return &planet.Layers[TopLayer].Tiles[index]
	} else {
		return nil
	}
}

func (planet *Planet) TileIsSolid(x, y int) bool {
	tile := planet.GetTopLayerTile(x, y)
	return tile != nil &&
		tile.TileCategory != TileCategoryNone &&
		tile.TileCollisionType == TileCollisionTypeSolid
}

func (planet *Planet) TileTopIsSolid(x, y int, ignorePlatforms bool) bool {
	tile := planet.GetTopLayerTile(x, y)
	return tile != nil && tile.TileCategory != TileCategoryNone &&
		(tile.TileCollisionType == TileCollisionTypeSolid || !ignorePlatforms)
}

func (planet *Planet) TileExists(layerTiles []Tile, x, y int) bool {
	index := planet.GetTileIndex(x, y)
	if index < 0 {
		return false
	}
	tile := layerTiles[index]
	return tile.TileCategory != TileCategoryNone
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
			if planet.TileIsSolid(i, j) {
				fx := float32(i) - 0.5
				fy := float32(j) - 0.5
				fxw := fx + 1.0
				fyh := fy + 1.0

				// only draw the tiles outline instead of every single one
				if planet.TileIsSolid(i+1, j) { // right
					lines = append(lines, []float32{
						fxw, fyh,
						fxw, fy,
					}...)
				}
				if planet.TileIsSolid(i-1, j) { // left
					lines = append(lines, []float32{
						fx, fyh,
						fx, fy,
					}...)
				}
				if planet.TileIsSolid(i, j+1) { // up
					lines = append(lines, []float32{
						fx, fyh,
						fxw, fyh,
					}...)
				}
				if planet.TileIsSolid(i, j-1) { // down
					lines = append(lines, []float32{
						fx, fy,
						fxw, fy,
					}...)
				}
			}
		}
	}

	// save array
	planet.collidingLines = lines
	return planet.collidingLines
}

func (planet *Planet) DamageTile(
	x, y int, layerID LayerID,
) (tileCopy Tile, destroyed bool) {
	tileIdx := planet.GetTileIndex(x, y)
	if tileIdx < 0 {
		// invalid tile; nothing to damage
		return
	}
	tile := &planet.GetLayerTiles(layerID)[tileIdx]
	_tileCopy := *tile
	// TODO create tile chunk from collision point rather than tile center
	particles.CreateTileChunks(float32(x), float32(y), tile.SpriteID)
	// TODO use more robust system
	tileType, ok := GetTileTypeByID(tile.TileTypeID)
	if ok && !tileType.Invulnerable {
		tile.Durability--
	}
	_destroyed := tile.Durability <= 0
	if _destroyed {
		*tile = NewEmptyTile()
		planet.updateSurroundingTiles(planet.GetLayerTiles(layerID), x, y)
	}
	return _tileCopy, _destroyed
}

func (planet *Planet) WrapAround(pos mgl32.Vec2) mgl32.Vec2 {
	w := float32(planet.Width)
	if pos.X() < 0 {
		return mgl32.Vec2{pos.X() + w, pos.Y()}
	} else if pos.X() > w {
		return mgl32.Vec2{pos.X() - w, pos.Y()}
	}
	return pos
}

func (planet *Planet) WrapAroundOffset(before mgl32.Vec2) mgl32.Vec2 {
	after := planet.WrapAround(before)
	return after.Sub(before)
}

func (planet *Planet) MinimizeDistance(
	from, to mgl32.Vec2,
) (mgl32.Vec2, mgl32.Vec2) {
	toCloser := mgl32.Vec2{
		from.X() + cxmath.NewModular(float32(planet.Width)).Disp(from.X(), to.X()),
		to.Y(),
	}
	return from, toCloser
}

func (planet *Planet) Update(dt float32) {
	planet.Time += dt
}
