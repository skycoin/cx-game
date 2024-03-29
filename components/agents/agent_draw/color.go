package agent_draw

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/render"
)

var debugAgentColor = mgl32.Vec4{0, 1, 0, 1}

// for debugging
func ColorDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	for _, agent := range agents {
		translate := mgl32.Translate3D(
			agent.Transform.Pos.X,
			agent.Transform.Pos.Y,
			0)
		scale := mgl32.Scale3D(
			agent.Transform.Size.X,
			agent.Transform.Size.Y,
			1)
		modelView := ctx.Camera.GetViewMatrix().Mul4(translate.Mul4(scale))
		render.DrawColorQuad(modelView, debugAgentColor)

	}
}
