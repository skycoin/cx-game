package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

func PhysicsHandlerDissappearOnHitCallback(particleList []*particles.Particle, planet *world.Planet) {
	for _, par := range particleList {
		par.MoveNoBounceGravity(planet, constants.TimeStep)
		if par.Collisions.Collided() {
			par.Callback(par)
			par.Die()
		}
	}
}
