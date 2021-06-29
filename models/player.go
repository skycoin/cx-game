package models

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/physics/movement"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/world"
)

type Player struct {
	physics.Body
	movement.MovementComponent
	Controlled bool
	// RGBA            *image.RGBA
	// ImageSize       image.Point
	spriteId   int
	XDirection float32 // 1 when facing right, -1 when facing left
}

var (
	maxVerticalSpeed           float32 = 5
	minJumpSpeed, maxJumpSpeed float32
	wallSlideSpeed             float32 = -5
)

func NewPlayer() *Player {
	spriteId := spriteloader.LoadSingleSprite(
		"./assets/character/character.png", "player")
	player := Player{
		Body: physics.Body{
			Size: physics.Vec2{X: 2, Y: 3},
		},
		MovementComponent: movement.NewPlayerMovementComponent(),
		spriteId:          spriteId,
		XDirection:        1, // start facing right
		Controlled:        true,
	}
	maxJumpSpeed = cxmath.Sqrt(2 * cxmath.Abs(physics.Gravity) * player.MovementMeta.Jumpheight)
	minJumpSpeed = maxJumpSpeed / 4
	physics.RegisterBody(&player.Body)

	return &player
}

func (player *Player) Draw(cam *camera.Camera, planet *world.Planet) {
	disp := planet.ShortestDisplacement(
		mgl32.Vec2{cam.X, cam.Y},
		player.InterpolatedTransform.Col(3).Vec2())

	spriteloader.DrawSpriteQuad(
		disp.X(), disp.Y(),
		// player sprite actually faces left so throw an extra (-) here
		-player.Size.X*player.XDirection, player.Size.Y, player.spriteId,
	)
}

func (player *Player) FixedTick(planet *world.Planet) {
	//todo separate more logic
	//
	player.MovementBeforeTick()

	if player.Controlled {
		inputXAxis := input.GetAxis(input.HORIZONTAL)
		player.Vel.X += inputXAxis * player.MovementMeta.Acceleration * player.ActiveMovementType.GetMovementSpeedModifier()
		if inputXAxis != 0 {
			player.XDirection = math32.Sign(inputXAxis)
		}

		if player.ActiveMovementType == movement.FLYING {
			inputYAxis := input.GetAxis(input.VERTICAL)
			player.Vel.Y = inputYAxis * maxVerticalSpeed
		}
	}
	player.Vel.Y -= physics.Gravity * physics.TimeStep * player.ActiveMovementType.GetGravityModifier()

	if player.Vel.X != 0 {
		friction := cxmath.Sign(player.Vel.X) * player.MovementMeta.Acceleration * player.MovementMeta.DynamicFriction * player.ActiveMovementType.GetFrictionModifier()

		//to stop player from jiggling
		if cxmath.Abs(player.Vel.X) <= player.MovementMeta.Acceleration*player.MovementMeta.DynamicFriction*player.ActiveMovementType.GetFrictionModifier() && input.GetAxis(input.HORIZONTAL) == 0 {
			player.Vel.X = 0
		} else {
			player.Vel.X -= friction
			// fmt.Println(player.Vel.X)
		}
	}
	fmt.Println(player.Vel.X)
	player.Vel.X = utility.ClampF(player.Vel.X, -player.MovementMeta.MovSpeed*player.ActiveMovementType.GetMovementSpeedModifier(), player.MovementMeta.MovSpeed*player.ActiveMovementType.GetMovementSpeedModifier())
	player.MovementAfterTick(planet)
}

