package particle_constants

import "github.com/skycoin/cx-game/components/types"

const (
	MAX_PARTICLE_COUNT = 100000
)

const (
	//draw_handlers constants
	DRAW_HANDLER_NULL types.ParticleDrawHandlerId = iota
	DRAW_HANDLER_SOLID
	DRAW_HANDLER_ALPHA_BLENDED
)

const (
	//physics_handlers constants
	PHYSICS_HANDLER_NULL types.ParticlePhysicsHandlerID = iota
	PHYSICS_HANDLER_GRAVITY_BOUNCE
	PHYSICS_HANDLER_GRAVITY_NO_BOUNCE
	PHYSICS_HANDLER_NO_GRAVITY
)
