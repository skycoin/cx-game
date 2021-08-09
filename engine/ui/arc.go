package ui

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

var arcProgram render.Program

var arcSprite spriteloader.SpriteID

func initArc() {
	arcProgram = render.CompileProgram(
		"./assets/shader/mvp.vert", "./assets/shader/arc.frag")
	arcSprite = spriteloader.
		LoadSingleSprite("./assets/hud/hud_status_fill.png", "hud_status_fill")
}

func DrawArc(mvp mgl32.Mat4, fullness float32) {
	gl.ActiveTexture(gl.TEXTURE0)
	arcProgram.Use()
	defer arcProgram.StopUsing()
	metadata := spriteloader.GetSpriteMetadata(arcSprite)
	arcProgram.SetUint("tex", metadata.GpuTex)
	gl.BindTexture(gl.TEXTURE_2D, metadata.GpuTex)

	arcProgram.SetMat4("mvp", &mvp)
	arcProgram.SetVec4F("color", 1, 1, 1, 1)
	// FIXME note that this only works because sprite is in its own image.
	// Update this with a more robust strategy at some point.
	arcProgram.SetVec2F("offset", 0, 0)
	arcProgram.SetVec2F("scale", 1, 1)

	arcProgram.SetFloat("value", 2*math.Pi*fullness)

	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(render.QuadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
