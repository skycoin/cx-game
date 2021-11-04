package inputhandler

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/common"
	"github.com/skycoin/cx-game/engine/input"
)

type AgentInputHandler struct {
	AgentControlState common.Bitset
	Mapper            InputMapper

	Actions map[input.EventInfo]func()
}

func NewAgentInputHandler() *AgentInputHandler {
	newInputHandler := AgentInputHandler{
		AgentControlState: common.NewBitSet(AGENT_KEYSTATE_LENGTH),
		Mapper:            NewAgentInputMapper(),
		Actions:           map[input.EventInfo]func(){},
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

	down, ok := input.KeyPressed[i.Mapper.GetKey(FALL_DOWN)]

	if !ok {
		down = false
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
	i.AgentControlState.SetBit(FALL_DOWN, down)
	i.AgentControlState.SetBit(JUMP, jump)
	i.AgentControlState.SetBit(CROUCH, crouch)
	i.AgentControlState.SetBit(JET_PACK, jetpack)
}

type ActionInfo struct {
	Type   input.EventType
	Key    int
	Enum   ActionEnum
	Action func()
}

func (i *AgentInputHandler) RegisterAction(actionInfo ActionInfo) {
	//assume key is already converted to int
	if actionInfo.Type == input.KEY_DOWN {
		i.Mapper.BindKey(glfw.Key(actionInfo.Key), int(actionInfo.Enum))
	} else if actionInfo.Type == input.KEY_UP {
		log.Fatalf("NOT IMPLEMENTED!")
	}

	if actionInfo.Action == nil {
		log.Panicln("ACTION SHOULD NOT BE NIL")
	}
	i.Actions[input.EventInfo{
		Type:   actionInfo.Type,
		Button: actionInfo.Key,
	}] = actionInfo.Action
}

func (i *AgentInputHandler) ProcessEvents() {
	if len(input.InputEvents) == 0 {
		return
	}
	for _, event := range input.InputEvents {
		action, ok := i.Actions[event]
		if ok {
			action()
			// log.Fatal("Not Registered event! ", event.Type)
		}

	}

}
