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
	win := render.NewWindow(800, 800, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.
		LoadSpriteSheetByColRow("../assets/blackcat_sprite.png", 13, 4)
	spriteloader.
		LoadSprite(spriteSheetId, "star", 2, 1)
	spriteId := spriteloader.
		GetSpriteIdByName("star")
	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		spriteloader.DrawSpriteQuad(0, 0, 1, 1, spriteId)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
