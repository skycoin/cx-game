package main

import (
	"fmt"

	perlin "github.com/skycoin/cx-game/procgen"
)

type PerlinSettings struct {
}

func main() {
	myPerlin := perlin.NewPerlin2D(
		1,
		512,
		5,
		256,
	)

	result := myPerlin.Noise(15, 26, 0.5, 2, 8)
	fmt.Printf("%f", result)

}
