package model

import (
	"image"
	"image/png"
	"log"
	"os"
)

type Cat struct {
	image  *image.RGBA
	width  int
	height int
}

func NewCat() *Cat {
	imgFile, err := os.Open("./assets/cat.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	im, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}
	img := image.NewRGBA(im.Bounds())

	cat := Cat{
		image:  img,
		width:  2,
		height: 2,
	}

	return &cat
}
