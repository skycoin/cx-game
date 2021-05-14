package main

import (
	"runtime"

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

	if a == glfw.Press && k == glfw.KeyA {
		catBlack.SitStop()
		// go catBlack.Walk()
	}
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	win := render.NewWindow(400, 300, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	catBlack = models.NewCatBlack(&win, window)
	go catBlack.Sit()
	for !window.ShouldClose() {
		glfw.PollEvents()
	}
}
