package pathfinding

import (
	"github.com/skycoin/cx-game/cxmath/math32"
)

type WalkingBehaviour struct {
	walkSpeed float32
	jumpSpeed float32
}

func (wb WalkingBehaviour) shouldJump(ctx BehaviourContext) bool {
	needsToJumpUpLeftWall := ctx.Self.Collisions.Left && ctx.Self.Vel.X < 0
	needsToJumpUpRightWall := ctx.Self.Collisions.Right && ctx.Self.Vel.X > 0
	needsToJump := needsToJumpUpLeftWall || needsToJumpUpRightWall
	canJump := ctx.Self.Collisions.Below
	return canJump && needsToJump
}

func (wb WalkingBehaviour) Follow(ctx BehaviourContext) {
	// TODO fix this
	// dt := float32(1.0 / 30.0)
	directionX := math32.Sign(ctx.PlayerPos.X() - ctx.Self.Pos.X)
	ctx.Self.Vel.X = directionX * wb.walkSpeed

	if wb.shouldJump(ctx) {
		ctx.Self.Vel.Y = wb.jumpSpeed
	} else {
		// ctx.Self.Vel.Y -= physics.Gravity * dt
	}
}

var WalkingBehaviourID BehaviourID

func init() {
	WalkingBehaviourID = RegisterBehaviour(WalkingBehaviour{
		walkSpeed: 1,
		jumpSpeed: 15,
	})
}
