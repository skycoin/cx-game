package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	mouseCoords MouseCoords
)

func cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	mouseCoords.X = xpos
	mouseCoords.Y = ypos
}
