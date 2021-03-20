package main

import (
	"time"

	"github.com/skycoin/cx-game/world"
)

var (
	FPS int
)

var CurrentPlanet *world.Planet

func main() {
	for {
		Tick()
		time.Sleep(100 * time.Millisecond)
		Draw()
	}
}

func Tick() {

}

func Draw() {

}
