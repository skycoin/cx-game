package pathfinding

import (
	"github.com/skycoin/cx-game/cxmath/math32"
)

type LeapingBehaviour struct {
	verticalJumpSpeed   float32
	horizontalJumpSpeed float32
}

func (lb LeapingBehaviour) Follow(ctx BehaviourContext) {
	// TODO fix this
	// dt := float32(1.0 / 30.0)
	directionX := math32.Sign(ctx.PlayerPos.X() - ctx.Self.Pos.X)

	if ctx.Self.Collisions.Below {
		ctx.Self.Vel.X = directionX * lb.horizontalJumpSpeed
		ctx.Self.Vel.Y = lb.verticalJumpSpeed
	} else {
		// ctx.Self.Vel.Y -= physics.Gravity * dt
	}
}

var LeapingBehaviourID BehaviourID

func init() {
	LeapingBehaviourID = RegisterBehaviour(LeapingBehaviour{
		horizontalJumpSpeed: 5,
		verticalJumpSpeed:   15,
	})
}
