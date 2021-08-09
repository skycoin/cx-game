package models

import (
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/physics/movement"
	"github.com/skycoin/cx-game/world"
)

var (
	maxVerticalSpeed           float32 = 5
	minJumpSpeed, maxJumpSpeed float32
	wallSlideSpeed             float32 = -5
)

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
			/*
				if input.GetButtonUp("jump") {
					if player.Vel.Y > minJumpSpeed {
						player.Vel.Y = minJumpSpeed
					}
				}
			*/
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
	case movement.CLIMBING:
		if player.Collisions.Below {
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
