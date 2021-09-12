package agent_ai

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
)

const (
	verticalJumpSpeed     float32 = 8
	attackSpeed           float32 = 5
	secondsBetweenLeaps   float32 = 2
	distanceBetweenPlayer float32 = 8
	distanceBeforeJump    float32 = 7
	collissionDistanceVal float32 = 1
)

func slimeIdle(slimePhysicsState physics.Body, playback *anim.Playback) {
	playback.PlayRepeating("Idle")
	slimePhysicsState.Vel.X = 0
}

var (
	isAttacking         = false
	isBack              = false
	playerPosY  float32 = verticalJumpSpeed
	oldDistance float32 = 0
)

func slimeAttack(distance float32, directionX float32, slimePhysicsState *physics.Body, playback *anim.Playback, agentIsWaiting bool, playerPos mgl32.Vec2) {
	playback.PlayRepeating("Pre")
	canAttack := slimePhysicsState.Collisions.Below && !agentIsWaiting && (distance <= distanceBeforeJump)
	if canAttack || isAttacking {
		playback.PlayRepeating("Collide")
		isAttacking = true
		playerPosY -= 1
		currentDistance := playerPos.X() - slimePhysicsState.Pos.X
		if currentDistance < 0 {
			currentDistance = currentDistance * -1
		}
		oldDistance = directionX * (currentDistance * attackSpeed)
		slimePhysicsState.Vel.X = oldDistance
		slimePhysicsState.Vel.Y = playerPosY
		fmt.Println("jump: ", currentDistance <= collissionDistanceVal, " -> currentDistance ", currentDistance, " ", collissionDistanceVal, " playerPOSY ", playerPosY)
		if currentDistance <= collissionDistanceVal {
			// hit player
			playback.PlayRepeating("Collide")
			slimePhysicsState.Vel.Y = verticalJumpSpeed
			slimePhysicsState.Vel.X = 0
			isBack = true
		}

		if isBack {
			playback.PlayOnce("Fall")
			slimePhysicsState.Vel.Y = 0
			slimePhysicsState.Vel.X = (distanceBeforeJump) * (directionX * -1)
			if currentDistance >= distanceBeforeJump {
				playback.PlayRepeating("Pre")
				slimePhysicsState.Vel.X = 0
				slimePhysicsState.Vel.Y = 0
				// reset from the beginning attack
				playerPosY = verticalJumpSpeed
				isBack = false
				isAttacking = false
			}
		}
	} else {
		playback.PlayRepeating("Pre")
		slimePhysicsState.Vel.X = 0
	}
}

func AiHandlerLeap(agent *agents.Agent, ctx AiContext) {
	distance := ctx.PlayerPos.X() - agent.PhysicsState.Pos.X

	directionX := math32.Sign(distance)
	if math32.Abs(distance) > ctx.WorldWidth/2 {
		directionX *= -1
	}

	if distance <= distanceBetweenPlayer {
		slimeAttack(distance, directionX, &agent.PhysicsState, &agent.AnimationPlayback, agent.IsWaiting(), ctx.PlayerPos)
	} else {
		slimeIdle(agent.PhysicsState, &agent.AnimationPlayback)
	}
}
