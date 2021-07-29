package agent_ai

import (
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/world"
)

// TODO since AI will be run server-side, 
// we should not be passing any particular player to this function.
// Rather, the list of players should be computed
// by filtering the world agents
func UpdateAgents(World *world.World, player *models.Player) {
	ctx := AiContext { PlayerPos: player.Pos.Mgl32() }
	for _, agent := range World.Entities.Agents.Get() {
		aiHandlers[agent.AiHandlerID](agent,ctx)
	}
}
