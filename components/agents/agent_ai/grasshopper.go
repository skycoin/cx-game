package agent_ai

import "github.com/skycoin/cx-game/components/agents"

func AiHandlerGrassHopper(agent *agents.Agent, ctx AiContext) {
	if agent.PhysicsState.Collisions.Below {
		agent.PhysicsState.Vel.X = 0
	}
}
