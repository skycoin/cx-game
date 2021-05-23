package main

import (
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/ui"
)

func init() {
	runtime.LockOSThread()
}

var (
	x, y  float32 = 2, 2
	color         = mgl32.Vec4{0.1, 0.5, 0.8, 1.0}
)

func main() {
	window := render.NewWindow(800, 600, false)

	program := window.Program

	glfwWindow := window.Window
	glfwWindow.SetKeyCallback(keyCallBack)
	spriteloader.InitSpriteloader(&window)

	spriteSheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/planets.png")

	spriteloader.LoadSprite(spriteSheetId, "planet", 1, 1)

	spriteId := spriteloader.GetSpriteIdByName("planet")

	ui.InitTextRendering()

	for !glfwWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(program)
		spriteloader.DrawSpriteQuad(1, 1, 1, 1, spriteId)

		world := mgl32.Translate3D(
			x, y, -10,
		)

		ui.DrawString("abcd", world, color, ui.AlignRight)

		glfw.PollEvents()
		window.Window.SwapBuffers()
	}
}

func keyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyA:
		x -= 0.05
	case glfw.KeyD:
		x += 0.05
	case glfw.KeyW:
		y += 0.05
	case glfw.KeyS:
		y -= 0.05
	case glfw.KeyTab:
		if action == glfw.Press {
			color = mgl32.Vec4{rand.Float32(), rand.Float32(), rand.Float32(), rand.Float32()}
		}
	}
}
