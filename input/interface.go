package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/utility"
)

//continuos keys, holding
func GetButton(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY IS NOT MAPPED!")
		return false
	}
	pressed, ok := KeysPressed[key]
	if !ok {
		// log.Printf("ERROR!")
		return false
	}
	return pressed
}

//action keys, if pressed once
func GetButtonDown(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY IS NOT MAPPED!")
		return false
	}
	pressed, ok := KeysPressedDown[key]
	if !ok {
		return false
	}
	KeysPressedDown[key] = false
	return pressed
}

func GetLastKey() glfw.Key {
	return lastKeyPressed
}

func GetAxis(axis Axis) float32 {
	if axis == HORIZONTAL {
		return utility.BoolToFloat(GetButton("right")) - utility.BoolToFloat(GetButton("left"))
	} else { // VERTICAL
		return utility.BoolToFloat(GetButton("up")) - utility.BoolToFloat(GetButton("down"))
	}

}

func GetMouseX() float64 {
	return mouseCoords.X
}

func GetMouseY() float64 {
	return mouseCoords.Y
}