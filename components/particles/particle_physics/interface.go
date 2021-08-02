package particle_physics

import (
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

func Init() {
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_NULL,
		PhysicsHandlerNull)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_BOUNCE_GRAVITY,
		PhysicsHandlerBounceGravity)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_GRAVITY,
		PhysicsHandlerGravity)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_DRIFT,
		PhysicsHandlerDrift)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT,
		PhysicsHandlerDissappearOnHit)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT_CALLBACK,
		PhysicsHandlerDissappearOnHitCallback)

	AssertAllParticleHandlersRegistered()
}

func Update(World *world.World) {
	particleList := &World.Entities.Particles
	planet := &World.Planet
	bins := BinByPhysicsHandlerID(particleList.Particles)

	for physicsType, par := range bins {
		GetParticlePhysicsHandler(physicsType)(par, planet)
	}
}
