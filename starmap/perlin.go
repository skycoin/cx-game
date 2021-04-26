package starmap

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	perlin "github.com/skycoin/cx-game/procgen"
)

func Generate(size int, scale float32, levels uint8) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	perlinField := perlin.NewPerlin2D(1, 512, 4, 256)
	max := float32(math.Sqrt2 / 1.9)
	min := float32(-math.Sqrt2 / 1.9)

	// Set color for each pixel.
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			val := perlinField.Noise(float32(x)*scale, float32(y)*scale, 0.5, 2, 8)
			val = (val - min) / (max - min) // normalized aproximation
			//brightness := uint8(val*float32(levels)) * (255 / levels)
			brightness := uint8(val * 255)
			img.Set(x, y, color.RGBA{brightness, brightness, brightness, 255})
		}
	}

	file, _ := os.Create("test_noise.png")
	png.Encode(file, img)
}
