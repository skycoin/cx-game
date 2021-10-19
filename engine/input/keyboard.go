package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	//for actions
	KeyPressed = make(map[glfw.Key]bool)

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
	modifierKey = mk

	if action == glfw.Press {
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}

		ProcessFlags(key, mk)

		lastKeyPressed = key
		KeyPressed[key] = true

	} else if action == glfw.Repeat {
		//nothing
	} else if action == glfw.Release {
		modifierKey = 0
		KeyPressed[key] = false
	}
}
