package game

import (
	"github.com/skycoin/cx-game/components/agents"
)

// logic that only makes sense from client's perspective

var playerAgentID int

func findPlayer() *agents.Agent {
	return World.Entities.Agents.FromID(playerAgentID)
}
