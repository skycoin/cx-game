package Texture

import (
	"image"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/skycoin/cx-game/engine/spriteloader"
)

type Texture struct {
	M_renderID       uint32
	M_filePath       string
	M_localBufferImg *image.RGBA
	M_width          int
	M_height         int
	M_BPP            int
}

func SetUpTexture(path string) *Texture {

	spriteloader.InitSpriteloaderDev()

	status, img, _ := spriteloader.LoadPng(path)
	if status != spriteloader.LoadOk {
		log.Fatalf("cannot upload [%v] to GPU", path)
	}

	texture := &Texture{M_renderID: 0, M_filePath: path, M_localBufferImg: img, M_width: img.Rect.Dx(), M_height: img.Rect.Dy()}
	gl.GenTextures(1, &texture.M_renderID)
	gl.BindTexture(gl.TEXTURE_2D, texture.M_renderID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Dx()), int32(img.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return texture

}

func (tex *Texture) DeleteTexture() {
	gl.DeleteTextures(1, &tex.M_renderID)
}
func (tex *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, tex.M_renderID)
}
func (tex *Texture) UnBind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
