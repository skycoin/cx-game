package constants

import "github.com/skycoin/cx-game/components/types"

const (
	AI_HANDLER_NULL types.AgentAiHandlerID = iota

	NUM_AI_HANDLERS
)

//agent drawhandler constants
const (
	DRAW_HANDLER_NULL types.AgentDrawHandlerID = iota
	DRAW_HANDLER_QUAD

	NUM_AGENT_DRAW_HANDLERS = 2
)

//agent physics constants
const (
	PHYSICS_HANDLER_NULL types.AgentPhysicsHandlerID = iota
	PHYSICS_HANDLER_1
	PHYSICS_HANDLER_2

	NUM_AGENT_PHYSICS_HANDLERS = 3
)

//particle drawhandler constants
const (
	PARTICLE_DRAW_HANDLER_1 types.ParticleDrawHandlerId = iota
	PARTICLE_DRAW_HANDLER_2

	NUM_PARTICLE_DRAW_HANDLERS = 2
)

//particle physicshandler constants
const (
	PARTICLE_PHYSICS_HANDLER_1 types.ParticlePhysicsHandlerID = iota
	PARTICLE_PHYSICS_HANDLER_2
	PARTICLE_PHYSICS_HANDLER_3

	NUM_PARTICLE_PHYSICS_HANDLERS = 3
)
