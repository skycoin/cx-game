package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

var (
	spriteAnimated *spriteloader.SpriteAnimated
	win            render.Window
	window         *glfw.Window
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	win = render.NewWindow(600, 800, true)
	window = win.Window
	spriteloader.InitSpriteloader(&win)
	spriteAnimated = spriteloader.NewSpriteAnimated("./assets/spiderDrill.json")

	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		spriteAnimated.Draw()

		glfw.PollEvents()
		window.SwapBuffers()
	}

}
