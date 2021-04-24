package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called
	// from the main thread.
	runtime.LockOSThread()
}

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}

func main() {
	log.Print("running test")
	log.Print("You should see an orange square rock.")
	win := render.NewWindow(640, 480, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/planets.png")
	spriteloader.LoadSprite(spriteSheetId, "blue_star", 2, 2)
	spriteloader.LoadSprite(spriteSheetId, "purple_star", 2, 0)
	blueStarId := spriteloader.
		GetSpriteIdByName("blue_star")
	purpleStarId := spriteloader.
		GetSpriteIdByName("purple_star")
	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		spriteloader.DrawSpriteQuad(0, 0, 1, 1, blueStarId)
		spriteloader.DrawSpriteQuad(2, 2, 3, 1, purpleStarId)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
