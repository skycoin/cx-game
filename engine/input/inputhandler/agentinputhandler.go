package inputhandler

import (
	"github.com/skycoin/cx-game/common"
	"github.com/skycoin/cx-game/engine/input"
)

type AgentInputHandler struct {
	AgentControlState common.Bitset
	Mapper            InputMapper
}

func NewAgentInputHandler() *AgentInputHandler {
	newInputHandler := AgentInputHandler{
		AgentControlState: common.NewBitSet(AGENT_KEYSTATE_LENGTH),
		Mapper:            NewAgentInputMapper(),
	}
	return &newInputHandler
}

//poll at physics rate
func (i *AgentInputHandler) UpdateKeyState() {
	var left, right, jump, crouch, jetpack bool
	//todo have customizable keybindings
	left, ok := input.KeyPressed[i.Mapper.GetKey(MOVE_LEFT)]
	if !ok {
		left = false
	}
	right, ok = input.KeyPressed[i.Mapper.GetKey(MOVE_RIGHT)]
	if !ok {
		right = false
	}

	jump, ok = input.KeyPressed[i.Mapper.GetKey(JUMP)]
	if !ok {
		jump = false
	}
	crouch, ok = input.KeyPressed[i.Mapper.GetKey(CROUCH)]
	if !ok {
		crouch = false
	}
	jetpack, ok = input.KeyPressed[i.Mapper.GetKey(JET_PACK)]
	if !ok {
		jetpack = false
	}

	i.AgentControlState.SetBit(MOVE_LEFT, left)
	i.AgentControlState.SetBit(MOVE_RIGHT, right)
	i.AgentControlState.SetBit(JUMP, jump)
	i.AgentControlState.SetBit(CROUCH, crouch)
	i.AgentControlState.SetBit(JET_PACK, jetpack)
}
