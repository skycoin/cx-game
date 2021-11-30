package world

import (
	"log"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/constants"
)

type CircuitID uint32
type Circuit []int // list of tile indices
type Circuits map[CircuitID]Circuit

func (c Circuit) Bind(planet *Planet) BoundCircuit {
	return BoundCircuit { Circuit:c, Planet: planet }
}

// circuit bound to a planet
type BoundCircuit struct {
	Circuit Circuit
	Planet  *Planet
}

func (bc *BoundCircuit) Tiles() []*Tile {
	topLayerTiles := bc.Planet.GetLayerTiles(MidLayer)
	tiles := []*Tile{}
	for _,tileIdx := range bc.Circuit {
		tile := &topLayerTiles[tileIdx]
		tiles = append(tiles,tile)
	}
	return tiles
}

func (bc *BoundCircuit) Wattage() int {
	wattage := 0
	for _,tile := range bc.Tiles() {
		log.Printf("tile %s has wattage %d", tile.Name, tile.Power.Wattage)
		wattage += tile.Power.Wattage
	}
	return wattage
}

func (bc *BoundCircuit) FixedTick() {
	log.Printf("circuit wattage is %d", bc.Wattage())
	log.Printf("circuit has %d electric tiles", len(bc.Tiles()))
	active := bc.Wattage() > 0
	if active { log.Printf("circuit is ON") }
	bc.Toggle(active)
}

func (bc *BoundCircuit) Toggle(active bool) {
	for _,tile := range bc.Tiles() {
		tile.Power.On = active
		tileType,ok := GetTileTypeByID(tile.TileTypeID)
		if ok {
			tileType.UpdateTile(TileUpdateOptions{Tile:tile})
		}
	}
}

func (planet *Planet) UpdateCircuits() {
	for _,circuit := range planet.Circuits {
		boundCircuit := circuit.Bind(planet)
		boundCircuit.FixedTick()
	}
}

func (planet *Planet) electricTilePositions() []cxmath.Vec2i {
	positions := []cxmath.Vec2i{}
	for y := 0 ; y < int(planet.Height) ; y++ {
		for x := 0 ; x < int(planet.Width) ; x++ {
			tile,ok := planet.GetTile(x,y, MidLayer)
			if ok && tile.IsElectric() {
				position := cxmath.Vec2i { int32(x), int32(y) }
				positions = append(positions, position)
			}
		}
	}
	return positions
}

func (planet *Planet) DetectCircuits() {
	positions := planet.electricTilePositions()
	clusters := cxmath.FindClusters(positions, constants.POWER_REACH_RADIUS)

	planet.Circuits = map[CircuitID]Circuit{}
	for clusterIdx,cluster := range clusters {
		clusterID := CircuitID(clusterIdx)
		planet.Circuits[clusterID] = Circuit{}
		for _,point := range cluster {
			tileIdx := planet.GetTileIndex(int(point.X), int(point.Y))
			planet.Circuits[clusterID] =
				append(planet.Circuits[clusterID], tileIdx)
		}
	}
}

func (planet *Planet) AddCircuitTile(at cxmath.Vec2i) {
	for circuitID,tileIndices := range planet.Circuits {
		for _,tileIdx := range tileIndices {
			pos := planet.TileIdxToPos(tileIdx)
			// TODO remove magic number
			if pos.Sub(at).ManhattanDist() < 10 {
				hereTileIdx := planet.GetTileIndex(int(at.X), int(at.Y))
				// add new tile to circuit
				planet.Circuits[circuitID] =
					append(planet.Circuits[circuitID], hereTileIdx)
				log.Printf("added tile to circuit")
			}
		}
	}
	log.Printf("could not find home for circuit")
	// TODO create new circuit
}
