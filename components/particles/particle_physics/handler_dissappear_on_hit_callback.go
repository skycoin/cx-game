package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

func PhysicsHandlerDissappearOnHitCallback(
		Particles []*particles.Particle, World *world.World,
) {
	for _, par := range Particles {
		par.MoveNoBounceRaytrace(
			&World.Planet, World.Entities.Agents.GetAllAgents(),
			constants.PHYSICS_TICK, cxmath.Vec2{} )
		if par.Collisions.Collided() {
			par.OnCollideCallback(par)
			par.Die()
		}
	}
}
