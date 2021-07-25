package particles

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
)

//for now keep one global particles list, redo later
type ParticleList struct {
	Particles []*Particle
}

func (pl *ParticleList) AddParticle(
	position cxmath.Vec2,
	velocity cxmath.Vec2,
	texture uint32,
	duration float32,
	drawHandlerId types.ParticleDrawHandlerId,
	physiscHandlerID types.ParticlePhysicsHandlerID,
) {
	newParticle := Particle{
		Position:         position,
		Velocity:         velocity,
		Duration:         duration,
		TimeToLive:       duration,
		Texture:          spriteloader.GetSpriteIdByNameUint32("particle"),
		DrawHandlerID:    drawHandlerId,
		PhysicsHandlerID: physiscHandlerID,
	}

	pl.Particles = append(pl.Particles, &newParticle)
}

func (pl *ParticleList) Update(dt float32) {
	for _, par := range pl.Particles {
		par.Update(dt)

	}
}

func (pl *ParticleList) deleteParticles(indexes []int) {

}
