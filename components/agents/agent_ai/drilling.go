package agent_ai

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
)

const (
	drillSpeed float32 = 1
)

func AiHandlerDrill(agent *agents.Agent, ctx AiContext) {
	directionX := math32.Sign(ctx.PlayerPos.X() - agent.PhysicsState.Pos.X)
	agent.PhysicsState.Vel.X = directionX * walkSpeed
}
