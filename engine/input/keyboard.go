package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	//for actions
	KeysPressed     = make(map[glfw.Key]bool)
	KeysPressedDown = make(map[glfw.Key]bool)
	KeysPressedUp   = make(map[glfw.Key]bool)

	//for each context, have map of registered buttons to keys
	ButtonsMap          = make(map[InputContext]map[string]glfw.Key)
	ActiveButtonsToKeys map[string]glfw.Key
	lastKeyPressed      glfw.Key
)

type Axis int

const (
	HORIZONTAL Axis = iota
	VERTICAL
)

func keyCallback(
	w *glfw.Window,
	key glfw.Key, scancode int, action glfw.Action, mk glfw.ModifierKey,
) {
	KeysPressedDown = make(map[glfw.Key]bool)
	KeysPressedUp = make(map[glfw.Key]bool)

	if action == glfw.Press {
		//remap game quit to combination of  buttons
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)

		}
		lastKeyPressed = key
		KeysPressedDown[key] = true
		KeysPressed[key] = true

	} else if action == glfw.Repeat {
		//nothing
	} else if action == glfw.Release {
		KeysPressed[key] = false
		KeysPressedUp[key] = true
	}
}
