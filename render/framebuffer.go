package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/skycoin/cx-game/constants"
)

type Framebuffer uint32

func (f Framebuffer) Gl() uint32 { return uint32(f) }
func (f Framebuffer) Bind(target uint32) {
	gl.BindFramebuffer(target, f.Gl())
}
func (f Framebuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
func (f Framebuffer) SetTexture2D(texture uint32) {
	gl.FramebufferTexture2D(
		gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D,
		texture, 0,
	)
}

func (f Framebuffer) Attach(attachment uint32, texture uint32) {
	gl.FramebufferTexture2D(
		gl.FRAMEBUFFER, attachment, gl.TEXTURE_2D, texture, 0,
	)
}

func GenFramebuffer() Framebuffer {
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	return Framebuffer(fbo)
}

var (
	FRAMEBUFFER_SCREEN = Framebuffer(0)
	//world is drawn into this framebuffer
	FRAMEBUFFER_MAIN     Framebuffer
	FRAMEBUFFER_PLANET   Framebuffer
	RENDERTEXTURE_PLANET uint32
	RENDERTEXTURE_MAIN   uint32
)

func InitPlanetFrameBuffer() {
	fb := GenFramebuffer()
	fb.Bind(gl.FRAMEBUFFER)
	defer fb.Unbind()

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		constants.VIRTUAL_VIEWPORT_WIDTH, constants.VIRTUAL_VIEWPORT_HEIGHT,
		0, gl.RGBA, gl.UNSIGNED_BYTE, nil,
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	var depth uint32
	gl.GenTextures(1, &depth)
	gl.BindTexture(gl.TEXTURE_2D, depth)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.DEPTH24_STENCIL8,
		constants.VIRTUAL_VIEWPORT_WIDTH, constants.VIRTUAL_VIEWPORT_HEIGHT, 0,
		gl.DEPTH_STENCIL, gl.UNSIGNED_INT_24_8, nil,
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// unbind texture such that we can attach it to the framebuffer
	gl.BindTexture(gl.TEXTURE_2D, 0)
	fb.SetTexture2D(tex)
	fb.Attach(gl.DEPTH_STENCIL_ATTACHMENT, depth)

	FRAMEBUFFER_PLANET = fb
	RENDERTEXTURE_PLANET = tex
}

func InitMainFramebuffer() {
	fb := GenFramebuffer()
	fb.Bind(gl.FRAMEBUFFER)
	defer fb.Unbind()

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		constants.VIRTUAL_VIEWPORT_WIDTH, constants.VIRTUAL_VIEWPORT_HEIGHT,
		0, gl.RGBA, gl.UNSIGNED_BYTE, nil,
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	var depth uint32
	gl.GenTextures(1, &depth)
	gl.BindTexture(gl.TEXTURE_2D, depth)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.DEPTH24_STENCIL8,
		constants.VIRTUAL_VIEWPORT_WIDTH, constants.VIRTUAL_VIEWPORT_HEIGHT, 0,
		gl.DEPTH_STENCIL, gl.UNSIGNED_INT_24_8, nil,
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// unbind texture such that we can attach it to the framebuffer
	gl.BindTexture(gl.TEXTURE_2D, 0)
	fb.SetTexture2D(tex)
	fb.Attach(gl.DEPTH_STENCIL_ATTACHMENT, depth)

	FRAMEBUFFER_MAIN = fb
	RENDERTEXTURE_MAIN = tex
}

func BlitFramebuffer(
	readFramebuffer, drawFramebuffer Framebuffer,
	srcX0, srcY0, srcX1, srcY1 int32,
	dstX0, dstY0, dstX1, dstY1 int32,
	mask uint32, filter uint32,
) {
	checkError("before bind read")
	readFramebuffer.Bind(gl.READ_FRAMEBUFFER)
	checkError("after bind read")
	drawFramebuffer.Bind(gl.DRAW_FRAMEBUFFER)
	checkError("after bind draw")

	gl.BlitFramebuffer(
		srcX0, srcY0, srcX1, srcY1,
		dstX0, dstY0, dstX1, dstY1,
		mask, filter,
	)

	checkError("after blit")
}

func checkError(prefix string) {
	errorCode := gl.GetError()
	if errorCode != gl.NO_ERROR {
		log.Printf("%v: got GL error code [%v]", prefix, errorCode)
	}
}
