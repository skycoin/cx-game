package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
)

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press {
		if k == glfw.KeyEscape {
			//w.SetShouldClose(true)
		}
	}
}

func main() {
	win := render.NewWindow(800, 600, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	// gl.ClearColor(1, 1, 1, 1)
	// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for !window.ShouldClose() {
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
