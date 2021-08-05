package particle_draw

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/spriteloader"
)

func DrawParticles(particleList *particles.ParticleList, cam *camera.Camera) {
	particlesToDraw := FrustumCull(particleList.Particles, cam)

	bins := BinByDrawHandlerID(particlesToDraw)

	for drawHandlerID, agents := range bins {
		GetDrawHandler(drawHandlerID)(agents, cam)
	}

}
func Init() {
	RegisterDrawHandler(constants.PARTICLE_DRAW_HANDLER_NULL, DrawNull)
	RegisterDrawHandler(
		constants.PARTICLE_DRAW_HANDLER_SOLID, DrawSolid)
	RegisterDrawHandler(
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT, DrawTransparent)
	RegisterDrawHandler(constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED, DrawTransparentInstanced)

	AssertAllDrawHandlersRegistered()

	spriteloader.LoadSingleSprite("./assets/particles/particle.png", "particle")
	spriteloader.LoadSingleSprite("./assets/particles/star.png", "star")
	quad_vao = makeQuadVao()
	initDrawInstanced()
}
