package world

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/components/particles"
)

type WorldState struct {
	AgentList *agents.AgentList
	//particle list
	ParticleList *particles.ParticleList
}

func NewDevWorldState() *WorldState {
	worldState := WorldState{
		AgentList:    agents.NewDevAgentList(),
		ParticleList: &particles.ParticleList{},
	}

	return &worldState
}
