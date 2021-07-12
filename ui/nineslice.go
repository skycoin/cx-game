package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

func newStretchingNineSliceVao() (tex uint32, verts int32) {
	geometry := render.NewGeometry()

	// TODO read from struct / config file
	/*
	*/

	s := float32(100) // make large VAO to accomodate varying sizes
	geometry.AddQuadFromCorners(
		render.Vert { 0, 0,0, 0,0 },
		render.Vert { s,-s,0, s,s },
	)

	return geometry.Upload(),geometry.Verts()
}

// texture dimensions of nineslice
type NineSliceDims struct {
	Left,Right,Top,Bottom float32
}

type StretchingNineSlice struct {
	sprite spriteloader.SpriteID
	vao uint32
	verts int32
	dims NineSliceDims
	shader *utility.Shader
}

func NewStretchingNineSlice(
		sprite spriteloader.SpriteID, dims NineSliceDims,
) StretchingNineSlice {
	vao,verts := newStretchingNineSliceVao()
	return StretchingNineSlice {
		sprite: sprite, vao: vao, verts: verts, dims: dims,
		shader: utility.NewShader(
			"./assets/shader/nineslice.vert",
			"./assets/shader/nineslice.frag",
		),
	}
}

func (nine StretchingNineSlice) Draw(ctx render.Context, size mgl32.Vec2) {
	metadata := spriteloader.GetSpriteMetadata(nine.sprite)
	gl.ActiveTexture(gl.TEXTURE0)
	nine.shader.Use()
	defer nine.shader.StopUsing()

	nine.shader.SetUint("tex", metadata.GpuTex)
	gl.BindTexture(gl.TEXTURE_2D, metadata.GpuTex)

	nine.shader.SetVec4F("color",1,1,1,1)
	nine.shader.SetVec4F("colour",1,1,1,1)
	mvp := ctx.MVP()
	nine.shader.SetMat4("mvp", &mvp)

	nine.shader.SetVec2F("offset", 0,0)
	nine.shader.SetVec2F("scale", 1,1)

	left := float32(1.0/6.0)
	right := left
	top := float32(1.0/8.0)
	bottom := float32(2.0/8.0)

	nine.shader.SetFloat("left",left)
	nine.shader.SetFloat("right",right)
	nine.shader.SetFloat("top",top)
	nine.shader.SetFloat("bottom",bottom)

	nine.shader.SetFloat("width",size.X())
	nine.shader.SetFloat("height",size.Y())

	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(nine.vao)
	gl.DrawArrays(gl.TRIANGLES,0,nine.verts)
}
