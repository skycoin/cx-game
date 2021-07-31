package particles

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
)

//for now keep one global particles list, redo later
type ParticleList struct {
	Particles []*Particle
}

func (pl *ParticleList) AddParticle(
	position cxmath.Vec2,
	velocity cxmath.Vec2,
	size float32,
	elasticity float32,
	friction float32,
	texture uint32,
	duration float32,
	drawHandlerId types.ParticleDrawHandlerId,
	physiscHandlerID types.ParticlePhysicsHandlerID,
) {
	newParticle := Particle{
		ParticleBody: ParticleBody{
			Pos:        position,
			Vel:        velocity,
			Size:       cxmath.Vec2{size, size},
			Elasticity: elasticity,
			Friction:   friction,
		},
		Duration:         duration,
		TimeToLive:       duration,
		Texture:          texture,
		DrawHandlerID:    drawHandlerId,
		PhysicsHandlerID: physiscHandlerID,
	}

	pl.Particles = append(pl.Particles, &newParticle)
}

func (pl *ParticleList) Update(dt float32) {
	particlesToDelete := make([]int, 0)
	for i, par := range pl.Particles {

		// -1 is infinite
		if par.Duration != -1 {
			par.TimeToLive -= dt
			if par.TimeToLive <= 0 {
				particlesToDelete = append(particlesToDelete, i)
			}
		}
	}
	pl.deleteParticles(particlesToDelete)
}

func (pl *ParticleList) deleteParticles(indexes []int) {
	var newParticleList []*Particle
	for i, par := range pl.Particles {
		toBeDeleted := false
		for _, j := range indexes {
			if i == j {
				toBeDeleted = true
				break
			}
		}
		if !toBeDeleted {
			newParticleList = append(newParticleList, par)
		}
	}
	pl.Particles = newParticleList
}
