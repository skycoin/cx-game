package agents

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type AgentList struct {
	// profile this to see if reducing indirection
	// would help with performance in a significant way
	Agents []*Agent
}

func NewAgentList() *AgentList {
	return &AgentList{
		Agents: make([]*Agent, 0),
	}
}

func NewDevAgentList() *AgentList {
	agentList := NewAgentList()
	// player := newAgent(len(agentList.Agents))
	// player.AgentCategory = constants.AGENT_CATEGORY_PLAYER
	// agentList.CreateAgent(player)
	// enemy := newAgent(len(agentList.Agents))
	// enemy.AgentCategory = constants.AGENT_CATEGORY_ENEMY_MOB
	// agentList.CreateAgent(enemy)

	return agentList
}

func (al *AgentList) CreatelAgent(agent *Agent) bool {
	//for now
	if len(al.Agents) > constants.MAX_AGENTS {
		return false
	}
	al.Agents = append(al.Agents, agent)
	return true
}

func (al *AgentList) DestroyAgent(agentId int) bool {
	if agentId < 0 || agentId >= len(al.Agents) {
		return false
	}

	al.Agents = append(al.Agents[:agentId], al.Agents[agentId+1:]...)
	return true
}

func (al *AgentList) Spawn(
	agentTypeID constants.AgentTypeID, opts AgentCreationOptions,
) types.AgentID {
	agent := GetAgentType(agentTypeID).CreateAgent(opts)
	agent.FillDefaults()
	agent.Validate()
	agent.AgentId = types.AgentID(len(al.Agents))
	al.Agents = append(al.Agents, agent)
	return types.AgentID(agent.AgentId)
}

func (al *AgentList) Get() []*Agent { return al.Agents }

// TODO something better than linear search
func (al *AgentList) FromID(id types.AgentID) *Agent { 
	for _,agent := range al.Get() {
		if agent.AgentId == id { return agent }
	}
	return nil
}

func (al *AgentList) TileIsClear(x, y int) bool {
	for _, agent := range al.Get() {
		if agent.PhysicsState.Contains(float32(x), float32(y)) {
			return false
		}
	}
	return true
}
