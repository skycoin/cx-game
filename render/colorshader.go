// super simple shader for drawing colors directly.
// intended for UI.
package render;
import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func NewColorShader() Program {
	return CompileProgram(
		"./assets/shader/color.vert",
		"./assets/shader/color.frag",
	)
}
var colorProgram Program
var colorProgramInit bool = false

func DrawColorQuad(ctx Context, colour mgl32.Vec4) {
	if !colorProgramInit {
		colorProgram = NewColorShader()
		colorProgramInit = true
	}
	// setup features
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// update uniforms
	colorProgram.Use()
	defer colorProgram.StopUsing()
	colorProgram.SetVec4("colour",&colour)
	mvp := ctx.MVP()
	colorProgram.SetMat4("mvp",&mvp)
	// draw
	gl.BindVertexArray(QuadVao)
	gl.DrawArrays(gl.TRIANGLES,0,6)
}
