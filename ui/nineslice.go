package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

func newStretchingNineSliceVao(w,h float32) (tex uint32, verts int32) {
	geometry := utility.NewGeometry()

	// TODO read from struct / config file
	left := float32(1.0/6.0)
	right := left
	top := float32(1.0/8.0)
	bottom := float32(2.0/8.0)

	geometry.AddQuadFromCorners( // top left
		utility.Vert { 0,0,0, 0,0 },
		utility.Vert { left,-top,0, left,top },
	)
	geometry.AddQuadFromCorners( // top middle
		utility.Vert { left,0,0, left,0 },
		utility.Vert { w-right,-top,0, 1-right,top },
	)
	geometry.AddQuadFromCorners( // top righ
		utility.Vert { w-right,0,0, 1-right,0 },
		utility.Vert { w,-top,0, 1,top },
	)
	
	geometry.AddQuadFromCorners( // middle left
		utility.Vert { 0,-top,0, 0,top },
		utility.Vert { left,-h+bottom,0, left,1-bottom },
	)
	geometry.AddQuadFromCorners( // center
		utility.Vert { left,-top,0, left,top },
		utility.Vert { w-right, -h+bottom,0, 1-right,1-bottom },
	)
	geometry.AddQuadFromCorners( //middle right
		utility.Vert { w-right, -top, 0, 1-right, top },
		utility.Vert { w, -h+bottom, 0, 1, 1-bottom },
	)

	geometry.AddQuadFromCorners( // bottom left
		utility.Vert { 0, -h+bottom, 0, 0,1-bottom },
		utility.Vert { left, -h, 0, 0, 1 },
	)
	geometry.AddQuadFromCorners( // bottom middle
		utility.Vert { left,-h+bottom, 0, left,1-bottom },
		utility.Vert { w-right, -h, 0, 1-right, 1 },
	)
	geometry.AddQuadFromCorners( // bottom right
		utility.Vert { w-right, -h+bottom, 0, 1-right, 1-bottom },
		utility.Vert { w,-h,0, 1,1},
	)

	return geometry.Upload(),geometry.Verts()
}

type StretchingNineSlice struct {
	Width float32
	sprite spriteloader.SpriteID
	vao uint32
	verts int32
	shader *utility.Shader
}

func NewStretchingNineSlice(
		sprite spriteloader.SpriteID, w,h float32,
) StretchingNineSlice {
	vao,verts := newStretchingNineSliceVao(w,h)
	return StretchingNineSlice {
		Width: w,
		sprite: sprite, vao: vao, verts: verts,
		shader: utility.NewShader(
			"./assets/shader/mvp.vert", "./assets/shader/tex.frag" ),
	}
}

func (nine StretchingNineSlice) Draw(ctx render.Context) {
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

	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(nine.vao)
	gl.DrawArrays(gl.TRIANGLES,0,nine.verts)
}
