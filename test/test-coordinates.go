package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	win := render.NewWindow(800, 600, false)
	window := win.Window
	window.SetKeyCallback(keyCallback)
	window.SetCursorPosCallback(mouseCallback)

	spriteloader.InitSpriteloader(&win)
	spriteId := spriteloader.LoadSingleSprite("./assets/sprite.png", "background")
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		spriteloader.DrawSpriteQuad(0, 0, 1, 1, spriteId)

		spriteloader.DrawSpriteQuad(-5, -5, 1, 1, spriteId)
		spriteloader.DrawSpriteQuad(5, -5, 1, 1, spriteId)
		spriteloader.DrawSpriteQuad(5, 5, 1, 1, spriteId)
		spriteloader.DrawSpriteQuad(-5, 5, 1, 1, spriteId)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

}

func mouseCallback(w *glfw.Window, xpos, ypos float64) {
	
}
