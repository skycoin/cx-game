package agent_physics

import (
	"github.com/skycoin/cx-game/constants/physicsconstants"
	"github.com/skycoin/cx-game/world"
)

func UpdateAgents(worldState *world.WorldState, planet *world.Planet) {
	for _, agent := range worldState.AgentList.Agents {
		agent.PhysicsState.Move(planet, physicsconstants.PHYSICS_TIMESTEP)
	}
}
