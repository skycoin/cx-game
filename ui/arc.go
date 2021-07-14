package ui

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/spriteloader"
)

var arcShader *utility.Shader

var arcSprite spriteloader.SpriteID
func initArc() {
	arcShader = utility.NewShader(
		"./assets/shader/mvp.vert", "./assets/shader/arc.frag" )
	arcSprite = spriteloader.
		LoadSingleSprite("./assets/hud/hud_status_fill.png","hud_status_fill")
}

func DrawArc(mvp mgl32.Mat4, fullness float32) {
	gl.ActiveTexture(gl.TEXTURE0)
	arcShader.Use()
	defer arcShader.StopUsing()
	metadata := spriteloader.GetSpriteMetadata(arcSprite)
	arcShader.SetUint("tex", metadata.GpuTex)
	gl.BindTexture(gl.TEXTURE_2D, metadata.GpuTex)

	arcShader.SetMat4("mvp",&mvp)
	arcShader.SetVec4F("color",1,1,1,1)
	// FIXME note that this only works because sprite is in its own image.
	// Update this with a more robust strategy at some point.
	arcShader.SetVec2F("offset",0,0)
	arcShader.SetVec2F("scale",1,1)

	arcShader.SetFloat("value",2*math.Pi*fullness)

	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(spriteloader.QuadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
