package physics

import (
	"github.com/skycoin/cx-game/world/worldcollider"
)

var bodies = []*Body{}

func RegisterBody(body *Body) {
	bodies = append(bodies, body)
}

type TickInfo struct {
	Dt            float32
	WorldCollider worldcollider.WorldCollider
}

var prevTime, nextTime = float32(0), float32(0)
var now = float32(0)



// will a physics tick run after this time delta?
// use to decide whether we consume keyboard inputs etc.
func WillTick(dt float32) bool {
	return now+dt > nextTime
}

// run the necessary number of physics ticks for a given time delta
func Simulate(dt float32, worldcollider worldcollider.WorldCollider) {
	now += dt
	// run physics ticks until the current time lies between
	// the previous physics state and the next physics state.
	// Then, we can interpolate
	for nextTime < now {
		tick(worldcollider)
	}

	// physics simulation is done; save interpolated values for rendering
	alpha := (now - prevTime) / TimeStep
	for idx, _ := range bodies {
		bodies[idx].UpdateInterpolatedTransform(alpha)
	}
}

func tick(worldcollider worldcollider.WorldCollider) {
	prevTime = nextTime
	nextTime += TimeStep

	newBodies := []*Body{}
	for _, body := range bodies {
		body.SavePreviousTransform()
		body.Tick(worldcollider, TimeStep)
		if !body.Deleted {
			newBodies = append(newBodies, body)
		}
	}
	bodies = newBodies
}
