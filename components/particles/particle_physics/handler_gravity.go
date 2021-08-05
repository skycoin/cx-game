package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

//  Does not bounce, stops on impact, gravity

func PhysicsHandlerGravity(particleList []*particles.Particle, planet *world.Planet) {
	// for _, par := range particleList {
	// }
	for _, par := range particleList {
		par.MoveNoBounce(planet, constants.TimeStep, cxmath.Vec2{0, -constants.Gravity})
	}
}
