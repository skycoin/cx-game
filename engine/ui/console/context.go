package console

import (
	"github.com/skycoin/cx-game/world"
)

type CommandContext struct {
	World *world.World
}

func NewCommandContext() CommandContext {
	return CommandContext {}
}
