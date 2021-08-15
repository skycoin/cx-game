package particle_emitter

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
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

type OxygenEmitter struct {
	//ttl in ticks

	//-1 is unset, ready to by set, 0 means it expired
	TTL          float32
	Position     cxmath.Vec2
	agentId      types.AgentID
	particleList *particles.ParticleList
}

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
	e.TTL = 30
	//emit
}

//assume dt is not fixed
func (emitter *OxygenEmitter) Update(dt float32, currentWorld *world.World) {
	//change position
	agent := currentWorld.Entities.Agents.FromID(emitter.agentId)
	emitter.Position = agent.PhysicsState.Pos

	if currentWorld.Planet.NearOxygenGenerator(emitter.Position) {
		if emitter.TTL > 0 {
			//emit logic
			emitter.Emit()
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
	emitter.TTL = 0

}

func (emitter *OxygenEmitter) Reset() {
	emitter.TTL = -1
}

func (emitter *OxygenEmitter) Emit() {
	emitter.TTL -= 1

	id := emitter.particleList.AddParticle(
		emitter.Position,
		cxmath.Vec2{rand.Float32(), 35},
		1,
		0,
		0,
		spriteloader.GetSpriteIdByNameUint32("star"),
		5,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
		constants.PARTICLE_PHYSICS_HANDLER_OXYGEN,
		nil,
	)

	particle := emitter.particleList.GetParticle(id)
	particle.SlowdownFactor = 1
}

func (emitter *OxygenEmitter) Detach() {
	emitter.agentId = -1
}

func (emitter *OxygenEmitter) IsFree() bool {
	return emitter.agentId == -1
}
