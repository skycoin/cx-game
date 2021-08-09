package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
)

const (
	verticalJumpSpeed   float32 = 15
	horizontalJumpSpeed float32 = 5
	secondsBetweenLeaps float32 = 2
)

func AiHandlerLeap(agent *agents.Agent, ctx AiContext) {
	directionX := math32.Sign(ctx.PlayerPos.X() - agent.PhysicsState.Pos.X)

	onGround := agent.PhysicsState.Collisions.Below
	canJump := onGround && !agent.IsWaiting()
	if canJump {
		agent.PhysicsState.Vel.X = directionX * horizontalJumpSpeed
		agent.PhysicsState.Vel.Y = verticalJumpSpeed
		agent.WaitFor(secondsBetweenLeaps)
		agent.AnimationPlayback.PlayOnce("Jump")
	}

	// disable sliding
	if onGround && !canJump {
		agent.PhysicsState.Vel.X = 0
	}
}
