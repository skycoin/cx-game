package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

func init() {
	runtime.LockOSThread()
}

const (
	width  = 800
	height = 600
)

var (
	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}
)

func main() {
	win := render.NewWindow(width, height, false)
	window := win.Window
	window.SetKeyCallback(keyCallback)
	// program := win.Program
	defer glfw.Terminate()

	// program := win.Program
	spriteloader.InitSpriteloader(&win)
	vao := spriteloader.MakeQuadVao()
	gl.BindVertexArray(vao)
	starmap.Init(&win)
	starmap.Generate(256, 0.14, 8)

	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		starmap.Draw()

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
func keyCallback(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyTab {
		starmap.Generate(256, 0.08, 4)
	}
}
