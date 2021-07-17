package agent_draw

import (
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
)

func QuadDrawHandler(agents []*agents.Agent) {
	// TODO is this assumed??? can we omit this check?
	if len(agents)==0 { return }
	spriteID := getSpriteID(agents[0].AgentType)
	for _, agent := range agents {
		spriteloader.DrawSpriteQuad(
			agent.PhysicsState.Pos.X,
			agent.PhysicsState.Pos.Y,
			agent.PhysicsState.Size.X,
			agent.PhysicsState.Size.Y,
			spriteID,
		)
	}
}

func getSpriteID(agentType constants.AgentType) spriteloader.SpriteID {
	switch agentType {
	default:
		return spriteloader.GetSpriteIdByName("basic-agent")
	}
}
