package particle_physics

import (
	"log"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type ParticlePhysicsHandler func()

var ParticlePhysicsHandlerList [constants.NUM_PARTICLE_PHYSICS_HANDLERS]ParticlePhysicsHandler

// 1> bounces, gravity

// 2> Does not bounce, stops on impact, gravity

// 3> "drifts" at fixed velocity, no gravity

func Init() {
	RegisterDrawHandler(constants.PARTICLE_PHYSICS_HANDLER_1, PhysicsHandler1)
	RegisterDrawHandler(constants.PARTICLE_PHYSICS_HANDLER_2, PhysicsHandler2)
	RegisterDrawHandler(constants.PARTICLE_PHYSICS_HANDLER_3, PhysicsHandler3)

	AssertAllParticleHandlersRegistered()
}

func Update(dt float32) {
	
}

func RegisterDrawHandler(id types.ParticlePhysicsHandlerID, handler ParticlePhysicsHandler) {
	ParticlePhysicsHandlerList[id] = handler
}

func AssertAllParticleHandlersRegistered() {
	for _, handler := range ParticlePhysicsHandlerList {
		if handler == nil {
			log.Fatalln("Did not initialize all particle physics handlers")
		}
	}
}

func GetParticlePhysicsHandler(id types.ParticlePhysicsHandlerID) ParticlePhysicsHandler {
	return ParticlePhysicsHandlerList[id]
}
