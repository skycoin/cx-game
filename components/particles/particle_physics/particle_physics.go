package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
)

func BinByPhysicsHandlerID(particleList []*particles.Particle) map[types.ParticlePhysicsHandlerID][]*particles.Particle {
	bins := make(map[types.ParticlePhysicsHandlerID][]*particles.Particle)

	for _, par := range particleList {
		bins[par.PhysicsHandlerID] = append(bins[par.PhysicsHandlerID], par)
	}
	return bins
}
