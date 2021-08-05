package particle_emitter

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
)

var bulletEmitter *BulletEmitter

type BulletEmitter struct {
	particleList *particles.ParticleList
}

func NewBulletEmitter(particlelist *particles.ParticleList) *BulletEmitter {
	bulletEmitter := BulletEmitter{
		particleList: particlelist,
	}

	return &bulletEmitter
}

func (emitter *BulletEmitter) Emit(position, velocity cxmath.Vec2) {
	emitter.particleList.AddParticle(
		position,
		//set velocity up
		velocity,
		0.3,
		0.5,
		0.1,
		spriteloader.GetSpriteIdByNameUint32("star"),
		3,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
		constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT_CALLBACK,
		emitter.OnHitCallback(),
	)
}

func (emitter *BulletEmitter) OnHitCallback() func(*particles.Particle) {
	return func(particle *particles.Particle) {
		sparkEmitter.Emit(particle)
	}
}
