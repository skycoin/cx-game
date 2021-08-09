package agent_draw

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

var debugAgentColor = mgl32.Vec4{0, 1, 0, 1}

// for debugging
func ColorDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	if len(agents) == 0 {
		return
	}
	spriteID := getSpriteID(agents[0].AgentCategory)
	drawOpts := spriteloader.NewDrawOptions()
	for _, agent := range agents {
		drawOpts.Alpha = alphaForAgent(agent)
		spriteloader.DrawSpriteQuadOptions(
			agent.PhysicsState.Pos.X-ctx.Camera.X,
			agent.PhysicsState.Pos.Y-ctx.Camera.Y,
			agent.PhysicsState.Size.X,
			agent.PhysicsState.Size.Y,
			spriteID, drawOpts,
		)
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
