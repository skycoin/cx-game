package input

import "github.com/go-gl/glfw/v3.3/glfw"

type MouseCoords struct {
	X float64
	Y float64
}

var (
	window_ *glfw.Window
)

func Init(window *glfw.Window) {
	KeysPressed = make(map[glfw.Key]bool)
	ButtonsToKeys = make(map[string]glfw.Key)
	mouseCoords = &MouseCoords{}

	window_ = window
	registerCallbacks()

	MapKeyToButton("right", glfw.KeyD)
	MapKeyToButton("left", glfw.KeyA)
	MapKeyToButton("up", glfw.KeyW)
	MapKeyToButton("down", glfw.KeyS)
	MapKeyToButton("jump", glfw.KeySpace)
	MapKeyToButton("mute", glfw.KeyM)
	MapKeyToButton("freecam", glfw.KeyF2)
	MapKeyToButton("cycle-palette", glfw.KeyF3)

}

func registerCallbacks() {
	window_.SetKeyCallback(keyCallback)
	window_.SetCursorPosCallback(cursorPosCallback)
}
