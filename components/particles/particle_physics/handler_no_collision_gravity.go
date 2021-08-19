package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

// Does not collide, apply gravity

func PhysicsHandlerNoCollisionGravity(particleList []*particles.Particle, planet *world.Planet) {
	// for _, par := range particleList {
	// }
	for _, par := range particleList {
		par.MoveNoCollision(planet, constants.PHYSICS_TICK, cxmath.Vec2{
			0,
			-constants.Gravity * 5,
		})
	}
}
