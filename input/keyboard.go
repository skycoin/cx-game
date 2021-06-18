package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	//for actions
	KeysPressed     map[glfw.Key]bool
	KeysPressedDown map[glfw.Key]bool
	//for
	ButtonsToKeys  map[string]glfw.Key
	lastKeyPressed glfw.Key
)

type Axis int

const (
	HORIZONTAL Axis = iota
	VERTICAL
)

func Tick() {
	for key := range KeysPressed {
		KeysPressed[key] = false
	}
}

var repeatCounter float64 = 0

func keyCallback(w *glfw.Window, key glfw.Key, s int, action glfw.Action, mk glfw.ModifierKey) {

	if action == glfw.Press {
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)

		}
		lastKeyPressed = key
		KeysPressedDown[key] = true
		KeysPressed[key] = true
	} else if action == glfw.Repeat {
	} else if action == glfw.Release {
		KeysPressed[key] = false
	}
}

func MapKeyToButton(button string, key glfw.Key) {
	ButtonsToKeys[button] = key
}

func registerKeyMaps() {
	MapKeyToButton("right", glfw.KeyD)
	MapKeyToButton("left", glfw.KeyA)
	MapKeyToButton("up", glfw.KeyW)
	MapKeyToButton("down", glfw.KeyS)
	MapKeyToButton("jump", glfw.KeySpace)
	MapKeyToButton("mute", glfw.KeyM)
	MapKeyToButton("freecam", glfw.KeyF2)
	MapKeyToButton("cycle-palette", glfw.KeyF3)
	MapKeyToButton("scratch", glfw.KeyLeftShift)
	MapKeyToButton("inventory-grid", glfw.KeyI)
}
