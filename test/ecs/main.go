package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxecs"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

func init() {
	runtime.LockOSThread()
	rand.Seed(time.Now().Unix())
}
func main() {
	win := render.NewWindow(800, 600, false)
	window := win.Window

	Init(win)

	dt := float32(0)
	lastFrame := float32(glfw.GetTime())

	// spriteId := spriteloader.GetSpriteIdByName("enemy")
	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		dt = currentFrame - lastFrame
		lastFrame = currentFrame

		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw()
		// spriteloader.DrawSpriteQuad(1, 1, 1, 1, spriteId)
		Update(dt)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func Init(win render.Window) {
	spriteloader.InitSpriteloader(&win)
	cxecs.InitDev()
}

func Draw() {

}

func Update(dt float32) {
	//draws as well
	cxecs.UpdateDev(dt)
}
