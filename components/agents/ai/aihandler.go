package agent_ai

import "github.com/skycoin/cx-game/agents"

var AiHandlerList []func(*agents.Agent)

func Init() {
	RegisterAiHandler(AiHandler_1)
	RegisterAiHandler(AiHandler_2)
}

func RegisterAiHandler(newHandler func(*agents.Agent)) {
	AiHandlerList = append(AiHandlerList, newHandler)
}
