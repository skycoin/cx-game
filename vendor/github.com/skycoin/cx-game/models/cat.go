package models

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

type Cat struct {
	RGBA      *image.RGBA
	Size      image.Point
	width     int
	height    int
	XVelocity float32
	YVelocity float32
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
	draw.Draw(rgba, rgba.Bounds(), imageDecoded, image.Pt(0, 0), draw.Src)
	cat := Cat{
		RGBA:   rgba,
		Size:   rgba.Rect.Size(),
		width:  2,
		height: 2,
	}

	return &cat
}
