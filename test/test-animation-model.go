package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/engine/ui"

	//"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
)

var catBlack *models.CatBlack
var goroutineDelta = make(chan int)
var tilePaletteSelector ui.TilePaletteSelector

func init() {
	runtime.LockOSThread()
}

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}

	if a == glfw.Press && k == glfw.KeyS {
		catBlack.Sit()
	}

	if a == glfw.Press && k == glfw.KeyA {
		catBlack.StartRunning()
	}

	if a == glfw.Press && k == glfw.KeyQ {
		catBlack.Running()
	}

	if a == glfw.Press && k == glfw.KeyW {
		catBlack.Walk()
	}
	// if a == glfw.Press && k == glfw.KeyX {
	// 	Cam.Zoom += 0.5
	// }
	// if a == glfw.Press && k == glfw.KeyZ {
	// 	Cam.Zoom -= 0.5
	// }
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	win := render.NewWindow(600, 400, true)
	window := win.Window
	// starmap.Init(&win)
	// starmap.Generate(256, 0.04, 8)
	window.SetKeyCallback(keyCallBack)
	catBlack = models.NewCatBlack(&win, window)
	catBlack.Walk()
	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// starmap.Draw()

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
