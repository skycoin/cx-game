package particle_physics

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/engine/ui"
)

func Init() {
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_NULL,
		PhysicsHandlerNull)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_BOUNCE_GRAVITY,
		PhysicsHandlerBounceGravity)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_COLLISION_GRAVITY,
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
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_NO_COLLISION_GRAVITY,
		PhysicsHandlerNoCollisionGravity,
	)
	RegisterPhysicsHandler(
		constants.PARTICLE_PHYSICS_HANDLER_OXYGEN,
		PhysicsHandlerOxygen,
	)

	AssertAllParticleHandlersRegistered()
}

func Update(World *world.World) {
	particleList := &World.Entities.Particles
	bins := BinByPhysicsHandlerID(particleList.Particles)

	for physicsType, par := range bins {
		GetParticlePhysicsHandler(physicsType)(par, World)
	}

	// behaviour independent of particle type
	for _,particle := range particleList.Particles {
		isDamaging :=
			particle != nil && particle.IsHittingAgent && particle.Damage > 0
		if isDamaging {
			agent := World.Entities.Agents.
				FromID(particle.ParticleBody.HitAgentID)
			agent.TakeDamage(particle.Damage)

			ApplyKnockback(&agent.PhysicsState, particle.Vel.Mgl32())

			agentPos := agent.PhysicsState.Pos
			ui.CreateDamageIndicator(
				particle.Damage, agentPos.Mgl32().Add(mgl32.Vec2{0,1}) )

			particle.Damage = 0 // only hit agent once
		}
	}
}
