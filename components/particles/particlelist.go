package particles

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
)

//for now keep one global particles list, redo later
type ParticleList struct {
	Particles []*Particle
	idQueue   QueueI
}

type QueueI struct {
	queue []int
}

func NewQueue() QueueI {
	queue := QueueI{
		queue: make([]int, 0),
	}
	return queue
}

func (q *QueueI) Push(n int) {
	q.queue = append(q.queue, n)
}
func (q *QueueI) Pop() int {
	if len(q.queue) == 0 {
		particleIdCounter += 1
		return particleIdCounter
	}
	returnValue := q.queue[0]
	q.queue = q.queue[1:]
	return returnValue
}

var particleIdCounter int = -1

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
	callback func(*Particle),
) types.ParticleID {
	newParticle := Particle{
		ParticleId: types.ParticleID(pl.idQueue.Pop()),
		ParticleBody: ParticleBody{
			Pos:        position,
			Vel:        velocity,
			Size:       cxmath.Vec2{size, size},
			Elasticity: elasticity,
			Friction:   friction,
		},
		Duration:          duration,
		TimeToLive:        duration,
		Texture:           texture,
		DrawHandlerID:     drawHandlerId,
		PhysicsHandlerID:  physiscHandlerID,
		OnCollideCallback: callback,
	}

	pl.Particles = append(pl.Particles, &newParticle)
	return newParticle.ParticleId
}

func (pl *ParticleList) Update(dt float32) {
	// fmt.Println(pl.idQueue)
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
				pl.idQueue.Push(i)
				break
			}
		}
		if !toBeDeleted {
			newParticleList = append(newParticleList, par)
		}
	}
	pl.Particles = newParticleList
}

func (pl *ParticleList) GetParticle(id types.ParticleID) *Particle {

	return pl.Particles[int(id)]

}
