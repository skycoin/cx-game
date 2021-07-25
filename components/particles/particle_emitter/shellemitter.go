package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
)

type ParticleEmitter struct {
	position     cxmath.Vec2
	particleList *particles.ParticleList
	// gundirection cxmath.Vec2
}

func NewParticle(position cxmath.Vec2, particlelist *particles.ParticleList) *ParticleEmitter {
	particleEmitter := ParticleEmitter{
		position:     position,
		particleList: particlelist,
	}

	return &particleEmitter
}

func (emitter *ParticleEmitter) SetData(position cxmath.Vec2) {
	emitter.position = position

}
func (emitter *ParticleEmitter) Emit() {
	emitter.particleList.AddParticle(
		emitter.position,
		cxmath.Vec2{rand.Float32(), rand.Float32()},
		0,
		1,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT,
		constants.PARTICLE_PHYSICS_HANDLER_DRIFT,
	)
}
