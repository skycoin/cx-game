package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	//for actions
	KeysPressed     map[glfw.Key]bool
	KeysPressedDown map[glfw.Key]bool
	KeysPressedUp   map[glfw.Key]bool
	//for
	ButtonsToKeys  map[string]glfw.Key
	lastKeyPressed glfw.Key
)

type Axis int

const (
	HORIZONTAL Axis = iota
	VERTICAL
)

func Reset() {
	// for key := range KeysPressed {
	// 	KeysPressed[key] = false
	// }

}

var counter int

func keyCallback(w *glfw.Window, key glfw.Key, s int, action glfw.Action, mk glfw.ModifierKey) {
	for key := range KeysPressedDown {
		KeysPressedDown[key] = false
	}
	for key := range KeysPressedUp {
		KeysPressedUp[key] = false
	}

	if action == glfw.Press {
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

func MapKeyToButton(button string, key glfw.Key) {
	ButtonsToKeys[button] = key
}
