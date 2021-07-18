package agents

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/physics"
)

type Agent struct {
	AgentType         constants.AgentType
	AiHandlerID       types.AgentAiHandlerID
	PhysicsState      physics.Body
	PhysicsParameters physics.PhysicsParameters
	DrawHandlerID     types.AgentDrawHandlerID
	HealthComponent   HealthComponent
}

type HealthComponent struct {
	Health_amount int
	Health_max    int
	Died          bool
}

func newAgent() *Agent {
	agent := Agent{
		AgentType:         constants.AGENT_UNDEFINED,
		AiHandlerID:       constants.AI_HANDLER_NULL,
		DrawHandlerID:     constants.DRAW_HANDLER_NULL,
		PhysicsState:      physics.Body{},
		PhysicsParameters: physics.PhysicsParameters{Radius: 5},
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
