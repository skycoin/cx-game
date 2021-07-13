package agent_draw

import "github.com/skycoin/cx-game/agents"

var DrawFuncsArray []func([]*agents.Agent)

func DrawFunc1(agents []*agents.Agent) {

}

func DrawFunc2(agents []*agents.Agent) {

}

func Init() {
	DrawFuncsArray = append(DrawFuncsArray, DrawFunc1, DrawFunc2)
}
	