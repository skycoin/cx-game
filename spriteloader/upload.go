package spriteloader

import (
	"github.com/skycoin/cx-game/cxmath"
)

type GPUTexture struct {
	Gl uint32
	Width, Height int
}

func (tex GPUTexture) Dims() cxmath.Vec2i {
	return cxmath.Vec2i { int32(tex.Width), int32(tex.Height) }
}

func LoadTextureFromFileToGPU(fname string) GPUTexture {
	_,img,_ := LoadPng(fname)
	tex := MakeTexture(img)
	return GPUTexture {
		Gl: tex,
		Width: img.Bounds().Dx(), Height: img.Bounds().Dy(),
	}
}
