package particle_emitter

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spriteloader"
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

const (
	bulletSize float32 = 0.3
	bulletElasticity float32 = 0.5
	bulletFriction float32 = 0.1
	bulletDuration float32 = 3

	bulletDamage int = 1
)

func (emitter *BulletEmitter) Emit(position, velocity cxmath.Vec2) {
	body := particles.NewParticleBody(
		position, velocity,
		bulletSize,bulletElasticity,bulletFriction,
	)
	texture := spriteloader.GetSpriteIdByNameUint32("star")
	particle := particles.NewParticle(
		body, texture, bulletDuration,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
		constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT_CALLBACK,
		emitter.OnHitCallback(),
	)
	particle.Damage = bulletDamage

	emitter.particleList.AddParticle(particle)
}

func (emitter *BulletEmitter) OnHitCallback() func(*particles.Particle) {
	return func(particle *particles.Particle) {
		sparkEmitter.Emit(particle)
	}
}
