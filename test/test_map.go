package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world/mainmap"
)

func init() {
	// will crash without this
	runtime.LockOSThread()
}

func main() {
	window := render.NewWindow(800, 800, false)

	window.Window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()

	//init with data
	mainmap.InitMap(&window)

	for !window.Window.ShouldClose() {
		// gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		mainmap.DrawMap()

		glfw.PollEvents()
		window.Window.SwapBuffers()

	}
}

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a != glfw.Press {
		return
	}
	switch k {
	case glfw.KeyW:
		mainmap.GoTop()
	case glfw.KeyS:
		mainmap.GoBottom()
	case glfw.KeyA:
		mainmap.GoLeft()
	case glfw.KeyD:
		mainmap.GoRight()

	case glfw.KeyEscape:
		w.SetShouldClose(true)
	}

}
