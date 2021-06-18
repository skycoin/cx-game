package physics

import (
	"github.com/skycoin/cx-game/world"
)

var bodies = []*Body {}

func RegisterBody(body *Body) {
	bodies = append(bodies,body)
}

var Gravity float32 = 30.0
// length of physics tick in seconds.
// physics rate is independent of frame rate
const TimeStep float32 = 1.0/30

type TickInfo struct {
	Dt float32
	Planet *world.Planet
}

var prevTime,nextTime = float32(0),float32(0)
var now = float32(0)

// will a physics tick run after this time delta?
// use to decide whether we consume keyboard inputs etc.
func WillTick(dt float32) bool {
	return now + dt > nextTime
}

// run the necessary number of physics ticks for a given time delta
func Simulate(dt float32, planet *world.Planet) {
	now += dt
	// run physics ticks until the current time lies between
	// the previous physics state and the next physics state.
	// Then, we can interpolate
	for nextTime<now {
		tick(planet)
	}

	// physics simulation is done; save interpolated values for rendering
	alpha := (now-prevTime) / TimeStep
	for idx,_ := range bodies {
		bodies[idx].UpdateInterpolatedTransform(alpha)
	}
}

func tick(planet *world.Planet) {
	prevTime = nextTime
	nextTime += TimeStep

	for idx,_ := range bodies {
		bodies[idx].SavePreviousTransform()
		bodies[idx].Move(planet, TimeStep)
	}
}

