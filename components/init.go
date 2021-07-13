package components

import (
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_draw"
)

func Init() {
	agent_draw.Init()
	agent_ai.Init()
}
