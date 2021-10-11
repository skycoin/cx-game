package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/world"
)

/*
Give emitter

TTL, time to live, decremented each time tick
init function
precondition, check every tick
postcondition or teardown, what happens when ttl is expired
different parameters like "tracked agentId", and other misc parameters emitters can have, but not all parameters will be used
For example, oxygen bubble emitter will track player,
-- and will check if oxygen generator is within range

and if oxygen generator is not in range, will deconstruct itself or set TTL to zero
and if there is oxygen generator, set TTL to 30, meaning will expire automatically in 30 ticks
and in update function, emitter can change its position, track object etc
*/

//object pool
var oxygenEmitters []OxygenEmitter

var emitterTTL = 30

type OxygenEmitter struct {
	//ttl in ticks

	//-1 is unset, ready to by set, 0 means it expired
	TTL          int
	Position     cxmath.Vec2
	agentId      types.AgentID
	particleList *particles.ParticleList
}

// FIXME Why is the NewOxygenEmitter() method responsible 
// for both registering and creating the emitter?
// Additionally, why do we return a copy rather than a pointer/ID ?
func NewOxygenEmitter(trackedId types.AgentID, particleList *particles.ParticleList) OxygenEmitter {

	for i := range oxygenEmitters {
		if oxygenEmitters[i].IsFree() {
			oxygenEmitters[i].agentId = trackedId
			return oxygenEmitters[i]
		}
	}

	newOxygenEmitter := OxygenEmitter{
		agentId:      trackedId,
		particleList: particleList,
	}

	oxygenEmitters = append(oxygenEmitters, newOxygenEmitter)

	return newOxygenEmitter
}

//works when entering/reentering near oxygen generator
func (e *OxygenEmitter) Init() {
	e.TTL = emitterTTL
	//emit
}

//assume dt is not fixed
func (emitter *OxygenEmitter) Update(dt float32, currentWorld *world.World) {
	//change position
	agent := currentWorld.Entities.Agents.FromID(emitter.agentId)
	emitter.Position = agent.Transform.Pos

	if currentWorld.Planet.NearOxygenGenerator(emitter.Position) {
		if emitter.TTL > 0 {
			//emit logic
			//every fifth tick
			if emitter.TTL%5 == 0 {
				emitter.Emit()
			}
			emitter.TTL -= 1
		} else if emitter.TTL == -1 {
			emitter.Init()
			return
		} else {
			emitter.Teardown()
		}
	} else {
		//unset emitter
		emitter.Teardown()
		emitter.Reset()
	}
}

// works when exited oxygen generator
func (emitter *OxygenEmitter) Teardown() {
	//todo add bubble pop
	emitter.TTL = 0

}

func (emitter *OxygenEmitter) Reset() {
	emitter.TTL = -1
}

const (
	oxygenSize float32 = 1
	oxygenElasticity float32 = 0
	oxygenFriction float32 = 0
	oxygenDuration float32 = 3
)

func (emitter *OxygenEmitter) Emit() {
	emitter.TTL -= 1

	body := particles.NewParticleBody(
		emitter.Position,
		cxmath.Vec2{emitter.getBubbleX(), emitter.getBubbleY()},
		oxygenSize, oxygenFriction, oxygenDuration,
	)
	texture := spriteloader.GetSpriteIdByNameUint32("star")
	particle := particles.NewParticle(
		body, texture, oxygenDuration,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT,
		constants.PARTICLE_PHYSICS_HANDLER_OXYGEN,
		nil,
	)
	id := emitter.particleList.AddParticle(particle)
	particleInList := emitter.particleList.GetParticle(id)
	particleInList.SlowdownFactor = 0.9
}

func (emitter *OxygenEmitter) Detach() {
	emitter.agentId = -1
}

func (emitter *OxygenEmitter) IsFree() bool {
	return emitter.agentId == -1
}

func (emitter *OxygenEmitter) getBubbleX() float32 {
	direction := getDirection()

	return rand.Float32() * 5 * direction
}

func (emitter *OxygenEmitter) getBubbleY() float32 {
	return rand.Float32() + 0.5*2
}

func getDirection() float32 {
	return math32.Sign(rand.Float32() - 0.5)
}
