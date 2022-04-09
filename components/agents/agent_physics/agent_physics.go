package agent_physics

import (
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

func UpdateAgents(World *world.World) {
	for _, agent := range World.Entities.Agents.GetAllAgents() {
		//check for nil first
		if agent == nil {
			continue
		}
		agent.Transform.
			Move(&World.Planet, constants.MS_PER_TICK)
	}
}
