package pathfinding

import (
	"github.com/skycoin/cx-game/physics"
)

type FreeBehaviour struct {}

func (fb FreeBehaviour) Follow(ctx BehaviourContext) {
	// TODO fix this
	dt := float32(1.0/30.0)
	// apply gravity when not touching ground
	if !ctx.Self.Collisions.Below { ctx.Self.Vel.Y -= physics.Gravity * dt }
}

var FreeBehaviourID BehaviourID
func init() {
	FreeBehaviourID = RegisterBehaviour(FreeBehaviour{})
}
