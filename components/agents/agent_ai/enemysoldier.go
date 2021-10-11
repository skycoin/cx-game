package agent_ai

import "github.com/skycoin/cx-game/components/agents"

func AiHandlerEnemySoldier(agent *agents.Agent, ctx AiContext) {
	if agent.Transform.Collisions.Below {
		agent.Transform.Vel.X = 0
	}
}
