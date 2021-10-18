package world

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world/tiling"
)

const NUM_INSTANCES = 100

type Planet struct {
	Width           int32
	Height          int32
	Layers          Layers
	collidingLines  []float32
	collidingLinesX int
	collidingLinesY int
	Time            float32
	LightingValues  []LightValue

	program, liquidProgram render.Program
}

// TODO move to render
func newPlanetLiquidProgram() render.Program {
	config := render.NewShaderConfig(
		"./assets/shader/world.vert", "./assets/shader/liquid.frag")
	config.Define("NUM_INSTANCES", strconv.Itoa(NUM_INSTANCES))
	program := config.Compile()

	return program
}

func NewPlanet(x, y int32) *Planet {
	planet := Planet{
		Width:          x,
		Height:         y,
		Layers:         NewLayers(x * y),
		liquidProgram:  newPlanetLiquidProgram(),
		LightingValues: make([]LightValue, x*y),
	}
	return &planet
}

func (planet *Planet) GetNeighbours(
	layer []Tile, x, y int, id TileTypeID,
) tiling.DetailedNeighbours {
	return tiling.DetailedNeighbours{
		Up:        planet.GetNeighbour(layer, x, y+1, id),
		UpRight:   planet.GetNeighbour(layer, x+1, y+1, id),
		Right:     planet.GetNeighbour(layer, x+1, y, id),
		DownRight: planet.GetNeighbour(layer, x+1, y-1, id),
		Down:      planet.GetNeighbour(layer, x, y-1, id),
		DownLeft:  planet.GetNeighbour(layer, x-1, y-1, id),
		Left:      planet.GetNeighbour(layer, x-1, y, id),
		UpLeft:    planet.GetNeighbour(layer, x-1, y+1, id),
	}
}

