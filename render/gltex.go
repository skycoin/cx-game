package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// wrapper over an opengl texture

type Texture struct {
	Target  uint32
	Texture uint32
}

// may need to adjust this interface if we ever need simultaneous texture units

func (t Texture) Bind() {
	gl.BindTexture(t.Target, t.Texture)
}

func (t Texture) Unbind() {
	gl.BindTexture(t.Target, 0)
}

func (t Texture) SetTextureFiltering(minFilter, magFilter int32) {
	// assume we already bind the texture
	gl.TexParameteri(t.Target, gl.TEXTURE_MIN_FILTER, minFilter)
	gl.TexParameteri(t.Target, gl.TEXTURE_MAG_FILTER, magFilter)
}
