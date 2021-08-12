package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	mouseCoords  MouseCoords
	widthOffset  int32
	heightOffset int32
	scale        float32 = 1
)

func cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	mouseCoords.X = xpos
	mouseCoords.Y = ypos

}

func UpdateMouseCoords(widthOff, heightOff int32, scl float32) {
	scale = scl
	widthOffset = widthOff
	heightOffset = heightOff
}
