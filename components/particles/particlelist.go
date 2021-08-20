package particles

import (
	"github.com/skycoin/cx-game/components/types"
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
		return -1
	}
	returnValue := q.queue[0]
	q.queue = q.queue[1:]
	return returnValue
}

func (pl *ParticleList) AddParticle(particle Particle) types.ParticleID {

	//id should match particle id
	newId := pl.idQueue.Pop()

	if newId == -1 {
		pl.Particles = append(pl.Particles, &particle)
		particle.ParticleId = types.ParticleID(len(pl.Particles) - 1)
	} else {
		particle.ParticleId = types.ParticleID(newId)
		pl.Particles[newId] = &particle
	}
	return particle.ParticleId
}

func (pl *ParticleList) Update(dt float32) {
	// fmt.Println(pl.idQueue)
	particlesToDelete := make([]int, 0)
	for i, par := range pl.Particles {
		if par != nil {
			// -1 is infinite
			if par.Duration != -1 {
				par.TimeToLive -= dt
				if par.TimeToLive <= 0 {
					particlesToDelete = append(particlesToDelete, i)
				}
			}
		}

	}
	pl.deleteParticles(particlesToDelete)
}

func (pl *ParticleList) deleteParticles(indexes []int) {
	for _, i := range indexes {
		pl.Particles[i] = nil
		pl.idQueue.Push(i)
	}
	// var newParticleList []*Particle
	// for i, par := range pl.Particles {
	// 	toBeDeleted := false
	// 	for _, j := range indexes {
	// 		if i == j {
	// 			toBeDeleted = true
	// 			pl.idQueue.Push(int(par.ParticleId))
	// 			break
	// 		}
	// 	}
	// 	if !toBeDeleted {
	// 		newParticleList = append(newParticleList, par)
	// 	}
	// }
	// pl.Particles = newParticleList
}

func (pl *ParticleList) GetParticle(id types.ParticleID) *Particle {

	return pl.Particles[int(id)]

}
