package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
)

var catBlack *models.CatBlack

func init() {
	runtime.LockOSThread()
}

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}

	if a == glfw.Press && k == glfw.KeyA {
		// catBlack.Walk()
	}
}

func main() {
	win := render.NewWindow(400, 300, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	catBlack := models.NewCatBlack(&win, window)
	catBlack.Running()
}
