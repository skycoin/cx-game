package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	perlin "github.com/skycoin/cx-game/procgen"
)

func genNoiseMap() {
	width := 256
	height := 256
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	pNoise := perlin.NewPerlin2D(10, 256, 1, 255)

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			//val := uint8(math.Abs(float64(pNoise.Noise(float32(x)*0.02, float32(y)*0.02) * 255)
			normVal := (pNoise.Noise(float32(x)*0.02, float32(y)*0.02) + 1.0) / 2.0
			val := uint8(normVal * 255)
			img.Set(x, y, color.RGBA{val, val, val, 255})
		}
	}

	file, _ := os.Create("test_noise.png")
	png.Encode(file, img)
}

func main() {
	genNoiseMap()
}
