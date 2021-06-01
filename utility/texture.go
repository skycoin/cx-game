package utility

// Helper texture generator functions
import "github.com/go-gl/gl/v4.1-core/gl"

//generate 1d_texture provided filename
func Make1DTexture(fname string) uint32 {
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_1D, tex)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_1D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	img := LoadPng(fname)
	gl.TexImage1D(gl.TEXTURE_1D, 0, gl.RGBA, int32(img.Rect.Size().X), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return tex
}

//generate 2d_texture provided filename
func Make2DTexture(fname string) uint32 {
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	img := LoadPng(fname)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Size().X), int32(img.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return tex
}
