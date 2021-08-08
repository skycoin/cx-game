package objects

import (
	"fmt"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/events"
)

type spiderDrillJumpNotifier struct {
	agent *agents.Agent
}

func SpiderDrillInit() {
	createNotifier := spiderDrillJumpNotifier{
		agent: nil,
	}

	events.OnSpiderJump.Register(createNotifier)
}

func (s spiderDrillJumpNotifier) OnSpiderDrillJump(data events.SpiderEventData) {
	fmt.Println("Jump: ", data)
}
