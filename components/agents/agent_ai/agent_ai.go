package agent_ai

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/models"
)

func UpdateAgents(agentlist *agents.AgentList, player *models.Player) {
	for _, agent := range agentlist.Agents {
		AiHandlerList[agent.AiHandlerId](agent)
	}
}
