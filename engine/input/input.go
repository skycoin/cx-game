package input

//https://www.gamedev.net/blogs/entry/2250186-designing-a-robust-input-handling-system-for-games/

/*
	Performance is important; input lag is a bad thing.
	It should be easy to have new systems tap into the input stream.
	The system must be very flexible and capable of handling a wide variety of game situations.
	Configurability (input mapping) is essential for modern games.
*/

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	window_          *glfw.Window
	inputContext     InputContext
	prevInputContext InputContext
	// camZoom          float32 = 1
	// camX, camY       float32
)

type InputContext uint8

const (
	GAME InputContext = iota
	FREECAM
	CONSOLE
	SWITCH_COSTUME
)

func Init(window *glfw.Window) {
	window_ = window
	registerCallbacks()
	registerKeyMaps()
}

func registerCallbacks() {
	window_.SetKeyCallback(keyCallback)
}

func registerKeyMaps() {
	//init buttons map for each input context
	ButtonsMap[GAME] = make(map[string]KeyComb)
	ButtonsMap[FREECAM] = make(map[string]KeyComb)
	ButtonsMap[CONSOLE] = make(map[string]KeyComb)
	ButtonsMap[SWITCH_COSTUME] = make(map[string]KeyComb)

	SetInputContext(GAME)
	MapKeyToButton("right", glfw.KeyD, 0)
	MapKeyToButton("left", glfw.KeyA, 0)
	MapKeyToButton("up", glfw.KeyW, 0)
	MapKeyToButton("down", glfw.KeyS, 0)
	MapKeyToButton("jump", glfw.KeySpace, 0)
	MapKeyToButton("mute", glfw.KeyM, 0)
	MapKeyToButton("freecam-on", glfw.KeyKP0, 0)
	MapKeyToButton("cycle-palette", glfw.KeyF3, 0)
	MapKeyToButton("inventory-grid", glfw.KeyTab, 0)
	MapKeyToButton("fly", glfw.KeyT, 0)
	MapKeyToButton("crouch", glfw.KeyC, 0)
	MapKeyToButton("action", glfw.KeyE, 0)
	MapKeyToButton("switch-helmet", glfw.Key0, 0)
	MapKeyToButton("switch-suit", glfw.Key9, 0)
	MapKeyToButton("shoot", glfw.KeyP, 0)
	MapKeyToButton("toggle-zoom", glfw.KeyF2, 0)
	MapKeyToButton("bubbles", glfw.KeyU, 0)
	MapKeyToButton("toggle-texture-filtering", glfw.KeyF6, 0)
	MapKeyToButton("toggle-bbox", glfw.KeyF1, 0)
	MapKeyToButton("cycle-pixel-snap", glfw.KeyF8, 0)
	MapKeyToButton("cycle-camera-snap", glfw.KeyF9, 0)
	MapKeyToButton("toggle-log", glfw.KeyF10, 0)
	MapKeyToButton("set-camera-player", glfw.KeyKP1, 0)
	MapKeyToButton("set-camera-target", glfw.KeyKP2, 0)
	MapKeyToButton("enemy-tool-scroll-down", glfw.KeyDown, 0)
	MapKeyToButton("enemy-tool-scroll-up", glfw.KeyUp, 0)
	MapKeyToButton("switch-skylight", glfw.KeyF11, 0)

	//freecam
	SetInputContext(FREECAM)
	MapKeyToButton("right", glfw.KeyD, 0)
	MapKeyToButton("left", glfw.KeyA, 0)
	MapKeyToButton("up", glfw.KeyW, 0)
	MapKeyToButton("down", glfw.KeyS, 0)
	MapKeyToButton("freecam-off", glfw.KeyKP0, 0)
	MapKeyToButton("toggle-bbox", glfw.KeyF1, 0)
	MapKeyToButton("cycle-pixel-snap", glfw.KeyF8, 0)
	MapKeyToButton("cycle-camera-snap", glfw.KeyF9, 0)

	//revert to game input context
	SetInputContext(GAME)

}

func MapKeyToButton(button string, key glfw.Key, mk glfw.ModifierKey) {
	ActiveButtonsToKeys[button] = KeyComb{key, mk}
}

func SetInputContext(ctx InputContext) {
	prevInputContext = inputContext
	inputContext = ctx
	ActiveButtonsToKeys = ButtonsMap[ctx]
}
func SetPreviousInputContext() {
	inputContext = prevInputContext
	ActiveButtonsToKeys = ButtonsMap[prevInputContext]
}
func GetInputContext() InputContext {
	return inputContext
}

func PrintMk() {
	fmt.Println(modifierKey)
}
