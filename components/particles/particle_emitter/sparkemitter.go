package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

type SparkEmitter struct {
	particleList *particles.ParticleList
	program      *render.Program
	minduration  float32
	maxduration  float32
}

func NewSparkEmitter(particleList *particles.ParticleList) *SparkEmitter {
	return &SparkEmitter{
		particleList: particleList,
		minduration:  1.3,
		maxduration:  1.5,
	}
}
func (emitter *SparkEmitter) Emit(particle *particles.Particle) {

	// direction := cxmath.Vec2{0, 1}

	for i := 0; i < 10; i++ {
		velocity := particle.PrevVel.Mult(-1).Normalize().Add(cxmath.Vec2{0, 0}).Mult(35)
		emitter.particleList.AddParticle(
			particle.Pos,
			velocity,
			1,
			0,
			0,
			spriteloader.GetSpriteIdByNameUint32("star"),
			rand.Float32()*(emitter.maxduration-emitter.minduration)+emitter.minduration,
			constants.PARTICLE_DRAW_HANDLER_TRANSPARENT,
			constants.PARTICLE_PHYSICS_HANDLER_BOUNCE_GRAVITY,
			nil,
		)
	}
}
