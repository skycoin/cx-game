package world

import (
	"log"
)

type Layers struct {
	Background []Tile
	Mid []Tile
	Top []Tile
}

type Planet struct {
	Width  int32
	Height int32
	Layers Layers
}

func NewPlanet(x, y int32) *Planet {
	tiles := make([]Tile,x*y)
	planet := Planet{
		Width:  x,
		Height: y,
		Layers:  Layers{
			Background: tiles,
			Mid: tiles,
			Top: tiles,
		},
	}

	return &planet
}

func (planet *Planet) DrawLayer(tiles []Tile) {
	for idx,tile := range tiles {
		y := int32(idx)/planet.Width
		x := int32(idx)%planet.Width
		log.Print(x,y)
		log.Print(tile)
	}
}

func (planet *Planet) Draw() {
	log.Print("drawing background")
	planet.DrawLayer(planet.Layers.Background)
	log.Print("drawing mid")
	planet.DrawLayer(planet.Layers.Mid)
	log.Print("drawing top")
	planet.DrawLayer(planet.Layers.Top)
}
