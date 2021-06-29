package agents

import (
	"github.com/skycoin/cx-game/physics"
)

type Agent struct {
	// InventoryId uint32
	AgentMeta
	physics.Body
	AgentType AgentType
	AgentID   int32
	// CollisionWidth  float32
	// CollisionHeight float32
}

type AgentType int

var (
	agentIdCounter int32 = 0
)

func newAgent(agentType AgentType) *Agent {
	agentIdCounter += 1
	// return &Agent{
	// 	AgentID:     agentIdCounter,
	// 	AgentType:   agentType,
	// 	InventoryId: inventoryId,
	// }
	return &Agent{}
}

func (a *Agent) FixedTick() {
	//move the agent
	//resolve collisions
	a.Pos = a.Pos.Add(a.Vel.Mult(physics.TimeStep))

}

func (a *Agent) Draw() {
	//interpolate position between physics ticks and draw the agent
}
