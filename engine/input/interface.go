package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
)

var DEBUG = false

//continuos keys, holding
func GetButton(button string) bool {
	keyComb, ok := ActiveButtonsToKeys[button]
	if !ok && DEBUG {
		log.Printf("BUTTON [%s] IS NOT MAPPED!\n", button)
		return false
	}
	pressed, ok := KeyPressed[keyComb.key]
	if !ok {
		// log.Printf("ERROR!")
		return false
	}
	if modifierKey == keyComb.modifierKey {
		return pressed
	}
	return false
}

func GetKey(key glfw.Key) bool {
	if modifierKey == 0 {
		return KeyPressed[key]
	}
	return false
}

func GetKeyIsUp(key glfw.Key) bool {
	return window_.GetKey(key) == glfw.Press
}

func GetLastKey() glfw.Key {
	key := lastKeyPressed
	//consume and set to false
	lastKeyPressed = glfw.KeyUnknown
	return key
}

func GetAxis(axis Axis) float32 {
	if axis == HORIZONTAL {
		return cxmath.BoolToFloat(GetButton("right")) - cxmath.BoolToFloat(GetButton("left"))
	} else { // VERTICAL
		return cxmath.BoolToFloat(GetButton("up")) - cxmath.BoolToFloat(GetButton("down"))
	}

}

func GetMouseX() float32 {
	return float32(MouseCoords.X)
}
func GetMouseY() float32 {
	return float32(MouseCoords.Y)
}

func GetMousePos() mgl32.Vec2 {
	physicalX := float32(MouseCoords.X)
	physicalY := float32(MouseCoords.Y)

	physicalPos := mgl32.Vec2{physicalX, physicalY}
	physicalPosHomogenous :=
		mgl32.Vec4{physicalPos.X(), physicalPos.Y(), 0, 1}

	transform := physicalToViewPortTransform
	virtualPos := transform.Mul4x1(physicalPosHomogenous).Vec2()

	return virtualPos
}

func Reset() {
	//reset lastkeyPressed
	return
}
