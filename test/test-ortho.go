package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

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

var (
	xpos, ypos, xwidth, ywidth float32 = 1, 1, 1, 1
	// xpos, ypos, xwidth, ywidth float32 = 24, 18, 7, 7
	zoomValue float32 = 0.05
)

func main() {
	win := render.NewWindow(600, 800, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/galaxy_256x256.png")
	spriteloader.
		LoadSprite(spriteSheetId, "star", 6, 5)
	spriteId := spriteloader.
		GetSpriteIdByName("star")
	rand.Seed(time.Now().UnixNano())
	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.2, 0.5, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		spriteloader.DrawSpriteQuadOrtho(xpos, ypos, xwidth, ywidth, spriteId)
		// spriteloader.DrawSpriteQuadOrtho(24, 18, 7, 7, spriteId)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	if a == glfw.Press {
		switch k {
		case glfw.KeyW:
			ypos += 1
		case glfw.KeyS:
			ypos -= 1
		case glfw.KeyA:
			xpos -= 1
		case glfw.KeyD:
			xpos += 1
		}
	}

	//zoom
	if k == glfw.KeyQ {
		xwidth += zoomValue
		ywidth += zoomValue
	}
	if k == glfw.KeyZ {
		xwidth -= zoomValue
		ywidth -= zoomValue
	}

	fmt.Printf("X: %f,     Y:    %f,        Zoom:          %f\n", xpos, ypos, xwidth)
}
