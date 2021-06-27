package input

//https://www.gamedev.net/blogs/entry/2250186-designing-a-robust-input-handling-system-for-games/

/*
	Performance is important; input lag is a bad thing.
	It should be easy to have new systems tap into the input stream.
	The system must be very flexible and capable of handling a wide variety of game situations.
	Configurability (input mapping) is essential for modern games.
*/

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
	KeysPressedDown = make(map[glfw.Key]bool)
	KeysPressedUp = make(map[glfw.Key]bool)
	ButtonsToKeys = make(map[string]glfw.Key)

	mouseCoords = &MouseCoords{}

	window_ = window
	registerCallbacks()

	registerKeyMaps()

}

func registerCallbacks() {
	window_.SetKeyCallback(keyCallback)
	window_.SetCursorPosCallback(cursorPosCallback)
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
	MapKeyToButton("fly", glfw.KeyT)
	MapKeyToButton("crouch", glfw.KeyC)
	MapKeyToButton("action", glfw.KeyE)
}
