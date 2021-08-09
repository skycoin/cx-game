package agent_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/render"
)

func AnimatedDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	anim.Program.Use()
	defer anim.Program.StopUsing()

	gl.BindVertexArray(render.QuadVao)

	for _, agent := range agents {
		tex := agent.AnimationPlayback.Animation.Texture
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tex)
		translate := mgl32.Translate3D(
			agent.PhysicsState.Pos.X-ctx.Camera.X,
			agent.PhysicsState.Pos.Y-ctx.Camera.Y,
			0,
		)
		scale := mgl32.Scale3D(
			agent.PhysicsState.Size.X*agent.PhysicsState.Direction,
			agent.PhysicsState.Size.Y,
			1,
		)
		projection := spriteloader.Window.GetProjectionMatrix()
		mvp := projection.Mul4(translate).Mul4(scale)
		anim.Program.SetMat4("mvp", &mvp)
		texTransform := agent.AnimationPlayback.Frame().Transform
		anim.Program.SetMat3("texTransform", &texTransform)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}
}
