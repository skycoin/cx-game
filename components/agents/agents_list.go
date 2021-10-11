package agents

import (
	"log"

	"github.com/skycoin/cx-game/common"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type AgentList struct {
	// profile this to see if reducing indirection
	// would help with performance in a significant way
	Agents []*Agent
	//to hold all freed ids
	idQueue common.QueueI
}

func NewAgentList() *AgentList {
	return &AgentList{
		Agents: make([]*Agent, 0),
	}
}

func NewDevAgentList() *AgentList {
	agentList := NewAgentList()

	return agentList
}

func (al *AgentList) AddAgent(agent *Agent) bool {
	//for now
	if len(al.Agents) > constants.MAX_AGENTS {
		// return false
		log.Fatalln("TOO MUCH AGENTS")
	}
	newId := al.idQueue.Pop()

	if newId == -1 {
		al.Agents = append(al.Agents, agent)
		agent.AgentId = types.AgentID(len(al.Agents) - 1)
	} else {
		agent.AgentId = types.AgentID(newId)
		al.Agents[newId] = agent
	}
	return true
}

//throw log.fatal at errors, rather than returning false
func (al *AgentList) DestroyAgent(agentId types.AgentID) {
	if int(agentId) < 0 || int(agentId) >= len(al.Agents) {
		log.Fatalln("ERROR AT AGENT DELETION")
	}

	al.Agents[agentId] = nil
	al.idQueue.Push(int(agentId))
}

func (al *AgentList) Spawn(
	agentTypeID types.AgentTypeID, opts AgentCreationOptions,
) types.AgentID {
	agent := GetAgentType(agentTypeID).CreateAgent(opts)
	agent.FillDefaults()
	agent.Validate()
	al.AddAgent(agent)
	return types.AgentID(agent.AgentId)
}

func (al *AgentList) GetAllAgents() []*Agent { return al.Agents }

// TODO something better than linear search
func (al *AgentList) FromID(id types.AgentID) *Agent {
	agent := al.Agents[id]
	if agent == nil {
		log.Fatalln("Expect the agent to not be nil")
	}
	return agent
}

func (al *AgentList) TileIsClear(x, y int) bool {
	for _, agent := range al.GetAllAgents() {
		if agent == nil {
			continue
		}
		if agent.Transform.Contains(float32(x), float32(y), 0.5, 0.5) {
			return false
		}
	}
	return true
}

func (al *AgentList) GetFirstSpiderDrill() *Agent {
	for _, agent := range al.Agents {
		if agent == nil {
			continue
		}
		if agent.Meta.Type == constants.AGENT_TYPE_SPIDER_DRILL {
			return agent
		}
	}
	return nil
}
