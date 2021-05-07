package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	win := render.NewWindow(400, 300, true)
	defer glfw.Terminate()
	catBlack := models.NewCatBlack(&win)
	// catBlack.Walk()
	catBlack.Sit()
	// catBlack.StartRunning()
	// catBlack.Running()
}
