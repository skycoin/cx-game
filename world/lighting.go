package world

import "log"

type LightValue uint16

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

type tilePos [2]int

var skyLightUpdateQueue []tilePos = make([]tilePos, slLengthMax)
var slStartIndex int = 0
var slNum int = 0
var slLengthMax int = 10000

func (planet *Planet) InitLighting() {
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
		}
	}
}

func (planet *Planet) PushSkyLight(xTile, yTile int) {
	if len(skyLightUpdateQueue) == slLengthMax {
		return
	} else if len(skyLightUpdateQueue) > slLengthMax {
		log.Fatalf("Length exceeded")
	}
	//adds to queue
	slStartIndex = (slStartIndex + 1) % slLengthMax
	slNum++

	skyLightUpdateQueue[slStartIndex] = tilePos{xTile, yTile}
}

func (planet *Planet) UpdateSkyLight(iterations int) {
	//do update tile logic on tiles in update queue
	if slNum == 0 {
		return
	} else if slNum < 0 {
		log.Fatalln("Skylight update error 1")
	}

	// for _, tile := range skyLightUpdateQueue {
	// 	slStartIndexпш
	// 	idx := planet.GetTileIndex(tile[0], tile[1])

	// 	topTileIdx := planet.GetTileIndex(tile[0], tile[1]+1)

	// }
}

func (planet *Planet) PushEnvLight(xTile, yTile int) {

}

func (planet *Planet) UpdateEnvLight(iterations int) {

}
