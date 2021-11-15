package world

func (planet *Planet) TogglePower(x,y int, on bool) {
	tile,ok := planet.GetTile(x,y, MidLayer)
	if !ok { return }
	tile.Power.On = on
	tilesInLayer := planet.GetLayerTiles(MidLayer)
	planet.updateTile(tilesInLayer, x,y)
}
