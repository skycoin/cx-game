package agent_ai

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/models"
)

func UpdateAgents(agentlist *agents.AgentList, player *models.Player) {
	ctx := AiContext { PlayerPos: player.Pos.Mgl32() }
	for _, agent := range agentlist.Agents {
		aiHandlers[agent.AiHandlerID](agent,ctx)
	}
}
