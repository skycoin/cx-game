package inputhandler

import "github.com/go-gl/glfw/v3.3/glfw"

const (
	//states
	MOVE_LEFT = iota
	MOVE_RIGHT
	FALL_DOWN
	JUMP
	CROUCH
	JET_PACK

	//len
	AGENT_KEYSTATE_LENGTH

	// //actions
	// ENABLE_BACKGROUND
)

type ActionEnum int

const (
	ENABLE_BACKGROUND ActionEnum = 20 + iota
	DISABLE_BACKGROUND
	TOGGLE_FREECAM
)

func NewAgentInputMapper() InputMapper {
	newInputMapper := InputMapper{
		enumToKey: make(map[int]glfw.Key),
		keyToEnum: make(map[glfw.Key]int),
	}

	newInputMapper.BindKey(glfw.KeyA, MOVE_LEFT)
	newInputMapper.BindKey(glfw.KeyD, MOVE_RIGHT)
	newInputMapper.BindKey(glfw.KeyS, FALL_DOWN)
	newInputMapper.BindKey(glfw.KeySpace, JUMP)
	newInputMapper.BindKey(glfw.KeyLeftShift, CROUCH)
	newInputMapper.BindKey(glfw.KeyT, JET_PACK)

	return newInputMapper
}
