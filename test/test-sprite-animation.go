package main

import (
	"github.com/skycoin/cx-game/spriteloader"
)

var spriteAnimated *spriteloader.SpriteAnimated

func main() {
	// if err := glfw.Init(); err != nil {
	// 	log.Fatalln("failed to initialize glfw:", err)
	// }
	// defer glfw.Terminate()

	// win := render.NewWindow(600, 400, true)
	// window := win.Window
	spriteAnimated = spriteloader.NewSpriteAnimated("./assets/spiderDrill.json")

	// for !window.ShouldClose() {
	// 	gl.ClearColor(1, 1, 1, 1)
	// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// 	glfw.PollEvents()
	// 	window.SwapBuffers()
	// }

}
