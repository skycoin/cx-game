package Texture

import (
	"image"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/geometry"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/test/spine-animation/animation"
)

type Texture struct {
	M_renderID       uint32
	M_filePath       string
	M_localBufferImg *image.RGBA
	M_width          int
	M_height         int
	M_BPP            int
	M_metrix         geometry.Matrix
	M_name           string
	//	M_matrix geometry.Matrix
	GeoM   animation.GeoM
	ColorM animation.ColorM
}

// type Matrix [6]float64

// var IM = Matrix{1, 0, 0, 1, 0, 0}

func ReadTextureData(path string) *Texture {

	spriteloader.InitSpriteloaderDev()

	status, img, _ := spriteloader.LoadPng(path)

	if status != spriteloader.LoadOk {
		log.Fatalf("cannot upload [%v] to GPU", path)
	}

	texture := &Texture{M_renderID: 0, M_filePath: path, M_localBufferImg: img, M_width: img.Rect.Dx(), M_height: img.Rect.Dy(), M_name: ""}

	return texture

}
func SetUpTexture(path string) *Texture {

	spriteloader.InitSpriteloaderDev()

	status, img, _ := spriteloader.LoadPng(path)

	if status != spriteloader.LoadOk {
		log.Fatalf("cannot upload [%v] to GPU", path)
	}

	texture := &Texture{M_renderID: 0, M_filePath: path, M_localBufferImg: img, M_width: img.Rect.Dx(), M_height: img.Rect.Dy(), M_name: ""}

	gl.GenTextures(1, &texture.M_renderID)
	gl.BindTexture(gl.TEXTURE_2D, texture.M_renderID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// use to flip image
	imgData := FlipTexturePixels(img)
	// imgData := img.Pix
	//
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Size().X), int32(img.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(imgData))
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture

}

func FlipTexturePixels(img *image.RGBA) []byte {
	data := make([]byte, img.Rect.Dx()*img.Rect.Dy()*4)
	lineLen := img.Rect.Dx() * 4
	dest := len(data) - lineLen
	for src := 0; src < len(img.Pix); src += img.Stride {
		copy(data[dest:dest+lineLen], img.Pix[src:src+img.Stride])
		dest -= lineLen
	}

	return data
}

func (tex *Texture) DeleteTexture() {
	gl.DeleteTextures(1, &tex.M_renderID)
}
func (tex *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, tex.M_renderID)
}
func (tex *Texture) Bind2(renderID uint32) {
	// gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, renderID)
}
func (tex *Texture) UnBind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
