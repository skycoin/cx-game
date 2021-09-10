package world

import (
	"fmt"
	"log"
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
var slLengthMax int = 10000

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

		// for slNum != 0 {

		// }
	}

	for slNum != 0 {
		planet.UpdateSkyLight(64 * 1024)

		fmt.Println("COUNTER IS: ", counter)
	}

}

func (planet *Planet) PushSkylight(xTile, yTile int) {
	if slNum == slLengthMax {
		// fmt.Println("HIT SKYLIGHT LIMIT")
		return
	} else if len(skyLightUpdateQueue) > slLengthMax {
		log.Fatalf("Length exceeded")
	}
	//adds to queue
	slStartIndex = (slStartIndex + 1) % slLengthMax
	slNum++

	skyLightUpdateQueue[slStartIndex] = tilePos{xTile, yTile}
}

var neighboursOffsets = []int{
	// -1, 1, //nw
	0, 1, //n
	// 1, 1, //ne
	-1, 0, //w
	1, 0, //e
	// -1, -1, //sw
	0, 1, //s
	// 1, -1, //se
}

var counter int

func (planet *Planet) UpdateSkyLight(iterations int) {

	counter++
	//do update tile logic on tiles in update queue
	if slNum == 0 {
		return
	} else if slNum < 0 {
		log.Fatalln("Skylight update error 1")
	}
	// fmt.Println(slNum)

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
		valueChanged := false

		if lightTile.Name == "" {
			//is not solid block
			continue
		}
		// lightingValue := planet.LightingValues[idx]
		// topTileIdx := planet.GetTileIndex(value.X, value.Y+1)
		// if idx == -1 {
		// 	log.Fatalln("Top tile lighting error")
		// }
		// topTileLightingValue := planet.LightingValues[idx]
		var neighbourIndexes [4]int
		for i := 0; i < 4; i++ {
			idx := planet.GetTileIndex(pos.X+neighboursOffsets[i], pos.Y+neighboursOffsets[i+1])
			neighbourIndexes[i] = idx
		}

		var maxSkylightValue uint8 = 0

		//determine brightest neighbour out of 8
		for i := 0; i < 4; i++ {
			if neighbourIndexes[i] == -1 {
				continue
			}
			neighbourLightValue := planet.LightingValues[neighbourIndexes[i]]
			if neighbourLightValue.GetSkyLight() > maxSkylightValue {
				maxSkylightValue = neighbourLightValue.GetSkyLight()
			}
		}

		if lightSkyLightValue != 255 && lightSkyLightValue >= maxSkylightValue {
			lightVal.SetSkyLight(maxSkylightValue - 16)
			valueChanged = true
		}

		if valueChanged {
			for i := 0; i < 4; i++ {
				planet.PushSkylight(pos.X+neighboursOffsets[i], pos.Y+neighboursOffsets[i+1])
			}
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
