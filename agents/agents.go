package agents

import (
	"github.com/skycoin/cx-game/physics"
)

type Agent struct {
	//physics state
	PhysicsState physics.Body
	//physics parameters
	PhysicsParameters physics.PhysicsParameters
	AiHandlerId    int
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

//prefabs

func (a *Agent) SetPosition(x, y float32) {
	a.PhysicsState.Pos.X = x
	a.PhysicsState.Pos.Y = y
}

func (a *Agent) SetSize(x, y float32) {
	a.PhysicsState.Size.X = x
	a.PhysicsState.Size.Y = y
}

func (a *Agent) SetVelocity(x, y float32) {
	a.PhysicsState.Vel.X = x
	a.PhysicsState.Vel.Y = y
}
