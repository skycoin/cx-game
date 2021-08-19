package physics

import (
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/world/worldcollider"
)

var bodies = []*Body{}

func RegisterBody(body *Body) {
	bodies = append(bodies, body)
}

var Gravity float32 = 30.0

// length of physics tick in seconds.
// physics rate is independent of frame rate
const TimeStep float32 = 1.0 / 10

type TickInfo struct {
	Dt            float32
	WorldCollider worldcollider.WorldCollider
}

// run the necessary number of physics ticks for a given time delta
func Simulate(dt float32, worldcollider worldcollider.WorldCollider) {

	// run physics ticks until the current time lies between
	// the previous physics state and the next physics state.
	// Then, we can interpolate

	tick(worldcollider)

	// physics simulation is done; save interpolated values for rendering
	alpha := timer.GetTimeBetweenTicks() / constants.PHYSICS_TICK
	for idx, _ := range bodies {
		bodies[idx].UpdateInterpolatedTransform(alpha)
	}
}

func tick(worldcollider worldcollider.WorldCollider) {

	newBodies := []*Body{}
	for _, body := range bodies {
		body.SavePreviousTransform()
		body.Move(worldcollider, constants.PHYSICS_TICK)
		if !body.Deleted {
			newBodies = append(newBodies, body)
		}
	}
	bodies = newBodies
}
