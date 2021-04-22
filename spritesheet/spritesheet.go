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

var nextSpriteSheetId=1

func LoadSpriteSheet(fname string) int {
	log.Print("loading sprite sheet from "+fname)

	img := LoadPng(fname)
	_=img

	var spriteSheetId = nextSpriteSheetId
	nextSpriteSheetId++
	return spriteSheetId
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
