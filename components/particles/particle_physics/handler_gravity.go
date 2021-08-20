package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

//  Does not bounce, stops on impact, gravity

func PhysicsHandlerGravity(Particles []*particles.Particle, World *world.World) {
	// for _, par := range Particles {
	// }
	for _, par := range Particles {
		par.MoveNoBounce(&World.Planet, constants.PHYSICS_TICK, cxmath.Vec2{0, -constants.Gravity})
	}
}
