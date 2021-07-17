package agents

import "github.com/skycoin/cx-game/constants"

type AgentList struct {
	Agents []*Agent
}

func NewAgentList() *AgentList {
	return &AgentList{
		Agents: make([]*Agent, 0),
	}
}

func NewDevAgentList() *AgentList {
	agentList := NewAgentList()
	agentList.CreateAgent(constants.AGENT_PLAYER)
	agentList.CreateAgent(constants.AGENT_ENEMY_MOB)

	return agentList
}

//  agentType - constants.AGENT_*desired type*
func (al *AgentList) CreateAgent(agentType int) bool {
	//for now
	if len(al.Agents) > constants.MAX_AGENTS {
		return false
	}
	agent := newAgent(agentType)
	al.Agents = append(al.Agents, agent)
	return true
}

func (al *AgentList) DestroyAgent(agentId int) bool {
	if agentId < 0 || agentId >= len(al.Agents) {
		return false
	}

	al.Agents = append(al.Agents[:agentId], al.Agents[agentId+1:]...)
	return false
}
