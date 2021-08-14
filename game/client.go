package game

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
)

// logic that only makes sense from client's perspective

var playerAgentID types.AgentID

func findPlayer() *agents.Agent {
	return World.Entities.Agents.FromID(playerAgentID)
}
