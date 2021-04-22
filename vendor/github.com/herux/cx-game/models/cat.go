package models

import (
	"image"
	"image/png"
	"log"
	"os"
)

type Cat struct {
	RGBA   *image.RGBA
	Image  image.Image
	width  int
	height int
}

func NewCat() *Cat {
	imageFile, err := os.Open("./assets/cat.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer imageFile.Close()

	imageDecoded, err := png.Decode(imageFile)
	if err != nil {
		log.Fatalln(err)
	}
	rgba := image.NewRGBA(imageDecoded.Bounds())

	cat := Cat{
		RGBA:   rgba,
		Image:  imageDecoded,
		width:  2,
		height: 2,
	}

	return &cat
}