func (player *Player) MovementBeforeTick() {
	//handle movement state
	//todo change from crude switch state to more sophisticated and easier to read structure
	switch player.ActiveMovementType {
	case movement.NORMAL:
		if player.Collisions.Below {
			if input.GetButtonDown("crouch") {
				player.TryChangeMovementState(movement.CROUCHING)

			} else if input.GetButtonDown("action") && player.IsMovementTypePresent(movement.CLIMBING) {
				player.TryChangeMovementState(movement.CLIMBING)

			} else if input.GetButtonDown("fly") && player.IsMovementTypePresent(movement.FLYING) { // todo only from falling state
				player.TryChangeMovementState(movement.FLYING)
			}
		} else {
			if player.Collisions.Left || player.Collisions.Right { // && player.IsMovementTypePresent(movement.WALL_SLIDING) {
				player.TryChangeMovementState(movement.WALL_SLIDING)
				return
			}
			if input.GetButtonUp("jump") {
				if player.Vel.Y > minJumpSpeed {
					player.Vel.Y = minJumpSpeed
				}
			}
		}

	// handling crouching
	case movement.CROUCHING:
		if input.GetButtonDown("crouch") {
			player.TryChangeMovementState(movement.NORMAL)
		}
	// handling flying
	case movement.FLYING:
		if player.Collisions.Below || input.GetButtonDown("fly") {
			player.TryChangeMovementState(movement.NORMAL)
		}
	case movement.WALL_SLIDING:
		if player.Collisions.Below || input.GetButtonDown("jump") {
			player.TryChangeMovementState(movement.NORMAL)
		}

		if player.Collisions.Below {
			//change if only the previous position is not crouching or normal
			if player.PreviousActiveMovementType&movement.CROUCHING == 0 {
				player.ActiveMovementType = movement.NORMAL
			}
		} else if player.Collisions.Left || player.Collisions.Right {
			player.ActiveMovementType = movement.WALL_SLIDING
		} else {
			if player.ActiveMovementType == movement.WALL_SLIDING {
				player.TryChangeMovementState(movement.NORMAL)
			}
		}

		if !player.Collisions.Below && input.GetButtonUp("jump") {
			if player.Vel.Y > minJumpSpeed {
				player.Vel.Y = minJumpSpeed
			}
		}
	}
}

func (player *Player) Jump() bool {
	//jump only if there is jumps left or on ground or wall sliding
	if !player.CanJump() {
		return false
	}

	if player.ActiveMovementType == movement.WALL_SLIDING && player.IsMovementTypePresent(movement.CAN_WALL_JUMP) {
		player.Vel.X = cxmath.Sign(input.GetAxis(input.HORIZONTAL)) * -1 * 15
	}

	player.Vel.Y = maxJumpSpeed
	return true

	// if player.Collisions.Below {
	// 	// fmt.Println(jumpSpeed)

	// 	// jumpCounter = maxAdditionalJumps
	// 	return true
	// } else if player.Collisions.Left || player.Collisions.Right { //wall_jumping
	// 	player.Vel.Y = maxJumpSpeed
	// 	// player.Vel.X = cxmath.Min(-player.Vel.X, 15*cxmath.Sign(player.Vel.X))
	// 	player.Vel.X = cxmath.Sign(input.GetAxis(input.HORIZONTAL)) * -1 * 15
	// 	// jumpCounter = maxAdditionalJumps
	// 	return true
	// }

	// if jumpCounter > 0 {
	// 	jumpCounter -= 1
	// 	player.Vel.Y = maxJumpSpeed
	// 	return true
	// }
	// return false
}

//states - running
func (player *Player) MovementAfterTick(planet *world.Planet) {
	switch player.ActiveMovementType {
	case movement.NORMAL: // moving | idle
		if player.Collisions.Below {
			player.ResetJumpCounter()
		}
		//sprite normal
	case movement.WALL_SLIDING: // wall slide
		player.Vel.Y = cxmath.Max(player.Vel.Y, wallSlideSpeed)
		if player.IsMovementTypePresent(movement.CAN_WALL_JUMP) {
			player.ResetJumpCounter()
		}
	case movement.FLYING: // flying

		//todo

		// if player.PreviousActiveMovementType == movement.NORMAL {
		// 	if planet.GetDistanceFromGround(player.Pos) > 5 {
		// 		player.Vel.Y = maxVerticalSpeed
		// 	}
		// }
	}
}
