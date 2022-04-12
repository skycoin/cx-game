package agent_draw

import (
	"log"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/camera"
)

type DrawHandlerContext struct {
	// TODO maybe pass in framebuffers here???
	Camera *camera.Camera
	Dt     float32
}

type AgentDrawHandler func([]*agents.Agent, DrawHandlerContext)

var drawHandlers [constants.NUM_AGENT_DRAW_HANDLERS]AgentDrawHandler

func Init() {
	// TODO move this
	//spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-agent")

	RegisterDrawHandler(constants.DRAW_HANDLER_NULL, NullDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_QUAD, QuadDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_ANIM, AnimatedDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_PLAYER, PlayerDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_COLOR, ColorDrawHandler)
	RegisterDrawHandler(constants.DRAW_HANDLER_SPINE, SpinePlayerDrawHandler)

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
