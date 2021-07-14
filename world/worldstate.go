package world

import "github.com/skycoin/cx-game/agents"

type WorldState struct {
	AgentList *agents.AgentList
	//particle list
}

func NewDevWorldState() *WorldState {
	worldState := WorldState{
		AgentList: agents.NewDevAgentList(),
	}

	return &worldState
}
