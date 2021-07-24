package particle_draw

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/render"
)

var particleShader *render.Shader

func FrustumCull(particleList []*particles.Particle, cam *camera.Camera) []*particles.Particle {
	particlesToDraw := make([]*particles.Particle, 0, len(particleList))

	for _, par := range particleList {
		//assume particle radius is 1
		if cam.IsInBoundsRadius(par.Position, 1) {
			particlesToDraw = append(particlesToDraw, par)
		}
	}

	return particlesToDraw
}

func BinByDrawHandlerID(particleList []*particles.Particle) map[types.ParticleDrawHandlerId][]*particles.Particle {
	bins := make(map[types.ParticleDrawHandlerId][]*particles.Particle)
	for _, par := range particleList {
		bins[par.DrawHandlerID] = append(bins[par.DrawHandlerID], par)
	}
	return bins
}
