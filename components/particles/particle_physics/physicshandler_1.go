package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

//bounces, gravity
func PhysicsHandlerBounceGravity(particleList []*particles.Particle, planet *world.Planet) {
	for _, par := range particleList {
		par.MoveBounceGravity(planet, constants.TimeStep)
	}
}
