package particle_physics

import "github.com/skycoin/cx-game/components/particles"

// "drifts" at fixed velocity, no gravity

func PhysicsHandler3(particleList []*particles.Particle) {
	for _, particle := range particleList {
		particle.Position = particle.Position.Add(particle.Velocity)
	}
}
