package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

func PhysicsHandlerDissappearOnHitCallback(particleList []*particles.Particle, planet *world.Planet) {
	for _, par := range particleList {
		par.MoveNoBounceRaytrace(planet, constants.PHYSICS_TICK, cxmath.Vec2{})
		if par.Collisions.Collided() {
			par.OnCollideCallback(par)
			par.Die()
		}
	}
}
