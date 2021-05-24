// super simple shader for drawing colors directly.
// intended for UI.
package utility;
import (
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func NewColorShader() *Shader {
	return NewShader(
		"./assets/shader/color.vert",
		"./assets/shader/color.frag",
	)
}
var colorShader *Shader

func DrawColorQuad(ctx render.Context, colour mgl32.Vec4) {
	if colorShader==nil {
		colorShader = NewColorShader()
	}
	// setup features
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// update uniforms
	gl.UseProgram(colorShader.ID)
	colorShader.SetVec4("colour",&colour)
	mvp := ctx.MVP()
	colorShader.SetMat4("mvp",&mvp)
	// draw
	gl.BindVertexArray(spriteloader.QuadVao)
	gl.DrawArrays(gl.TRIANGLES,0,6)
}
