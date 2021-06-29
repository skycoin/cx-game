package agents

import (
	"github.com/skycoin/cx-game/physics"
)

type AgentList struct {
	Agents []*Agent
}

var (
	accumulator float32
)

func NewAgentList() *AgentList {
	return &AgentList{
		Agents: make([]*Agent, MAX_AGENTS),
	}
}

func NewDevAgentList() *AgentList {
	agentList := NewAgentList()
	agentList.CreateAgent(PLAYER)
	agentList.CreateAgent(ENEMY_MOB)

	return agentList
}
func (al *AgentList) CreateAgent(agentType AgentType) bool {
	//for now
	if len(al.Agents) > MAX_AGENTS {
		return false
	}
	agent := newAgent(agentType)
	al.Agents = append(al.Agents, agent)
	return true
}

func (al *AgentList) DestroyAgent(agentId int32) bool {
	if agentId < 0 || agentId >= int32(len(al.Agents)) {
		return false
	}
	for i, v := range al.Agents {
		if v.AgentID == agentId {
			al.Agents = append(al.Agents[:i], al.Agents[i+1:]...)
			return true
		}
	}
	return false
}

func (al *AgentList) Draw() {
	for _, agent := range al.Agents {
		agent.Draw()
	}
}

func (al *AgentList) Tick(dt float32) {
	accumulator += dt

	//for physics
	for accumulator >= physics.TimeStep {
		for _, agent := range al.Agents {
			agent.FixedTick()
		}
		accumulator -= physics.TimeStep
	}
}
