package agents

import (
	"github.com/skycoin/cx-game/physics"
)

type Agent struct {
	//physics state
	PhysicsState physics.Body
	//physics parameters
	PhysicsParameters physics.PhysicsParameters

	//movementstate
	DrawFunctionUpdateId int
	AgentType            int
}

func newAgent(agentType int) *Agent {
	agent := Agent{
		PhysicsState:         physics.Body{},
		PhysicsParameters:    physics.PhysicsParameters{Radius: 5},
		DrawFunctionUpdateId: 1,
		AgentType:            agentType,
	}

	return &agent
}
