package particles

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics"
)

type Particle struct {
	Verlet
	physics.Body
	TimeToLive       float32
	Duration         float32
	Texture          uint32
	DrawHandlerID    types.ParticleDrawHandlerId
	PhysicsHandlerID types.ParticlePhysicsHandlerID
}

type Verlet struct {
	Position    cxmath.Vec2
	OldPosition cxmath.Vec2
}

func NewVerlet(position, velocity cxmath.Vec2) Verlet {
	verlet := Verlet{
		Position:    position,
		OldPosition: position.Sub(velocity),
	}
	return verlet
}

func (v *Verlet) Integrate(dt float32) {
	v.Position = v.Position.Add(v.Position.Sub(v.OldPosition).Mult(dt))
}
