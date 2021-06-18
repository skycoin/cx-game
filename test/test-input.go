package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
)

var (
	glfwWindow *glfw.Window
	actionKeys map[glfw.Key]bool = make(map[glfw.Key]bool)
	statekeys  map[glfw.Key]bool = make(map[glfw.Key]bool)
)

func init() {
	runtime.LockOSThread()
}

func main() {
	window := render.NewWindow(800, 600, false)

	glfwWindow = window.Window
	glfwWindow.SetKeyCallback(keyCallback)

	for !window.Window.ShouldClose() {
		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		InputTick()

		glfw.PollEvents()
		window.Window.SwapBuffers()
	}
}
func InputTick() {
	if actionKeys[glfw.KeyT] {
		fmt.Println("PRESSED T!")
	}

	for k := range actionKeys {
		actionKeys[k] = false
	}

}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		actionKeys[key] = true
		statekeys[key] = true
	} else if action == glfw.Release {
		statekeys[key] = false
	}
}
