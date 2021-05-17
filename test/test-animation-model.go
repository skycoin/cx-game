package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
)

var catBlack *models.CatBlack
var goroutineDelta = make(chan int)

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
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	win := render.NewWindow(400, 600, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	catBlack = models.NewCatBlack(&win, window)
	catBlack.Walk()
}
