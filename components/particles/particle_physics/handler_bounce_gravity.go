package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

//bounces, gravity
func PhysicsHandlerBounceGravity(
		Particles []*particles.Particle, World *world.World,
) {
	for _, par := range Particles {
		par.MoveBounce( &World.Planet,
			constants.PHYSICS_TICK, cxmath.Vec2{0, -constants.Gravity} )
	}
}
