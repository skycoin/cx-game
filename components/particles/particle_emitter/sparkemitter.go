package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

var sparkEmitter *SparkEmitter

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
	for i := 0; i < 10; i++ {
		velocity := cxmath.Vec2{
			X: (rand.Float32() - 0.5) * 3,
			Y: 30 * rand.Float32(),
		}
		emitter.particleList.AddParticle(
			particle.Pos,
			velocity,
			// (rand.Float32()+0.3)/4,
			1,
			0,
			0,
			spriteloader.GetSpriteIdByNameUint32("star"),
			rand.Float32()*(emitter.maxduration-emitter.minduration)+emitter.minduration,
			constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
			constants.PARTICLE_PHYSICS_HANDLER_NO_COLLISION_GRAVITY,
			nil,
		)
	}
}
