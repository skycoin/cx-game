package particles

import (
	"github.com/skycoin/cx-game/components/types"
)

type Particle struct {
	ParticleId types.ParticleID
	ParticleBody
	TimeToLive        float32
	Duration          float32
	Texture           uint32
	DrawHandlerID     types.ParticleDrawHandlerId
	PhysicsHandlerID  types.ParticlePhysicsHandlerID
	OnCollideCallback func(*Particle)
	//embed the struct for easier access
	ParticleMeta
}

type ParticleMeta struct {
	//should be between 0 and 1
	SlowdownFactor float32
}

func (p *Particle) Die() {
	p.TimeToLive = 0
}
