package particle_emitter

import (
	"github.com/skycoin/cx-game/components/particles"
)

func Init(particleList *particles.ParticleList) {
	bulletEmitter = NewBulletEmitter(
		particleList,
	)
	sparkEmitter = NewSparkEmitter(particleList)
}
