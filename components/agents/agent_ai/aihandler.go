package agent_ai

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type AiContext struct {
	PlayerPos mgl32.Vec2
}
type AiHandler func(*agents.Agent, AiContext)

var aiHandlers [constants.NUM_AI_HANDLERS]AiHandler

func Init() {
	RegisterAiHandler(constants.AI_HANDLER_NULL, AiHandlerNull)
	RegisterAiHandler(constants.AI_HANDLER_WALK, AiHandlerWalk)
	RegisterAiHandler(constants.AI_HANDLER_LEAP, AiHandlerLeap)
	RegisterAiHandler(constants.AI_HANDLER_DRILL, AiHandlerDrill)

	assertAllAiHandlersRegistered()
}

func assertAllAiHandlersRegistered() {
	for _, handler := range aiHandlers {
		if handler == nil {
			log.Fatalf("Did not initialize all agent ai handlers")
		}
	}
}

func RegisterAiHandler(id types.AgentAiHandlerID, newHandler AiHandler) {
	aiHandlers[id] = newHandler
}
