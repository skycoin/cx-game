package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
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
	log.Print("You should see a black cat walking")
	win := render.NewWindow(800, 600, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.
		LoadSpriteSheetByColRow("./assets/blackcat_sprite.png", 13, 4)

	j := 0
	for {
		if window.ShouldClose() {
			break
		}
		spriteloader.LoadSprite(spriteSheetId, "blackcat", 3, j)
		spriteId := spriteloader.GetSpriteIdByName("blackcat")
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		time.Sleep(100 * time.Millisecond)
		spriteloader.DrawSpriteQuad(0, 0, 2, 1, spriteId)
		glfw.PollEvents()
		window.SwapBuffers()
		j++
		if j == 11 {
			j = 0
		}
	}
}
