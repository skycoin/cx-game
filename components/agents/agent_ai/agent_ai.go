package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/world"
)

// TODO since AI will be run server-side,
// we should not be passing any particular player to this function.
// Rather, the list of players should be computed
// by filtering the world agents
func UpdateAgents(World *world.World, player *agents.Agent) {
	ctx := AiContext{
		World:      World,
		WorldWidth: float32(World.Planet.Width),
		PlayerPos:  player.Transform.Pos.Mgl32(),
	}
	for _, agent := range World.Entities.Agents.GetAllAgents() {
		if agent == nil {
			continue
		}
		aiHandlers[agent.Handlers.AI](agent, ctx)
	}
}
