package agent_draw

import (
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
)


const TimeBeforeFadeout = float32(1.0) // in seconds
const TimeDuringFadeout = float32(1.0) // in seconds

func alphaForAgent(agent *agents.Agent) float32 {
	if agent.TimeSinceDeath < TimeBeforeFadeout { return 1 }
	x := agent.TimeSinceDeath - TimeBeforeFadeout
	return 1 - x / TimeDuringFadeout
}

func QuadDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	// TODO is this assumed??? can we omit this check?
	if len(agents)==0 { return }
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
	}
}

func getSpriteID(agentType constants.AgentCategory) spriteloader.SpriteID {
	switch agentType {
	default:
		return spriteloader.GetSpriteIdByName("basic-agent")
	}
}