func (planet *Planet) GetNeighbour(
	layerTiles []Tile, x, y int, id TileTypeID,
) tiling.Neighbour {
	index := planet.GetTileIndex(x, y)
	if index < 0 {
		return tiling.None
	}
	tile := layerTiles[index]
	if tile.TileTypeID == id {
		return tiling.Self
	}
	if tile.TileCategory == TileCategoryNone {
		return tiling.None
	}
	return tiling.Solid
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

func (planet *Planet) PlaceTileTypeNoConnect(tileTypeID TileTypeID, x, y int) {
	tileType, ok := GetTileTypeByID(tileTypeID)
	if !ok {
		log.Fatalf("cannot find tile type for id [%v]", tileTypeID)
	}
	tilesInLayer := planet.GetLayerTiles(tileType.Layer)
	rootTileIdx := planet.GetTileIndex(x, y)
	if rootTileIdx == -1 {
		return
	}
	tilesInLayer[rootTileIdx] =
		tileType.CreateTile(TileCreationOptions{
			//Neighbours: planet.GetNeighbours(tilesInLayer, x, y, tileTypeID),
		})
	rect := cxmath.Rect{
		cxmath.Vec2i{int32(x), int32(y)},
		tileType.Size(),
	}

	// place child tiles (non-root) to prevent overlap
	for offsetX := int32(0); offsetX < rect.Size.X; offsetX++ {
		for offsetY := int32(0); offsetY < rect.Size.Y; offsetY++ {
			// don't overwrite root
			if offsetX != 0 || offsetY != 0 {
				tileIdx := planet.GetTileIndex(
					int(rect.Origin.X+offsetX),
					int(rect.Origin.Y+offsetY),
				)
				if tileIdx != -1 {

					tilesInLayer[tileIdx] = Tile{
						TileCategory: TileCategoryChild,
						OffsetX:      int8(offsetX), OffsetY: int8(offsetY),
						Name: fmt.Sprintf("%s (child)", tileType.Name),
					}
				}
			}
		}
	}

	planet.LightUpdateBlock(x, y)
}

func (planet *Planet) PlaceTileType(tileTypeID TileTypeID, x, y int) {
	tileType, ok := GetTileTypeByID(tileTypeID)
	if !ok {
		log.Fatalf("cannot find tile type for id [%v]", tileTypeID)
	}
	tilesInLayer := planet.GetLayerTiles(tileType.Layer)
	planet.PlaceTileTypeNoConnect(tileTypeID, x, y)
	rect := cxmath.Rect{
		cxmath.Vec2i{int32(x), int32(y)},
		tileType.Size(),
	}
	for _, neighbour := range rect.Neighbours() {
		planet.updateTile(tilesInLayer, int(neighbour.X), int(neighbour.Y))
	}
	planet.updateTile(tilesInLayer, x, y)
}

// cycle the pipe connection state at (x,y) to the next valid state.
// also updates connections of neighbouring pipe tiles.
func (planet *Planet) TryCyclePipeConnection(x, y int) {
	layerTiles := planet.GetLayerTiles(PipeLayer)
	tileIdx := planet.GetTileIndex(x, y)
	if tileIdx < 0 {
		return
	}
	tile := &layerTiles[tileIdx]
	if tile.TileCategory == TileCategoryNone {
		return
	}
	oldConnections := tile.Connections
	tile.Connections =
		tile.Connections.Next(planet.PipeConnectionCandidates(x, y))

	tileType := tile.TileTypeID.Get()
	tileType.UpdateTile(TileUpdateOptions{
		Tile:       tile,
		Cycling:    true,
		Neighbours: planet.GetNeighbours(layerTiles, x, y, tileType.ID),
	})

	neighbours := pipeNeighbours(x, y, tile.Connections, oldConnections)
	for _, neighbour := range neighbours {
		neighbourTileIdx :=
			planet.GetTileIndex(neighbour.X, neighbour.Y)
		if neighbourTileIdx >= 0 {
			neighbourTile := &layerTiles[neighbourTileIdx]
			if neighbourTile.TileCategory != TileCategoryNone {
				neighbourTile.Connections = neighbourTile.Connections.
					ApplyDiff(neighbour.ConnectionDiff)
				neighbourTileType := neighbourTile.TileTypeID.Get()
				neighbourTileType.UpdateTile(TileUpdateOptions{
					Tile:    neighbourTile,
					Cycling: true,
					Neighbours: planet.GetNeighbours(
						layerTiles,
						neighbour.X, neighbour.Y, tileType.ID,
					),
				})
			}
		}
	}
}

// which pipes can the pipe at (x,y) be connected to?
func (planet *Planet) PipeConnectionCandidates(x, y int) Connections {
	layerTiles := planet.GetLayerTiles(PipeLayer)
	return Connections{
		Up:    planet.TileExists(layerTiles, x, y+1),
		Down:  planet.TileExists(layerTiles, x, y-1),
		Left:  planet.TileExists(layerTiles, x-1, y),
		Right: planet.TileExists(layerTiles, x+1, y),
	}
}

// update surrounding tiles. called after place/destroy tile.
// necessary for auto-tiling mechanism.
func (planet *Planet) updateSurroundingTiles(
	tilesInLayer []Tile, x, y int,
) {
	rect := cxmath.Rect{
		cxmath.Vec2i{int32(x), int32(y)},
		cxmath.Vec2i{1, 1},
	}
	for _, neighbour := range rect.Neighbours() {
		planet.updateTile(tilesInLayer, int(neighbour.X), int(neighbour.Y))
	}
}

func (planet *Planet) updateTile(tilesInLayer []Tile, x, y int) {
	tileIdx := planet.GetTileIndex(x, y)
	if tileIdx >= 0 {
		tile := &tilesInLayer[tileIdx]
		tileType, ok := GetTileTypeByID(tile.TileTypeID)
		if ok {
			neighbours :=
				planet.GetNeighbours(tilesInLayer, x, y, tile.TileTypeID)
			tileType.UpdateTile(TileUpdateOptions{
				Tile: tile, Neighbours: neighbours,
			})
		}
	}
}

func (planet *Planet) GetLayerTiles(layerID LayerID) []Tile {
	return planet.Layers[layerID].Tiles
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

func (planet *Planet) DamageTileBy(tile Tile) (destroyed bool) {
	return
}

func (planet *Planet) DamageTile(
	x, y int, layerID LayerID,
) (tileCopy Tile, destroyed bool) {
	tileIdx := planet.GetTileIndex(x, y)
	if tileIdx < 0 {
		// invalid tile; nothing to damage
		return
	}
	child := &planet.GetLayerTiles(layerID)[tileIdx]
	tileCopy = *child

	parentX := x - int(child.OffsetX)
	parentY := y - int(child.OffsetY)
	parentIdx := planet.GetTileIndex(parentX, parentY)

	parent := &planet.GetLayerTiles(layerID)[parentIdx]
	// TODO create tile chunk from collision point rather than tile center
	//particles.CreateTileChunks(float32(x), float32(y), tile.SpriteID)
	// TODO use more robust system
	tileType, ok := GetTileTypeByID(parent.TileTypeID)
	if ok && !tileType.Invulnerable {
		parent.Durability--
	}
	destroyed = parent.Durability <= 0
	if destroyed {
		for offsetY := int32(0); offsetY < tileType.Height; offsetY++ {
			for offsetX := int32(0); offsetX < tileType.Width; offsetX++ {
				childX := parentX + int(offsetX)
				childY := parentY + int(offsetY)
				childIdx := planet.GetTileIndex(childX, childY)
				if childIdx >= 0 {
					planet.GetLayerTiles(layerID)[childIdx] = NewEmptyTile()
				}
			}
		}
		*parent = NewEmptyTile()
		planet.updateSurroundingTiles(planet.GetLayerTiles(layerID), x, y)
		if tileType.Layer == TopLayer {
			planet.LightUpdateBlock(x, y)
		}
	}
	return tileCopy, destroyed
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

func (planet *Planet) FixedUpdate() {
	planet.UpdateLighting()

}

func (planet *Planet) NearOxygenGenerator(position cxmath.Vec2) bool {
	//todo after, return always true for now
	if input.GetButton("bubbles") {
		// fmt.Println("Working")
		return true
	}
	return false
}
