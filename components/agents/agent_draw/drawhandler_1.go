package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
)

func DrawHandler_1(agents []*agents.Agent) {
	for _, agent := range agents {
		spriteloader.DrawSpriteQuad(
			agent.PhysicsState.Pos.X,
			agent.PhysicsState.Pos.Y,
			agent.PhysicsState.Size.X,
			agent.PhysicsState.Size.Y,
			getSpriteId(agent.AgentType),
		)
	}
}

func getSpriteId(agentType int) spriteloader.SpriteID {
	switch agentType {
	default:
		return spriteloader.GetSpriteIdByName("basic-agent")
	}
}
