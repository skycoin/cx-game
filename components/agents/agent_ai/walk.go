package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
)

const (
	walkSpeed float32 = 1
	jumpSpeed float32 = 15
)

/*
func shouldJump(agent *agents.Agent) {
	// TODO
	return false
}
*/

func AiHandlerWalk(agent *agents.Agent, ctx AiContext) {
	directionX := math32.Sign(ctx.PlayerPos.X() - agent.Transform.Pos.X)
	agent.Transform.Vel.X = directionX * walkSpeed
}
