package agents

import (
	"github.com/skycoin/cx-game/physics"
)

type Agent struct {
	AgentType            int
	AiHandlerId          int
	PhysicsState         physics.Body
	PhysicsParameters    physics.PhysicsParameters
	DrawFunctionUpdateId int
	HealthComponent      HealthComponent
}

type HealthComponent struct {
	Health_amount int
	Health_max    int
	Died          bool
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

func (a *Agent) TakeDamage(amount int) {
	a.HealthComponent.Health_amount -= amount
	if a.HealthComponent.Health_amount <= 0 {
		a.HealthComponent.Died = true
	}
}

func (a *Agent) Died() bool {
	return a.HealthComponent.Died
}
