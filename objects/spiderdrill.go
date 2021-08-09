package objects

import (
	"fmt"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/events"
)

type spiderDrillJumpNotifier struct {
	agent *agents.Agent
}

type spiderDrillBeforeJumpNotifier struct {
	agent *agents.Agent
}

func SpiderDrillInit() {
	createNotifier := spiderDrillJumpNotifier{
		agent: nil,
	}

	createBeforeJumpNotifier := spiderDrillBeforeJumpNotifier{
		agent: nil,
	}

	events.OnSpiderJump.Register(createNotifier)
	events.OnSpiderBeforeJump.Register(createBeforeJumpNotifier)
}

func (s spiderDrillJumpNotifier) OnSpiderDrillJump(data events.SpiderEventData) {
	fmt.Println("Jump: ", data)
}

func (s spiderDrillBeforeJumpNotifier) OnSpiderDrillBeforeJump(data events.SpiderEventData) {
	fmt.Println("Before Jump: ", data)
}
