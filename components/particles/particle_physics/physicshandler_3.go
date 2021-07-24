package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
)

// "drifts" at fixed velocity, no gravity

func PhysicsHandlerDrift(particleList []*particles.Particle) {
	// fmt.Println(len(particleList))
	for _, particle := range particleList {
		particle.Position = particle.Position.Add(particle.Velocity.Mult(constants.TimeStep))
	}
}
