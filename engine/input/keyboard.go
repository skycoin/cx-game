package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	//for actions
	keyPressed     = make(map[glfw.Key]bool)
	keyPressedDown = make(map[glfw.Key]bool)
	keysPressedUp  = make(map[glfw.Key]bool)

	modifierKey glfw.ModifierKey

	//for each context, have map of registered buttons to keys
	ButtonsMap          = make(map[InputContext]map[string]KeyComb)
	ActiveButtonsToKeys map[string]KeyComb
	lastKeyPressed      glfw.Key
)

type Axis int

type KeyComb struct {
	key         glfw.Key
	modifierKey glfw.ModifierKey
}

const (
	HORIZONTAL Axis = iota
	VERTICAL
)

func keyCallback(
	w *glfw.Window,
	key glfw.Key, scancode int, action glfw.Action, mk glfw.ModifierKey,
) {
	keyPressedDown = make(map[glfw.Key]bool)
	keysPressedUp = make(map[glfw.Key]bool)

	modifierKey = mk

	if action == glfw.Press {
		//remap game quit to combination of  buttons
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)

		}

		ProcessFlags(key, mk)

		lastKeyPressed = key
		keyPressedDown[key] = true
		keyPressed[key] = true

	} else if action == glfw.Repeat {
		//nothing
	} else if action == glfw.Release {
		modifierKey = 0
		keyPressed[key] = false
		keysPressedUp[key] = true
	}
}
