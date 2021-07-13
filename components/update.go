package components

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/agents/agent_physics"
	"github.com/skycoin/cx-game/world"
)

func Update(dt float32) {
	//call various functions
	agent_physics.UpdateAgents(currentWorldState, currentPlanet)
	agent_ai.UpdateAgents(currentWorldState.AgentList, currentPlayer)
}

func Draw(worldState *world.WorldState, cam *camera.Camera) {
	agent_draw.DrawAgents(worldState.AgentList, cam)
}
