package particle_physics

import (
	"log"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

type ParticlePhysicsHandler func([]*particles.Particle, *world.Planet)

var ParticlePhysicsHandlerList [constants.NUM_PARTICLE_PHYSICS_HANDLERS]ParticlePhysicsHandler

// 1> bounces, gravity

// 2> Does not bounce, stops on impact, gravity

// 3> "drifts" at fixed velocity, no gravity

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
