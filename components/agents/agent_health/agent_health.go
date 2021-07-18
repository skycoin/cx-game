package agent_health

import (
	"github.com/skycoin/cx-game/agents"
)

func UpdateAgents(agentList *agents.AgentList) {
	//todo right now only checks if agent is dead
	for i, agent := range agentList.Agents {
		if agent.Died() {
			agentList.DestroyAgent(i)
		}
	}
}

func Init() {

}
