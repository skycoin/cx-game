package particles

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
)

type Particle struct {
	Position         cxmath.Vec2
	Velocity         cxmath.Vec2
	TimeToLive       float32
	Duration         float32
	Texture          uint32
	DrawHandlerID    types.ParticleDrawHandlerId
	PhysicsHandlerID types.ParticlePhysicsHandlerID
}

func (p *Particle) Update(dt float32) {

}
