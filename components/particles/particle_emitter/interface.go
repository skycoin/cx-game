package particle_emitter

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
)

func CreateBullet(position, velocity mgl32.Vec2) {
	bulletEmitter.Emit(
		cxmath.Vec2{position.X(), position.Y()},
		cxmath.Vec2{velocity.X(), velocity.Y()},
	)
}

func EmitDust(position cxmath.Vec2) {
	dustEmitter.Emit(position)
}

func EmitOxygen(agentId types.AgentID, particleList *particles.ParticleList) {
	// create oxygen emitter.
	// it will be added to object pool and updated each tick,
	// checking for conditions
	NewOxygenEmitter(agentId, particleList)
}
