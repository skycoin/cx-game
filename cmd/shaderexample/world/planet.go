package world

type Planet struct {
	Width  int32
	Height int32
	Tiles  []Tile
}

func NewPlanet(x, y int32) *Planet {
	planet := Planet{
		Width:  x,
		Height: y,
	}

	return &planet
}
