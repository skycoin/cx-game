package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

var dustEmitter *DustEmitter

type DustEmitter struct {
	particleList *particles.ParticleList
	program      *render.Program
	minduration  float32
	maxduration  float32
}

func NewDustEmitter(particleList *particles.ParticleList) *DustEmitter {
	return &DustEmitter{
		particleList: particleList,
		minduration:  0.3,
		maxduration:  0.5,
	}
}

const (
	dustElasticity float32 = 0
	dustDuration   float32 = 0
)

func (emitter *DustEmitter) Emit(parent cxmath.Vec2) {
	for i := 0; i < 10; i++ {
		velocity := cxmath.Vec2{
			X: (rand.Float32() - 0.5) * 10,
			Y: 20 * rand.Float32(),
		}
		size := (rand.Float32() + 0.3) / 4
		duration :=
			rand.Float32()*(emitter.maxduration-emitter.minduration) +
				emitter.minduration
		particle := particles.NewParticle(
			particles.NewParticleBody(
				parent, velocity, size, dustElasticity, dustDuration,
			),
			spriteloader.GetSpriteIdByNameUint32("dust"),
			duration,
			constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
			constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT,
			nil,
		)
		emitter.particleList.AddParticle(particle)
	}
}
