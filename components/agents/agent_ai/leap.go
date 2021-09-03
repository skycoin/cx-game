package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
)

const (
	verticalJumpSpeed     float32 = 15
	horizontalJumpSpeed   float32 = 5
	secondsBetweenLeaps   float32 = 2
	glideSpeed            float32 = 1
	distanceBetweenPlayer float32 = 8
	distanceBeforeJump    float32 = 7
)

func attack(distance float32, directionX float32, slimePhysicsState *physics.Body, playback *anim.Playback, agentIsWaiting bool) {
	// slime start to attack the player from these distance
	playback.PlayRepeating("Pre")
	onGround := slimePhysicsState.Collisions.Below
	canJump := onGround && !agentIsWaiting
	if distance <= distanceBeforeJump {
		if canJump {
			playback.PlayRepeating("Jump")
			slimePhysicsState.Vel.X = directionX * horizontalJumpSpeed
			playback.PlayRepeating("Fall")
			slimePhysicsState.Vel.Y = verticalJumpSpeed
		}

	} else {
		slimePhysicsState.Vel.X = 0
		playback.PlayOnce("Pre")
	}

	if onGround && !canJump {
		slimePhysicsState.Vel.X = 0
	}

	// if !onGround && math32.Abs(slimePhysicsState.Vel.X) < glideSpeed {
	// 	slimePhysicsState.Vel.X = glideSpeed * directionX
	// }

}

func idle(slimePhysicsState physics.Body, playback *anim.Playback) {
	playback.PlayRepeating("Idle")
	slimePhysicsState.Vel.X = 0
}

func AiHandlerLeap(agent *agents.Agent, ctx AiContext) {
	distance := ctx.PlayerPos.X() - agent.PhysicsState.Pos.X

	directionX := math32.Sign(distance)
	if math32.Abs(distance) > ctx.WorldWidth/2 {
		directionX *= -1
	}

	if distance <= distanceBetweenPlayer {
		attack(distance, directionX, &agent.PhysicsState, &agent.AnimationPlayback, agent.IsWaiting())
	} else {
		idle(agent.PhysicsState, &agent.AnimationPlayback)
	}
}
