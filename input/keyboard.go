package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	KeysPressed    map[glfw.Key]bool
	ButtonsToKeys  map[string]glfw.Key
	lastKeyPressed glfw.Key
)

func keyCallback(w *glfw.Window, key glfw.Key, s int, action glfw.Action, mk glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)

		}
		lastKeyPressed = key
		KeysPressed[key] = true
	} else if action == glfw.Release {
		KeysPressed[key] = false
	}
}

func MapKeyToButton(button string, key glfw.Key) {
	ButtonsToKeys[button] = key
}

func GetButton(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY IS NOT MAPPED!")
		return false
	}
	pressed, ok := KeysPressed[key]
	if !ok {
		log.Printf("ERROR!")
		return false
	}
	return pressed
}

func GetLastKey() glfw.Key {
	return lastKeyPressed
}
