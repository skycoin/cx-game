package particle_draw

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/engine/camera"
)

const particleRadius int = 1

// filter only particles visible by camera
func FrustumCull(
	particleList []*particles.Particle, cam *camera.Camera,
) []*particles.Particle {
	particlesToDraw := make([]*particles.Particle, 0)

	for _, par := range particleList {
		if par == nil {
			continue
		}
		//assume particle radius is 1
		if cam.IsInBoundsRadius(par.Pos, particleRadius) {
			particlesToDraw = append(particlesToDraw, par)
		}
	}

	return particlesToDraw
}

func BinByDrawHandlerID(
	particleList []*particles.Particle,
) map[types.ParticleDrawHandlerId][]*particles.Particle {
	bins := make(map[types.ParticleDrawHandlerId][]*particles.Particle)
	for _, par := range particleList {
		bins[par.DrawHandlerID] = append(bins[par.DrawHandlerID], par)
	}
	return bins
}
