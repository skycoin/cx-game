package agent_physics

import (
	"github.com/skycoin/cx-game/constants/physicsconstants"
	"github.com/skycoin/cx-game/world"
)

func UpdateAgents(World *world.World) {
	for _, agent := range World.Entities.Agents.Get() {
		agent.PhysicsState.
			Move(&World.Planet, physicsconstants.PHYSICS_TIMESTEP)
	}
}
