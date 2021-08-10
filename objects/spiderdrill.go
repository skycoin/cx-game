package objects

import (
	"fmt"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/events"
)

const DEBUG = false

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
	if DEBUG { fmt.Println("Jump: ", data) }
}

func (s spiderDrillBeforeJumpNotifier) OnSpiderDrillBeforeJump(data events.SpiderEventData) {
	if DEBUG { fmt.Println("Before Jump: ", data) }
}
