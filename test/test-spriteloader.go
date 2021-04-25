package main

import (
	"log"
	"runtime"
	"time"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
<<<<<<< HEAD
	win := render.NewWindow(800, 800, true)
=======
	win := render.NewWindow(640,480,true)
>>>>>>> 16cb562... Revert "added simple draw map with move hotkeys"
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.
<<<<<<< HEAD
		LoadSpriteSheetByColRow("../assets/blackcat_sprite.png", 13, 4)

	j := 0
	for {
		if window.ShouldClose() {
			break
		}
		spriteloader.LoadSprite(spriteSheetId, "blackcat", 0, j)
		spriteId := spriteloader.GetSpriteIdByName("blackcat")
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		time.Sleep(100 * time.Millisecond)
		spriteloader.DrawSpriteQuad(0, 0, 1, 1, spriteId)
=======
		LoadSpriteSheet("./assets/starfield/stars/planets.png")
	spriteloader.
		LoadSprite(spriteSheetId, "star", 2,1)
	spriteId := spriteloader.
		GetSpriteIdByName("star")
	for !window.ShouldClose() {
		gl.ClearColor(1,1,1,1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		spriteloader.DrawSpriteQuad(0,0,2,2,spriteId)
>>>>>>> 16cb562... Revert "added simple draw map with move hotkeys"
		glfw.PollEvents()
		window.SwapBuffers()
		j++
		if j == 11 {
			j = 0
		}
	}
}
