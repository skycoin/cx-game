package input

//https://www.gamedev.net/blogs/entry/2250186-designing-a-robust-input-handling-system-for-games/

/*
	Performance is important; input lag is a bad thing.
	It should be easy to have new systems tap into the input stream.
	The system must be very flexible and capable of handling a wide variety of game situations.
	Configurability (input mapping) is essential for modern games.
*/

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
)

type MouseCoords struct {
	X float64
	Y float64
}

var (
	window_          *render.Window
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

func Init(window *render.Window) {
	window_ = window
	registerCallbacks()
	registerKeyMaps()
}

func registerCallbacks() {
	window_.Window.SetKeyCallback(keyCallback)
	window_.Window.SetCursorPosCallback(cursorPosCallback)
}

func registerKeyMaps() {
	//init buttons map for each input context
	ButtonsMap[GAME] = make(map[string]glfw.Key)
	ButtonsMap[FREECAM] = make(map[string]glfw.Key)
	ButtonsMap[CONSOLE] = make(map[string]glfw.Key)
	ButtonsMap[SWITCH_COSTUME] = make(map[string]glfw.Key)

	SetInputContext(GAME)
	MapKeyToButton("right", glfw.KeyD)
	MapKeyToButton("left", glfw.KeyA)
	MapKeyToButton("up", glfw.KeyW)
	MapKeyToButton("down", glfw.KeyS)
	MapKeyToButton("jump", glfw.KeySpace)
	MapKeyToButton("mute", glfw.KeyM)
	MapKeyToButton("freecam-on", glfw.KeyKP0)
	MapKeyToButton("cycle-palette", glfw.KeyF3)
	MapKeyToButton("inventory-grid", glfw.KeyI)
	MapKeyToButton("fly", glfw.KeyT)
	MapKeyToButton("crouch", glfw.KeyC)
	MapKeyToButton("action", glfw.KeyE)
	MapKeyToButton("switch-helmet", glfw.Key0)
	MapKeyToButton("switch-suit", glfw.Key9)
	MapKeyToButton("shoot", glfw.KeyP)
	MapKeyToButton("toggle-zoom", glfw.KeyF2)
	MapKeyToButton("toggle-texture-filtering", glfw.KeyF6)
	MapKeyToButton("bubbles", glfw.KeyU)

	MapKeyToButton("enemy-tool-scroll-down", glfw.KeyDown)
	MapKeyToButton("enemy-tool-scroll-up", glfw.KeyUp)

	//freecam
	SetInputContext(FREECAM)
	MapKeyToButton("right", glfw.KeyD)
	MapKeyToButton("left", glfw.KeyA)
	MapKeyToButton("up", glfw.KeyW)
	MapKeyToButton("down", glfw.KeyS)
	MapKeyToButton("freecam-off", glfw.KeyKP0)

	//revert to game input context
	SetInputContext(GAME)

}

func MapKeyToButton(button string, key glfw.Key) {
	ActiveButtonsToKeys[button] = key
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
