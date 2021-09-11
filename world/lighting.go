package world

import (
	"fmt"
	"log"
	"time"
)

type LightValue uint16

type LightAttenuationType uint8

//upper 8 bits
func (l LightValue) GetSkyLight() uint8 {
	return uint8((l & 0xff00) >> 8)
}

func (l *LightValue) SetSkyLight(value uint8) {
	*l = (*l & 0xff) | LightValue(value)<<8
}

//lower 8 bits
func (l LightValue) GetEnvLight() uint8 {
	return uint8(l & 0xff)
}
func (l *LightValue) SetEnvLight(value uint8) {
	*l = (*l & 0xff00) | LightValue(value)
}

type tilePos struct {
	X int
	Y int
}

var skyLightUpdateQueue []tilePos = make([]tilePos, slLengthMax)
var slStartIndex int = 0
var slNum int = 0
var slLengthMax int = 100000

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
			planet.LightingValues[idx].SetSkyLight(255)
			// planet.PushSkyLightUpdate(x, y)
		}
		for ; y >= 0; y-- {
			// tile := planet.GetTile(x, y, TopLayer)
			//tile not empty
			idx := planet.GetTileIndex(x, y)
			topTileIdx := planet.GetTileIndex(x, y+1)
			topTileValue := planet.LightingValues[topTileIdx]
			if topTileValue.GetSkyLight() < 2*16 {
				topTileValue.SetSkyLight(2 * 16)
			}

			planet.LightingValues[idx].SetSkyLight(topTileValue.GetSkyLight() - 2*16)
			planet.PushSkylight(x, y)
		}

		for slNum != 0 {
			planet.UpdateSkyLight(10000)
		}
	}

	startTimer := time.Now()
	for slNum != 0 && time.Since(startTimer).Seconds() < 5 {
		planet.UpdateSkyLight(64 * 1024)
	}

}

func (planet *Planet) LightAddBlock(xtile, yTile int) {
	planet.PushSkylight(xtile, yTile)
}

func (planet *Planet) PushSkylight(xTile, yTile int) {
	if slNum == slLengthMax {
		// fmt.Println("HIT SKYLIGHT LIMIT")
		return
	}
	//adds to queue
	// slStartIndex = (slStartIndex + 1) % slLengthMax
	// fmt.Println(xTile, yTile)
	skyLightUpdateQueue[(slStartIndex+slNum)%slLengthMax] = tilePos{xTile, yTile}
	slNum++
}

var neighboursOffsets = []int{
	// -1, 1, //nw
	0, 1, //n
	// 1, 1, //ne
	-1, 0, //w
	1, 0, //e
	// -1, -1, //sw
	0, -1, //s
	// 1, -1, //se
}

func (planet *Planet) UpdateSkyLight(iterations int) {
	// fmt.Println(slNum)
	//only top tile is accounted in calculations
	//do update tile logic on tiles in update queue
	if slNum == 0 {
		return
	} else if slNum < 0 {
		log.Fatalln("Skylight update error 1")
	}

	fmt.Println("updating: ", slNum)
	//update logic
	for i := 0; i < iterations; i++ {
		if slNum == 0 {
			return
		}

		pos := skyLightUpdateQueue[slStartIndex]

		slStartIndex = (slStartIndex + 1) % slLengthMax
		slNum--
		// fmt.Println(slNum, "   ttt")

		idx := planet.GetTileIndex(pos.X, pos.Y)
		if idx == -1 {
			continue
			// log.Fatalln("Lighting error")
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

		//placed solid block below sunlight
		if topTile.Name == "" && topTileLightValue.GetSkyLight() == 255 {
			//toptile is skylight
			if lightTile.Name != "" && lightVal.GetSkyLight() != 255-2*16 {
				//if checked tile is solid
				planet.LightingValues[idx].SetSkyLight(255 - 2*16)
				for i := 0; i < 4; i++ {
					planet.PushSkylight(pos.X, pos.Y)
					fmt.Print("NEIGHBOURS ARE: ", pos.X+neighboursOffsets[i*2], pos.Y+neighboursOffsets[i*2+1], "   ")

					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)

				}
				// fmt.Println("   POS: ", pos.X, pos.Y)
				// fmt.Println("FIRST")
			} else if lightTile.Name == "" && lightVal.GetSkyLight() != 255 {
				//removed block, waiting for sunlight instead
				fmt.Println(pos.X, pos.Y, "SECOND")
				planet.LightingValues[idx].SetSkyLight(255)
				for i := 0; i < 4; i++ {
					planet.PushSkylight(pos.X, pos.Y)

					if pos.X+neighboursOffsets[i*2] == 25 &&
						pos.Y+neighboursOffsets[i*2+1] == 20 {
						fmt.Println("YESYSEYESY SYE")
					}
					planet.PushSkylight(
						pos.X+neighboursOffsets[i*2],
						pos.Y+neighboursOffsets[i*2+1],
					)
				}
			} else {
				// fmt.Println("WTFF MANNA", "   ", pos.X, pos.Y, "     ", lightVal.GetSkyLight())
			}
			continue
		}

		// fmt.Println("GOT: ", pos.X, "   ", pos.Y)
		var neighbourIndexes [4]int
		for i := 0; i < 4; i++ {
			idx := planet.GetTileIndex(
				pos.X+neighboursOffsets[i*2],
				pos.Y+neighboursOffsets[i*2+1],
			)
			neighbourIndexes[i] = idx
		}

		var maxSkylightValue uint8 = 0

		//determine brightest neighbour out of 4
		for i := 0; i < 4; i++ {
			if neighbourIndexes[i] == -1 {
				continue
			}
			neighbourLightValue := planet.LightingValues[neighbourIndexes[i]]
			if neighbourLightValue.GetSkyLight() > maxSkylightValue {
				maxSkylightValue = neighbourLightValue.GetSkyLight()
			}
		}

		if lightSkyLightValue >= maxSkylightValue && lightSkyLightValue > 0 {
			// if lightSkyLightValue+2*16 != maxSkylightValue {
			// lightVal.SetSkyLight(maxSkylightValue - 2*16)
			min := uint8(2 * 16)
			if lightSkyLightValue == 31 {
				min -= 1
			}
			planet.LightingValues[idx].SetSkyLight(lightSkyLightValue - min)
			for i := 0; i < 4; i++ {
				planet.PushSkylight(
					pos.X+neighboursOffsets[i*2],
					pos.Y+neighboursOffsets[i*2+1],
				)
			}
			continue
		}
		if maxSkylightValue > lightSkyLightValue+2*16 {
			planet.LightingValues[idx].SetSkyLight(maxSkylightValue - 2*16)
			for i := 0; i < 4; i++ {
				planet.PushSkylight(
					pos.X+neighboursOffsets[i*2],
					pos.Y+neighboursOffsets[i*2+1],
				)
			}
			continue
		}
		// topTileIdx := planet.GetTileIndex(pos.X, pos.Y+1)
		// if topTileIdx == -1 {
		// 	log.Println("Warning update skylight")
		// 	continue
		// }

	}
}

func (planet *Planet) PushEnvLight(xTile, yTile int) {

}

func (planet *Planet) UpdateEnvLight(iterations int) {

}
