package agent_draw

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

var debugAgentColor = mgl32.Vec4{0, 1, 0, 1}

// for debugging
func ColorDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	for _, agent := range agents {
		translate := mgl32.Translate3D(
			agent.PhysicsState.Pos.X-ctx.Camera.X,
			agent.PhysicsState.Pos.Y-ctx.Camera.Y,
			0)
		scale := mgl32.Scale3D(
			agent.PhysicsState.Size.X,
			agent.PhysicsState.Size.Y,
			1)
		ctx := render.Context{
			World:      translate.Mul4(scale),
			Projection: spriteloader.Window.GetProjectionMatrix(),
		}

		render.DrawColorQuad(ctx, debugAgentColor)

	}
}
