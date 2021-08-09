package agent_draw

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/spriteloader"
)

const TimeBeforeFadeout = float32(1.0) // in seconds
const TimeDuringFadeout = float32(1.0) // in seconds

func alphaForAgent(agent *agents.Agent) float32 {
	if agent.TimeSinceDeath < TimeBeforeFadeout {
		return 1
	}
	x := agent.TimeSinceDeath - TimeBeforeFadeout
	return 1 - x/TimeDuringFadeout
}

func QuadDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	// TODO is this assumed??? can we omit this check?
	if len(agents) == 0 {
		return
	}
	spriteID := getSpriteID(agents[0].AgentCategory)
	drawOpts := spriteloader.NewDrawOptions()
	for _, agent := range agents {
		drawOpts.Alpha = alphaForAgent(agent)
		body := &agent.PhysicsState
		translate := mgl32.Translate3D(
			body.Pos.X-ctx.Camera.X,
			body.Pos.Y-ctx.Camera.Y,
			0,
		)
		scale := mgl32.Scale3D(body.Size.X,body.Size.Y,1)
		ctx := render.Context {
			World: translate.Mul4(scale),
			Projection: spriteloader.Window.GetProjectionMatrix(),
		}
		spriteloader.DrawSpriteQuadContext(ctx, spriteID, drawOpts)
	}
}

func getSpriteID(agentType constants.AgentCategory) spriteloader.SpriteID {
	switch agentType {
	default:
		return spriteloader.GetSpriteIdByName("basic-agent")
	}
}
