package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

// "drifts" at fixed velocity, no gravity

func PhysicsHandlerDrift(particleList []*particles.Particle, planet *world.Planet) {
	// fmt.Println(len(particleList))
	for _, particle := range particleList {
		// particle.Position = particle.Position.Add(particle.Verlet.Mult(constants.TimeStep))
		particle.Integrate(constants.TimeStep)
	}
}
