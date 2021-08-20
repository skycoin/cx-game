package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/sound"
)

const (
	ignorePlatformTime float32 = 0.15
	playerWalkAccel    float32 = 1.8
	maxPlayerWalkSpeed float32 = 8
	// playerJumpSpeed    float32 = 25
	playerJumpAcceleration float32 = 170
	jumpMaxFrames                  = 15
	jumpFrames                     = 3
	frictionFactor         float32 = 1
)

var (
	//give it initial big value so it will not jump immediately, try setting this to 0 and see result
	jumpFrame int = 10
)

func AiHandlerPlayer(player *agents.Agent, ctx AiContext) {

	var inputXAxis float32
	//added this check to prevent player from moving when in "freecam" or other mode
	if input.GetInputContext() == input.GAME {
		inputXAxis = input.GetAxis(input.HORIZONTAL)

		if player.PhysicsState.IsOnGround() && input.GetButtonDown("jump") {
			sound.PlaySound("player_jump")
			jumpFrame = 0
		}

		if jumpFrame < jumpFrames {
			player.PhysicsState.Vel.Y += playerJumpAcceleration * constants.PHYSICS_TICK
			jumpFrame += 1
		} else if !player.PhysicsState.IsOnGround() &&
			input.GetButton("jump") &&
			jumpFrame < jumpMaxFrames {
			player.PhysicsState.Vel.Y += playerJumpAcceleration / 3 * constants.PHYSICS_TICK
			jumpFrame += 1
		}

		if input.GetButtonDown("down") {
			player.PlayerData.IgnoringPlatformsFor = ignorePlatformTime
		} else {
			if player.PlayerData.IgnoringPlatformsFor > 0 {
				player.PlayerData.IgnoringPlatformsFor -= constants.PHYSICS_TICK
			}
		}
	}

	// player.PhysicsState.Vel.X += playerWalkAccel * inputXAxis

	if !player.PhysicsState.IsOnGround() {

		player.PhysicsState.Vel.X += playerWalkAccel * inputXAxis
	} else {

		player.PhysicsState.Vel.X += playerWalkAccel * inputXAxis
	}

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

	player.PhysicsState.Vel.X = math32.Clamp(player.PhysicsState.Vel.X, -maxPlayerWalkSpeed, maxPlayerWalkSpeed)
	// if math32.Abs(player.PhysicsState.Vel.X) > maxPlayerWalkSpeed {
	// 	player.PhysicsState.Vel.X =
	// 		math32.Sign(player.PhysicsState.Vel.X) * maxPlayerWalkSpeed
	// }

	player.PhysicsState.Vel.Y -= constants.Gravity * constants.PHYSICS_TICK

	player.PhysicsState.IsIgnoringPlatforms = player.PlayerData.IgnoringPlatformsFor > 0
}
