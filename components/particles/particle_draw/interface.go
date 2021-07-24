package particle_draw

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
)

func DrawParticles(particleList *particles.ParticleList, cam *camera.Camera) {
	particlesToDraw := FrustumCull(particleList.Particles, cam)

	bins := BinByDrawHandlerID(particlesToDraw)

	for drawHandlerID, agents := range bins {
		GetDrawHandler(drawHandlerID)(agents, cam)
	}

}
func Init() {
	RegisterDrawHandler(
		constants.PARTICLE_DRAW_HANDLER_SOLID, DrawSolid)
	RegisterDrawHandler(
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT, DrawTransparent)

	AssertAllDrawHandlersRegistered()

}
