package world

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

type tilePos [2]float32

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
			planet.LightingValues[idx].SetSkyLight(35)
		}

	}
}

func (planet *Planet) PushSkyLight(xTile, yTile int) {
	//adds to queue
	slStartIndex = (slStartIndex + 1) % slLengthMax
	slNum++
}

func (planet *Planet) UpdateSkyLight(iterations int) {
	//do update tile logic on tiles in update queue
}

func (planet *Planet) PushEnvLight(xTile, yTile int) {

}

func (planet *Planet) UpdateEnvLight(iterations int) {

}
