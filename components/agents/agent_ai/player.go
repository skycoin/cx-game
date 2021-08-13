package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/input"
)

const (
	ignorePlatformTime float32 = 0.4
	playerWalkAccel    float32 = 5
	maxPlayerWalkSpeed float32 = 8
	// playerJumpSpeed    float32 = 25
	playerJumpAcceleration float32 = 120
	jumpMaxFrames                  = 20
	frictionFactor         float32 = 2
)

var (
	jumpFrame int = 5
)

func AiHandlerPlayer(player *agents.Agent, ctx AiContext) {
	var inputXAxis float32
	if input.GetInputContext() == input.GAME {
		inputXAxis = input.GetAxis(input.HORIZONTAL)
	}
	player.PhysicsState.Vel.X +=
		inputXAxis * playerWalkAccel

	if inputXAxis != 0 {
		player.PhysicsState.Direction = math32.Sign(inputXAxis)
	}

	friction :=
		cxmath.Sign(player.PhysicsState.Vel.X) * frictionFactor

	if math32.Abs(friction) < math32.Abs(player.PhysicsState.Vel.X) {
		player.PhysicsState.Vel.X -= friction
	} else {
		player.PhysicsState.Vel.X = 0
	}

	if math32.Abs(player.PhysicsState.Vel.X) > maxPlayerWalkSpeed {
		player.PhysicsState.Vel.X =
			math32.Sign(player.PhysicsState.Vel.X) * maxPlayerWalkSpeed
	}

	player.PhysicsState.Vel.Y -= constants.Gravity * constants.TimeStep

	if player.PhysicsState.IsOnGround() && input.GetButtonDown("jump") {
		jumpFrame = 0
	}
	if jumpFrame < 7 {
		player.PhysicsState.Vel.Y += playerJumpAcceleration * constants.TimeStep
		jumpFrame += 1
	}
	if input.GetButton("jump") && jumpFrame < jumpMaxFrames {
		player.PhysicsState.Vel.Y += playerJumpAcceleration / 3 * constants.TimeStep
		jumpFrame += 1
	}

	if input.GetButton("down") {
		player.PlayerData.IgnoringPlatformsFor = ignorePlatformTime
	} else {
		if player.PlayerData.IgnoringPlatformsFor > 0 {
			player.PlayerData.IgnoringPlatformsFor -= constants.TimeStep
		}
	}
	player.PhysicsState.IsIgnoringPlatforms = player.PlayerData.IgnoringPlatformsFor > 0
}
