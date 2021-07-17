package components

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/agents/ai"
	"github.com/skycoin/cx-game/components/agents/draw"
	"github.com/skycoin/cx-game/components/agents/health"
	"github.com/skycoin/cx-game/components/agents/physics"
	"github.com/skycoin/cx-game/world"
)

func Update(dt float32) {
	//update health state first
	agent_health.UpdateAgents(currentWorldState.AgentList)
	//update physics state second
	agent_physics.UpdateAgents(currentWorldState, currentPlanet)

	agent_ai.UpdateAgents(currentWorldState.AgentList, currentPlayer)
}

func Draw(worldState *world.WorldState, cam *camera.Camera) {
	agent_draw.DrawAgents(worldState.AgentList, cam)
}
