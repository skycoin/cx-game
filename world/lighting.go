package world

import (
	"log"
	"time"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type LightValue uint8

//upper 8 bits
func (l LightValue) GetSkyLight() uint8 {
	return uint8((l & 0xf0) >> 4)
}

func (l *LightValue) SetSkyLight(value uint8) {
	*l = (*l & 0xf) | LightValue(value)<<4
}

//lower 8 bits
func (l LightValue) GetEnvLight() uint8 {
	return uint8(l & 0xf)
}
func (l *LightValue) SetEnvLight(value uint8) {
	*l = (*l & 0xf0) | LightValue(value)
}

type tilePos struct {
	X int
	Y int
}

var skyLightUpdateQueue []tilePos = make([]tilePos, slLengthMax)
var slStartIndex int = 0
var slNum int = 0
var slLengthMax int = 10000
var neighbourCount int = 4

func getLightAttenuation(tile *Tile) types.LightAttenuationType {
	switch tile.Name {
	case "":
		return constants.ATTENUATION_DEFAULT
	default:
		return constants.ATTENUATION_SOLID
	}
}
func isSolid(tile *Tile) bool {
	return tile.Name != ""
}

//assume light attenuation for all block is 1
func (planet *Planet) InitSkyLight() {
	//init skylight
	for x := 0; x < int(planet.Width); x++ {
		y := int(planet.Height) - 1
		for ; y >= 0; y-- {
			tile := planet.GetTile(x, y, TopLayer)
			//tile not empty
			if tile.Name != "" {
				break
			}
			idx := planet.GetTileIndex(x, y)
			planet.LightingValues[idx].SetSkyLight(15)
			// planet.PushSkyLightUpdate(x, y)
		}
		for ; y >= 0; y-- {
			// tile := planet.GetTile(x, y, TopLayer)
			//tile not empty
			idx := planet.GetTileIndex(x, y)

			topTileIdx := planet.GetTileIndex(x, y+1)
			topTileValue := planet.LightingValues[topTileIdx]
			lightTile := planet.GetTile(x, y, TopLayer)

			attenuation := uint8(getLightAttenuation(lightTile))
			//avoid overflow
			if topTileValue.GetSkyLight() >= attenuation {
				planet.LightingValues[idx].SetSkyLight(topTileValue.GetSkyLight() - attenuation)
			} else {
				planet.LightingValues[idx].SetSkyLight(topTileValue.GetSkyLight())
			}
			planet.PushSkylight(x, y)
		}

	}

	startTimer := time.Now()
	for slNum != 0 && time.Since(startTimer).Seconds() < 1 {
		planet.UpdateSkyLight(64 * 1024)
	}
	if time.Since(startTimer).Seconds() > 1 {
		log.Println("SKYLIGHT UPDATE TOOK TOO LONG")
	}

}

func (planet *Planet) LightAddBlock(xtile, yTile int) {
	planet.PushSkylight(xtile, yTile)
}

func (planet *Planet) PushSkylight(xTile, yTile int) {
	if slNum == slLengthMax {
		return
	}
	//adds to queue
	skyLightUpdateQueue[(slStartIndex+slNum)%slLengthMax] = tilePos{xTile, yTile}
	slNum++
}

var neighboursOffsets = []int{
	0, 1, //n
	-1, 0, //w
	1, 0, //e
	0, -1, //s
	-1, 1, //nw
	1, 1, //ne
	-1, -1, //sw
	1, -1, //se
}

func (planet *Planet) UpdateSkyLight(iterations int) {
	//only top tile is accounted in calculations
	//do update tile logic on tiles in update queue
	if slNum == 0 {
		return
	} else if slNum < 0 {
		log.Fatalln("Skylight update error 1")
	}
	//update logic
	for i := 0; i < iterations; i++ {
		if slNum == 0 {
			return
		}

		pos := skyLightUpdateQueue[slStartIndex]

		slStartIndex = (slStartIndex + 1) % slLengthMax
		slNum--
		idx := planet.GetTileIndex(pos.X, pos.Y)
		if idx == -1 {
			continue

		}
		lightVal := planet.LightingValues[idx]
		lightTile := planet.GetTile(pos.X, pos.Y, TopLayer)
		lightSkyLightValue := lightVal.GetSkyLight()

		topTileIdx := planet.GetTileIndex(pos.X, pos.Y+1)
		if topTileIdx == -1 {
			continue
		}

		topTileLightValue := planet.LightingValues[topTileIdx]
		topTile := planet.GetTile(pos.X, pos.Y+1, TopLayer)

		attenuation := uint8(getLightAttenuation(lightTile))
		//placed solid block below sunlight
		if !isSolid(topTile) && topTileLightValue.GetSkyLight() == 15 {
			//toptile is skylight
			if isSolid(lightTile) && lightVal.GetSkyLight() != 15-2 {
				//if checked tile is solid
				planet.LightingValues[idx].SetSkyLight(15 - attenuation)
				for i := 0; i < neighbourCount; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
				planet.PushSkylight(pos.X, pos.Y)
			} else if !isSolid(lightTile) && lightVal.GetSkyLight() != 15 {
				//removed block, waiting for sunlight instead
				planet.LightingValues[idx].SetSkyLight(15)
				for i := 0; i < neighbourCount; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
				planet.PushSkylight(pos.X, pos.Y)
			} else {
			}
			continue
		}
		// continue
		neighbourIndexes := make([]int, neighbourCount)
		for i := 0; i < neighbourCount; i++ {
			idx := planet.GetTileIndex(
				pos.X+neighboursOffsets[i*2],
				pos.Y+neighboursOffsets[i*2+1],
			)
			neighbourIndexes[i] = idx
		}

		var maxSkylightValue uint8 = 0

		//determine brightest neighbour out of 4
		for i := 0; i < neighbourCount; i++ {
			if neighbourIndexes[i] == -1 {
				continue
			}
			neighbourLightValue := planet.LightingValues[neighbourIndexes[i]]
			if neighbourLightValue.GetSkyLight() > maxSkylightValue {
				maxSkylightValue = neighbourLightValue.GetSkyLight()
			}
		}

		if maxSkylightValue >= attenuation {
			if lightSkyLightValue != maxSkylightValue-attenuation {
				planet.LightingValues[idx].SetSkyLight(maxSkylightValue - attenuation)
				for i := 0; i < 4; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
		} else {
			if lightSkyLightValue != 0 {
				planet.LightingValues[idx].SetSkyLight(0)
				for i := 0; i < 4; i++ {
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			}
		}

	}
}

func SwitchNeighbourCount(planet *Planet) {
	if neighbourCount == 4 {
		log.Println("LIGHTING: 8 neighbours")
		neighbourCount = 8
	} else {
		log.Println("LIGHTING: 4 neighbours")
		neighbourCount = 4
	}
	planet.InitSkyLight()
}

func (planet *Planet) PushEnvLight(xTile, yTile int) {

}

func (planet *Planet) UpdateEnvLight(iterations int) {

}
