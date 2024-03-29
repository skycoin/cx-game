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

var cachedGPUTextures = map[string]GPUTexture{}

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

func LoadTextureFromFileToGPUCached(fname string) GPUTexture {
	gpuTexture, ok := cachedGPUTextures[fname]
	if ok {
		return gpuTexture
	}
	gpuTexture = LoadTextureFromFileToGPU(fname)
	cachedGPUTextures[fname] = gpuTexture
	return gpuTexture
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
	if len(config.SpriteConfigs) == 0 && config.Autoname == "" {
		//the whole picture is one sprite
		tileW = width
		tileH = height
	}
	tilesX := width / tileW
	tilesY := height / tileH
	tex := MakeTextureArray(img, tileW, tileH, tilesX, tilesY)
	return GPUTexture{
		Gl:     tex,
		Width:  width,
		Height: height,
	}
}
