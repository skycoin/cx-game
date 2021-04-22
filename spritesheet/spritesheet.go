package spritesheet

import (
	"os"
	"image/png"
	"image"
	/*
	"image/draw"
	"image/png"
	*/
	"log"
)

type Spritesheet struct {
	img image.Image
	spriteWidth int
	spriteHeight int
}

var spritesheets = []Spritesheet{};

func LoadSpriteSheet(fname string) int {
	log.Print("loading sprite sheet from "+fname)

	img := LoadPng(fname)

	spritesheets = append(spritesheets, Spritesheet{
		img: img,
		spriteWidth: img.Bounds().Dx() / 32,
		spriteHeight: img.Bounds().Dx() / 32,
	})

	return len(spritesheets)-1
}

func LoadPng(fname string) image.Image {
	imgFile, err := os.Open(fname)

	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	im, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}

	return image.NewRGBA(im.Bounds())
}
