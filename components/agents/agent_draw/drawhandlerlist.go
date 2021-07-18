package agent_draw

import (
	"log"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/spriteloader"
)

type AgentDrawHandler func([]*agents.Agent)

var drawHandlers [constants.NUM_AGENT_DRAW_HANDLERS]AgentDrawHandler

func Init() {
	// TODO move this
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-agent")

	RegisterDrawHandler(constants.DRAW_HANDLER_NULL, NullDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_QUAD, QuadDrawHandler)

	AssertAllDrawHandlersRegistered()
}

func AssertAllDrawHandlersRegistered() {
	for _, handler := range drawHandlers {
		if handler == nil {
			log.Fatalf("Did not initialize all agent draw handlers")
		}
	}
}

func RegisterDrawHandler(id types.AgentDrawHandlerID, handler AgentDrawHandler) {
	drawHandlers[id] = handler
}

func ChangeDrawHandler(id types.AgentDrawHandlerID, newHandler AgentDrawHandler) {
	if id < 0 || len(drawHandlers) >= int(id) {
		log.Printf("invalid draw handler with ID=%v", id)
		return
	}
	drawHandlers[id] = newHandler
}

func GetDrawHandler(id types.AgentDrawHandlerID) AgentDrawHandler {
	return drawHandlers[id]
}
