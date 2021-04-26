package starmap

import (
	"image"
	"image/color"
	"image/png"
	"os"

	perlin "github.com/skycoin/cx-game/procgen"
)

func Generate(size int, levels uint8) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	perlinField := perlin.NewPerlin2D(10, 256, 1, 255)

	// Set color for each pixel.
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			normVal := (perlinField.One_over_f_pers(float32(x)*0.02, float32(y)*0.02, 0.5) + 1.0) / 2.0
			//val := uint8(normVal*float32(levels)) * (255 / levels)
			val := uint8(normVal * 255)
			img.Set(x, y, color.RGBA{val, val, val, 255})
		}
	}

	file, _ := os.Create("test_noise.png")
	png.Encode(file, img)
}
