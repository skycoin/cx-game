package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

// Does not collide, apply gravity

func PhysicsHandlerOxygen(
		Particles []*particles.Particle, World *world.World,
) {
	for _, par := range Particles {
		par.MoveSlowXAxis( &World.Planet,
			constants.MS_PER_TICK, cxmath.Vec2{}, false, par.SlowdownFactor)
	}
}
