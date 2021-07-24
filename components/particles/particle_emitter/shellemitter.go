package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants/particle_constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
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
		spriteloader.GetSpriteIdByName("basic-agent"),
		1,
		particle_constants.DRAW_HANDLER_SOLID,
		particle_constants.PHYSICS_HANDLER_NO_GRAVITY,
	)
}
