package agent_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
)

func AnimatedDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	anim.Program.Use()
	defer anim.Program.StopUsing()

	gl.Enable(gl.DEPTH_TEST)

	gl.BindVertexArray(render.QuadVao)

	for _, agent := range agents {
		tex := agent.AnimationPlayback.Animation.Texture
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tex)

		alpha := timer.GetTimeBetweenTicks() / constants.MS_PER_TICK

		var interpolatedPos cxmath.Vec2
		if !agent.Transform.PrevPos.Equal(agent.Transform.Pos) {
			interpolatedPos = agent.Transform.PrevPos.Mult(1 - alpha).Add(agent.Transform.Pos.Mult(alpha))

		} else {
			interpolatedPos = agent.Transform.Pos
		}
		translate := mgl32.Translate3D(
			interpolatedPos.X,
			interpolatedPos.Y,
			constants.AGENT_Z,
		)
		scale := mgl32.Scale3D(
			agent.Transform.Size.X*agent.Transform.Direction,
			agent.Transform.Size.Y,
			1,
		)
		transform := translate.Mul4(scale)
		wrappedTransform := wrapTransform(
			transform,
			ctx.Camera.PlanetWidth,
			ctx.Camera.GetTransform(),
		)
		projection := spriteloader.Window.GetProjectionMatrix()
		mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

		anim.Program.SetMat4("mvp", &mvp)
		texTransform := agent.AnimationPlayback.Frame().Transform
		anim.Program.SetMat3("texTransform", &texTransform)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}
}

func wrapTransform(raw mgl32.Mat4, worldWidth float32, cameraTransform mgl32.Mat4) mgl32.Mat4 {
	rawX := raw.At(0, 3)
	x := math32.PositiveModulo(rawX, worldWidth)
	camX := cameraTransform.At(0, 3)
	if x-camX > worldWidth/2 {
		x -= worldWidth
	}
	if x-camX < -worldWidth/2 {
		x += worldWidth
	}

	translate := mgl32.Translate3D(x-rawX, 0, 0)
	return translate.Mul4(raw)
}
