package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
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

func (emitter *ParticleEmitter) SetPosition(position cxmath.Vec2) {
	emitter.position = position
}

func (emitter *ParticleEmitter) Emit() {
	emitter.particleList.AddParticle(
		emitter.position,
		//set velocity up
		cxmath.Vec2{(rand.Float32()*5 - 2.5), rand.Float32() + 0.5},
		0.3,
		0.5,
		0.1,
		spriteloader.GetSpriteIdByNameUint32("star"),
		5,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
		constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT,
		nil,
	)
}
