package agent_physics

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

type PhysicsHandler func(agent *agents.Agent)

var PhysicsHandlerList [constants.NUM_AGENT_PHYSICS_HANDLERS]PhysicsHandler

func Init() {

}

func RegisterPhysicsHandler(id types.AgentPhysicsHandlerID, handler PhysicsHandler) {
	PhysicsHandlerList[id] = handler
}
