package spriteloader

import (
	"log"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
)

type GPUTexture struct {
	Gl            uint32
	Width, Height int
}

func (tex GPUTexture) Dims() cxmath.Vec2i {
	return cxmath.Vec2i{int32(tex.Width), int32(tex.Height)}
}

func LoadTextureFromFileToGPU(fname string) GPUTexture {
	status, img, _ := LoadPng(fname)
	if status != LoadOk {
		log.Fatalf("cannot upload [%v] to GPU", fname)
	}
	tex := MakeTexture(img)
	return GPUTexture{
		Gl:    tex,
		Width: img.Bounds().Dx(), Height: img.Bounds().Dy(),
	}
}

func LoadTextureArrayFromFileToGPU(fname string, config SpriteSheetConfig) GPUTexture {
	status, img, _ := LoadPng(fname)

	if status != LoadOk {
		log.Fatalf("cannot upload [%v] to GPU", fname)
	}

	tileW, tileH := config.CellWidth, config.CellHeight
	if config.CellWidth == 0 {
		tileW = constants.DEFAULT_SPRITE_SIZE
	}
	if config.CellHeight == 0 {
		tileH = constants.DEFAULT_SPRITE_SIZE
	}
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	tilesX := width / tileW
	tilesY := height / tileH
	tex := MakeTextureArray(img, tileW, tileH, tilesX, tilesY)
	return GPUTexture{
		Gl:     tex,
		Width:  width,
		Height: height,
	}
}
