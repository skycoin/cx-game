package agent_draw

import (
	"log"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/spriteloader"
)

type DrawHandler func([]*agents.Agent)
var drawHandlers [constants.NUM_DRAW_HANDLERS]DrawHandler

func Init() {
	// TODO move this
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-agent")

	RegisterDrawHandler(constants.DRAW_HANDLER_NULL, NullDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_QUAD, QuadDrawHandler)

	AssertAllDrawHandlersRegistered()
}

func AssertAllDrawHandlersRegistered() {
	for _,handler := range drawHandlers {
		if handler==nil { log.Fatalf("Did not initialize all draw handlers") }
	}
}

func RegisterDrawHandler(id constants.DrawHandlerID, handler DrawHandler) {
	drawHandlers[id] = handler
}

func ChangeDrawHandler(id constants.DrawHandlerID, newHandler DrawHandler) {
	if id < 0 || len(drawHandlers) >= int(id) {
		log.Printf("invalid draw handler with ID=%v",id)
		return
	}
	drawHandlers[id] = newHandler
}

func GetDrawHandler(id constants.DrawHandlerID) DrawHandler {
	return drawHandlers[id]
}
