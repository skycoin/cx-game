package console

import (
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/components/agents"
)

type CommandContext struct {
	World *world.World
	Player *agents.Agent
}

func NewCommandContext() CommandContext {
	return CommandContext {}
}
