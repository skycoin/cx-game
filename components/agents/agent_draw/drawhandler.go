package agent_draw

import "github.com/skycoin/cx-game/agents"

var DrawHandlerList []func([]*agents.Agent)

func DrawHandler_1(agents []*agents.Agent) {

}

func DrawHandler_2(agents []*agents.Agent) {

}

func Init() {
	DrawHandlerList = append(DrawHandlerList, DrawHandler_1, DrawHandler_2)
}
	