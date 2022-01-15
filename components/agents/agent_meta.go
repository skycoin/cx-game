package agents

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/physics"
)

// type AgentMeta struct {
// 	MovementSpeed    float32
// 	PlayerControlled bool
// 	AgentID          int
// 	AgentType        int
// 	AiControllerId   int
// 	PosX             float32
// 	PosY             float32
// 	MomentumX        float32
// 	MomentumY        float32
// 	MinJumpHeight    float32
// 	MaxJumpHeight    float32
// 	// StaticFriction float32
// 	DynamicFriction float32
// }

type AgentMeta struct {
	Category          types.AgentCategory
	Type              types.AgentTypeID
	PhysicsParameters physics.PhysicsParameters
	PlayerData
	SpineData
}

//agent struct holds
//
