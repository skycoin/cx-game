package input

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
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

func GetMouseWorldCoordsX() float32 {
	screenX := GetScreenX()
	tileX := camX + screenX/render.PixelsPerTile
	return tileX
}

func GetMouseWorldCoordsY() float32 {
	screenY := GetScreenY()
	tileY := camY + screenY/render.PixelsPerTile
	return tileY
}

func GetMouseWorldCoords() cxmath.Vec2 {
	return cxmath.Vec2{
		X: GetMouseWorldCoordsX(),
		Y: GetMouseWorldCoordsY(),
	}
}
func GetScreenX() float32 {
	screenX := ((float32(mouseCoords.X)-float32(widthOffset))/float32(scale) - float32(window_.Width)/2) / camZoom
	return screenX
}

func GetScreenY() float32 {
	screenY := (((float32(mouseCoords.Y)-float32(heightOffset))/float32(scale) - float32(window_.Height)/2) * -1) / camZoom
	return screenY
}
