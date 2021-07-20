package particle_physics

import (
	"log"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type ParticlePhysicsHandler func([]*particles.Particle)

var ParticlePhysicsHandlerList [constants.NUM_PARTICLE_PHYSICS_HANDLERS]ParticlePhysicsHandler

// 1> bounces, gravity

// 2> Does not bounce, stops on impact, gravity

// 3> "drifts" at fixed velocity, no gravity

func Init() {
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_NULL,
		PhysicsHandlerNull )
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_BOUNCE_GRAVITY,
		PhysicsHandlerBounceGravity )
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_GRAVITY,
		PhysicsHandlerGravity )
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_DRIFT,
		PhysicsHandlerDrift )

	AssertAllParticleHandlersRegistered()
}

func Update(dt float32) {

}

func RegisterPhysicsHandler(
		id types.ParticlePhysicsHandlerID,
		handler ParticlePhysicsHandler,
) {
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
