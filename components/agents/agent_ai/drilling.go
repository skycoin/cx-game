package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/events"
)

const (
	drillSpeed     float32 = 3
	drillJumpSpeed float32 = 15
)

func AiHandlerDrill(agent *agents.Agent, ctx AiContext) {
	directionX := math32.Sign(ctx.PlayerPos.X() - agent.PhysicsState.Pos.X)
	agent.PhysicsState.Direction = directionX * -1
	agent.PhysicsState.Vel.X = directionX * drillSpeed

	doJump :=
		agent.PhysicsState.Collisions.Horizontal() &&
			agent.PhysicsState.IsOnGround()

	if doJump {
		events.OnSpiderBeforeJump.Trigger(events.SpiderEventData{
			WaitingFor: agent.WaitingFor,
		})

		agent.PhysicsState.Vel.Y = drillJumpSpeed
		// trigger an event when spiderdrill jump
		events.OnSpiderJump.Trigger(events.SpiderEventData{
			WaitingFor: agent.WaitingFor,
		})
	} else {
		agent.PhysicsState.Vel.Y -= constants.Gravity * constants.TimeStep
	}
}
