package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

func newStretchingNineSliceVao() (tex uint32, verts int32) {
	geometry := render.NewGeometry()

	// TODO read from struct / config file
	/*
	 */

	s := float32(100) // make large VAO to accomodate varying sizes
	geometry.AddQuadFromCorners(
		render.Vert{0, 0, 0, 0, 0},
		render.Vert{s, -s, 0, s, s},
	)

	return geometry.Upload(), geometry.Verts()
}

// texture dimensions of nineslice
type NineSliceDims struct {
	Left, Right, Top, Bottom float32
}

type StretchingNineSlice struct {
	sprite  spriteloader.SpriteID
	vao     uint32
	verts   int32
	dims    NineSliceDims
	program render.Program
}

func NewStretchingNineSlice(
	sprite spriteloader.SpriteID, dims NineSliceDims,
) StretchingNineSlice {
	vao, verts := newStretchingNineSliceVao()
	return StretchingNineSlice{
		sprite: sprite, vao: vao, verts: verts, dims: dims,
		program: render.CompileProgram(
			"./assets/shader/nineslice.vert",
			"./assets/shader/nineslice.frag",
		),
	}
}

func (nine StretchingNineSlice) Draw(ctx render.Context, size mgl32.Vec2) {
	metadata := spriteloader.GetSpriteMetadata(nine.sprite)
	gl.ActiveTexture(gl.TEXTURE0)
	nine.program.Use()
	defer nine.program.StopUsing()

	nine.program.SetUint("tex", metadata.GpuTex)
	gl.BindTexture(gl.TEXTURE_2D, metadata.GpuTex)

	nine.program.SetVec4F("color", 1, 1, 1, 1)
	nine.program.SetVec4F("colour", 1, 1, 1, 1)
	mvp := ctx.MVP()
	nine.program.SetMat4("mvp", &mvp)

	nine.program.SetVec2F("offset", 0, 0)
	nine.program.SetVec2F("scale", 1, 1)

	left := float32(1.0 / 6.0)
	right := left
	top := float32(1.0 / 8.0)
	bottom := float32(2.0 / 8.0)

	nine.program.SetFloat("left", left)
	nine.program.SetFloat("right", right)
	nine.program.SetFloat("top", top)
	nine.program.SetFloat("bottom", bottom)

	nine.program.SetFloat("width", size.X())
	nine.program.SetFloat("height", size.Y())

	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(nine.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, nine.verts)
}
