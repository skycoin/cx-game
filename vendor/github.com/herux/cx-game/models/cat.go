package models

import (
	"image"
	"image/png"
	"log"
	"os"
)

type Cat struct {
	Image        *image.RGBA
	ImageDecoded *image.Image
	width        int
	height       int
}

func NewCat() *Cat {
	imgFile, err := os.Open("./assets/cat.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	imageDecoded, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}
	image := image.NewRGBA(imageDecoded.Bounds())

	cat := Cat{
		Image:        image,
		ImageDecoded: &imageDecoded,
		width:        2,
		height:       2,
	}

	return &cat
}
