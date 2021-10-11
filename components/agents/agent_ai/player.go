package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/particles/particle_emitter"
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

	//these variables are for when we hold down the down key to ignore platforms
	holdTimer float32 = 0
	holdDelay float32 = 0.6
)

func AiHandlerPlayer(player *agents.Agent, ctx AiContext) {

	var inputXAxis float32
	//added this check to prevent player from moving when in "freecam" or other mode
	if input.GetInputContext() == input.GAME {
		inputXAxis = input.GetAxis(input.HORIZONTAL)

		if player.Transform.IsOnGround() && input.GetButtonDown("jump") {
			sound.PlaySound("player_jump")
			jumpFrame = 0
		}

		if jumpFrame < jumpFrames {
			player.Transform.Vel.Y += playerJumpAcceleration * constants.PHYSICS_TICK
			jumpFrame += 1
		} else if !player.Transform.IsOnGround() &&
			input.GetButton("jump") &&
			jumpFrame < jumpMaxFrames {
			player.Transform.Vel.Y += playerJumpAcceleration / 3 * constants.PHYSICS_TICK
			jumpFrame += 1
		}

		// fmt.Println(holdTimer)
		if input.GetButtonDown("down") {
			player.Meta.IgnoringPlatformsFor = ignorePlatformTime
		} else {

			if player.Meta.IgnoringPlatformsFor > 0 {
				player.Meta.IgnoringPlatformsFor -= constants.PHYSICS_TICK
			}
		}

		if input.GetButton("down") {
			holdTimer += constants.PHYSICS_TICK
			if holdTimer > holdDelay {
				player.Meta.IgnoringPlatformsFor += constants.PHYSICS_TICK
			}
		} else {
			holdTimer = 0
		}
	}

	// player.PhysicsState.Vel.X += playerWalkAccel * inputXAxis

	if !player.Transform.IsOnGround() {

		player.Transform.Vel.X += playerWalkAccel * inputXAxis
	} else {

		player.PhysicsState.Vel.X += playerWalkAccel * inputXAxis

	}

	if inputXAxis != 0 {
		player.PhysicsState.Direction = math32.Sign(inputXAxis)

		if player.PhysicsState.IsOnGround() {
			var dustPos = player.PhysicsState.Pos
			dustPos.Y = dustPos.Y - 1.5
			particle_emitter.EmitDust(dustPos)
		}


	}

	friction :=
		cxmath.Sign(player.Transform.Vel.X) * frictionFactor

	if math32.Abs(friction) < math32.Abs(player.Transform.Vel.X) {
		player.Transform.Vel.X -= friction
	} else {
		player.Transform.Vel.X = 0
	}

	player.Transform.Vel.X = math32.Clamp(player.Transform.Vel.X, -maxPlayerWalkSpeed, maxPlayerWalkSpeed)
	// if math32.Abs(player.PhysicsState.Vel.X) > maxPlayerWalkSpeed {
	// 	player.PhysicsState.Vel.X =
	// 		math32.Sign(player.PhysicsState.Vel.X) * maxPlayerWalkSpeed
	// }

	player.Transform.Vel.Y -= constants.Gravity * constants.PHYSICS_TICK

	player.Transform.IsIgnoringPlatforms = player.Meta.IgnoringPlatformsFor > 0
}
