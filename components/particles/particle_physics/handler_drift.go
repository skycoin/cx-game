package particle_physics

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
)

// "drifts" at fixed velocity, no gravity

func PhysicsHandlerDrift(Particles []*particles.Particle, World *world.World) {
	for _, par := range Particles {
		par.MoveNoCollision( &World.Planet,
			constants.PHYSICS_TICK, cxmath.Vec2{} )
	}
}
