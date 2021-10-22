package player_agent

import (
	"log"

	"github.com/skycoin/cx-game/common"
	"github.com/skycoin/cx-game/components/agents"
)

type PlayerAgent struct {
	agent *agents.Agent
}

func NewPlayerAgent() *PlayerAgent {
	newPlayerAgent := PlayerAgent{}
	return &newPlayerAgent
}
func (p *PlayerAgent) SetPlayerAgent(agent *agents.Agent) {
	p.agent = agent
}

func (p *PlayerAgent) GetAgent() *agents.Agent {
	if p.agent == nil {
		log.Fatalln("Player agent is not set!")
	}
	return p.agent
}

func (p *PlayerAgent) SetControlState(cs common.Bitset) {
	p.agent.SetControlState(cs)
	// assertAgentNotNil()
}

// func assertAgentNotNil() {
// 	if playerAgent == nil {
// 		log.Fatalln("EXPECT AGENT TO NOT BE NULL")
// 	}
// }
