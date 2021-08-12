package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	
	"github.com/skycoin/cx-game/cxmath"
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
		log.Printf("KEY [%s] IS NOT MAPPED!", button)
		return false
	}
	pressed, ok := KeysPressedDown[key]
	if !ok {
		return false
	}
	KeysPressedDown[key] = false
	return pressed
}

func GetButtonUp(button string) bool {
	key, ok := ButtonsToKeys[button]
	if !ok {
		log.Printf("KEY IS NOT MAPPED")
		return false
	}
	return GetKeyDown(key)
}

func GetKey(key glfw.Key) bool {
	return KeysPressed[key]
}

func GetKeyDown(key glfw.Key) bool {
	pressed, ok := KeysPressedUp[key]
	if !ok {
		return false
	}
	KeysPressedUp[key] = false
	return pressed
}

func GetKeyIsUp(key glfw.Key) bool {
	return window_.Window.GetKey(key) == glfw.Press
}

func GetLastKey() glfw.Key {
	return lastKeyPressed
}

func GetAxis(axis Axis) float32 {
	if axis == HORIZONTAL {
		return cxmath.BoolToFloat(GetButton("right")) - cxmath.BoolToFloat(GetButton("left"))
	} else { // VERTICAL
		return cxmath.BoolToFloat(GetButton("up")) - cxmath.BoolToFloat(GetButton("down"))
	}

}

func GetMouseX() float32 {
	return float32(mouseCoords.X)
}
func GetMouseY() float32 {
	return float32(mouseCoords.Y)
}

func GetMousePos() mgl32.Vec2 {
	physicalX := float32(mouseCoords.X)
	physicalY := float32(mouseCoords.Y)

	physicalPos := mgl32.Vec2 { physicalX, physicalY }
	physicalPosHomogenous :=
		mgl32.Vec4 { physicalPos.X(), physicalPos.Y(), 0, 1 }

	transform := window_.PhysicalToViewportTransform
	virtualPos := transform.Mul4x1(physicalPosHomogenous).Vec2()

	return virtualPos
}
