package particle_emitter

import (
	"fmt"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

type SparkEmitter struct {
	particleList *particles.ParticleList
	program      *render.Program
}

func NewSparkEmitter(particleList *particles.ParticleList) *SparkEmitter {
	return &SparkEmitter{
		particleList: particleList,
	}
}
func (emitter *SparkEmitter) Emit(particle *particles.Particle) {

	// direction := cxmath.Vec2{0, 1}
	fmt.Println("EMITTED")

	for i := 0; i < 10; i++ {
		emitter.particleList.AddParticle(
			particle.Pos,
			cxmath.Vec2{0, 10},
			1,
			1,
			0,
			spriteloader.GetSpriteIdByNameUint32("star"),
			3,
			constants.PARTICLE_DRAW_HANDLER_TRANSPARENT,
			constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT_CALLBACK,
			nil,
		)
	}
}
