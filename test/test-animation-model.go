package main

import (
	// "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
)

func main() {
	win := render.NewWindow(800, 800, true)
	window := win.Window
	// defer glfw.Terminate()
	catBlack := models.NewCatBlack(&window)
	catBlack.walk()
}
