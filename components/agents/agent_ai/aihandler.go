package agent_ai

import "github.com/skycoin/cx-game/agents"

var AiHandlerList []func(*agents.Agent)

func AiHandler_1(*agents.Agent) {
	//defines behavior
}

func AiHandler_2(*agents.Agent) {
	//defines behavior
}

func Init() {
	AiHandlerList = append(AiHandlerList, AiHandler_1, AiHandler_2)
}
